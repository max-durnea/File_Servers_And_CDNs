package main

import (
	"fmt"
	"net/http"
	"io"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"strings"
	"path/filepath"
	"os"
	"mime"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
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


	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	//10MB max file data stored in RAM
	var maxMemory int64 = 10 << 20
	r.ParseMultipartForm(maxMemory)
	//Get the file data and headers:
	file,header,err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, 400, "No thumbnail present", err)
		return
	}
	content_type := header.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(content_type)
	if mediatype != "image/jpeg" && mediatype != "image/png"{
		respondWithError(w,400, "Wrong format", err)
		return
	}
	extension := strings.Split(content_type,"/")[1]
	path:=filepath.Join(cfg.assetsRoot,fmt.Sprintf("%v",videoID))
	new_file, err := os.Create(path+"."+extension)
	io.Copy(new_file,file)
	if err != nil {
		respondWithError(w, 500, "Could not create file", err)
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
	
	//encoded_image := base64.StdEncoding.EncodeToString(image_data)
	data_URL := fmt.Sprintf("http://localhost:%v/assets/%v.%v",cfg.port,fmt.Sprintf("%v",videoID),extension)
	video_meta_data.ThumbnailURL =  &data_URL
	cfg.db.UpdateVideo(video_meta_data)
	respondWithJSON(w, http.StatusOK,video_meta_data)
}
