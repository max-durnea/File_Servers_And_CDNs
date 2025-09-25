package main

/*import (
    "github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
    "strings"
    "time"
    "fmt"
    "net/url"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
    if video.VideoURL == nil {
        return video, fmt.Errorf("video URL is nil")
    }

    urlStr := *video.VideoURL
    var bucket, key string

    // Try old format first
    parts := strings.Split(urlStr, ",")
    if len(parts) == 2 {
        bucket, key = parts[0], parts[1]
    } else {
        // fallback: parse full S3 URL
        u, err := url.Parse(urlStr)
        if err != nil {
            return video, fmt.Errorf("invalid video URL: %v", urlStr)
        }
        hostParts := strings.Split(u.Host, ".")
        if len(hostParts) < 3 {
            return video, fmt.Errorf("invalid s3 host in URL: %v", urlStr)
        }
        bucket = hostParts[0]
        key = strings.TrimPrefix(u.Path, "/")
    }

    presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Minute)
    if err != nil {
        return database.Video{}, err
    }

    video.VideoURL = &presignedURL
    return video, nil
}
*/