package report

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	_ "image/png"
)

func embedScreenshot(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if jpg := pngToJPEGDataURL(data, 72); jpg != "" {
		return jpg
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
}

func pngToJPEGDataURL(pngData []byte, quality int) string {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	img, err := pngDecode(pngData)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return ""
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func pngDecode(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}
