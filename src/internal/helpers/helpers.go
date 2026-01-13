package helpers

import (
	"bytes"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

// GenerateThumbnail creates a thumbnail with max width of 300px
func GenerateThumbnail(imageData []byte, mimeType string) ([]byte, error) {
	// Decode the image
	img, err := imaging.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	// Resize to max width 300px (maintains aspect ratio)
	thumbnail := imaging.Resize(img, 300, 0, imaging.Lanczos)

	// Encode back to bytes
	var buf bytes.Buffer

	switch mimeType {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 85})
	case "image/png":
		err = png.Encode(&buf, thumbnail)
	default:
		// Default to JPEG for other formats
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
