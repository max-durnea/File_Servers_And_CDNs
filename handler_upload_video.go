package main

import (
	"net/http"
    "os"
	"fmt"
	"mime"
	"io"
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w,r.Body,1<<30)

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video_meta_data, err:= cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, 400, "Video does not exist", err)
		return
	}
	if video_meta_data.UserID != userID {
		respondWithJSON(w, http.StatusUnauthorized, struct{}{})
		return
	}
	file,header,err := r.FormFile("video")
	if err != nil {
		respondWithError(w, 400, "No video present", err)
		return
	}
	defer file.Close()
	content_type := header.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(content_type)
	if mediatype != "video/mp4"{
		respondWithError(w,400, "Wrong format", err)
		return
	}
	new_file, err := os.CreateTemp("","tubely-upload.mp4")
	if err != nil {
		respondWithError(w, 500, "Could not create file", err)
		return
	}
	defer os.Remove(new_file.Name())
	defer new_file.Close()
	io.Copy(new_file,file)
	new_file.Seek(0,io.SeekStart)
	aspect_ratio,_ := getVideoAspectRatio(new_file.Name())
	if aspect_ratio == "16:9"{
		aspect_ratio = "landscape"
	} else if aspect_ratio == "9:16"{
		aspect_ratio = "portrait"
	}
	fastStart,_:=processVideoForFastStart(new_file.Name())
	fastStartFile,_:=os.Open(fastStart)
	defer fastStartFile.Close()
	var bytes []byte = make([]byte,32)
	rand.Read(bytes)
	string_encoding := base64.RawURLEncoding.EncodeToString(bytes)
	file_name:= fmt.Sprintf("%v/%v.mp4",aspect_ratio,string_encoding)
	putObjectParams := s3.PutObjectInput{
		Bucket : &cfg.s3Bucket,
		Key : &file_name,
		Body : fastStartFile,
		ContentType : &mediatype,
	}
	_,err = cfg.s3Client.PutObject(r.Context(),&putObjectParams)
	if err != nil {
		respondWithError(w, 500, "Could not create file", err)
		return
	}
	file_url := fmt.Sprintf("%v/%v",cfg.s3CfDistribution,file_name)
	video_meta_data.VideoURL = &file_url
	cfg.db.UpdateVideo(video_meta_data)
	respondWithJSON(w, http.StatusOK,video_meta_data)
}
