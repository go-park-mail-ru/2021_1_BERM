package imgcreator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

func randomFilename16Char() (s string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	s = fmt.Sprintf("%x", b)
	return
}

func CreateImg(imgBase64 string) (string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgBase64))
	m, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	salt := make([]byte, 8)
	_, err = rand.Read(salt)
	if err != nil {
		return "", err
	}
	jpegFilename, err := randomFilename16Char()
	jpegFilename += ".jpeg"
	f, err := os.Create(jpegFilename)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", err
	}
	return jpegFilename, nil
}
