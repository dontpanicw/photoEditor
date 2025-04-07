package upload_photo

import (
	"github.com/google/uuid"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

func UploadPhoto(file io.Reader, id uuid.UUID) error {
	const uploadPath = "homework-dontpanicw/app/repository/upload_photo"

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	filename := id.String() + ".png"
	filePath := filepath.Join(uploadPath, filename)

	outfile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		return err
	}
	return nil
}
