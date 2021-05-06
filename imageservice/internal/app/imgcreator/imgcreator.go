package imgcreator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type ImgCreator struct {}

func randomFilename16Char() (s string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	s = fmt.Sprintf("%x", b)
	return
}

func (i *ImgCreator) CreateImg(imgBase64 string) (string, error) {
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
	f, err := os.Create("image/" + jpegFilename)
	if err != nil {
		return "", err
	}
	err = f.Chmod(0777)
	if err != nil {
		return "", err
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", err
	}
	return jpegFilename, nil
}

func (i *ImgCreator) CropImg(imgURL string) (string, error) {
	fi, err := os.Open("image/" + imgURL)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(fi)
	if err != nil {
		return "", err
	}

	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, err := analyzer.FindBestCrop(img, 250, 250)
	if err != nil {
		return "", err
	}
	fi.Close()
	sub, ok := img.(SubImager)
	if ok {
		fi, err := os.Create("image/" + imgURL)
		if err != nil {
			return "", err
		}
		cropImage := sub.SubImage(topCrop)
		if err := jpeg.Encode(fi, cropImage, nil); err != nil {
			return "", err
		}
		return imgURL, nil
	} else {
		return "", err
	}
}
