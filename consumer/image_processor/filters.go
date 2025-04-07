package image_processor

import (
	"fmt"
	"github.com/disintegration/imaging"
	"homework-dontpanicw/consumer/domain"
	"log"
)

func RedactPhoto(task domain.Task) error {
	photoName := fmt.Sprintf("repository/upload_photo/photo_storage/%s.png", task.PhotoId.String())
	src, err := imaging.Open(photoName)
	if err != nil {
		return err
	}
	switch task.Filter {
	case "blur":
		redactedImage := imaging.Blur(src, task.Parameter)
		err = imaging.Save(redactedImage, photoName)
		if err != nil {
			return err
		}
	case "sharpen":
		redactedImage := imaging.Sharpen(src, task.Parameter)
		err = imaging.Save(redactedImage, photoName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported filter: %s", task.Filter)
	}
	log.Printf("Redacted photo name: %s.png", task.PhotoId.String())
	return nil

}

//func Negative(id uuid.UUID, parameter float64) error {
//}

//func X(id uuid.UUID, parameter float64) error {
//}
