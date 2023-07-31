package image

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
)

var ImageContentTypes = []string{"image/png", "image/jpg", "image/jpeg"}

func IsTypeAllowed(contentType string) bool {
	for _, v := range ImageContentTypes {
		if v == contentType {
			return true
		}
	}

	return false
}

func IsSizeAllowed(size int64) bool {
	return size < (1000000 * config.IMAGE_MAX_SIZE_IN_MB)
}

func IsDimentionAllowed(r *http.Request) bool {
	// Need to get file again from http request, because if we dont it make demaged on image when saved
	src, _, _ := r.FormFile("image")
	img, _, _ := image.DecodeConfig(src)
	defer src.Close()

	if img.Width < config.IMAGE_MIN_WIDTH || img.Height < config.IMAGE_MIN_HEIGHT {
		return false
	}
	if img.Width > config.IMAGE_MAX_WIDTH || img.Height > config.IMAGE_MAX_HEIGHT {
		return false
	}

	return true
}
