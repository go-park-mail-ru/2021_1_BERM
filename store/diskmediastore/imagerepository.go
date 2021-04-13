package diskmediastore

import (
	"github.com/pkg/errors"
	"os"
)

type ImageRepository struct {
	workDir string
}

const (
	imageExtend = ".base64"
)

func (i *ImageRepository) GetImage(imageInfo interface{}) ([]byte, error) {
	if imageName := imageInfo.(string); imageName == "" {
		return nil, nil
	}
	imagePath := i.formImagePath(imageInfo.(string))
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, diskDbSourceError)
	}
	fileStat, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, diskDbSourceError)
	}
	image := make([]byte, fileStat.Size())
	_, err = file.Read(image)
	return image, err
}

func (i *ImageRepository) SetImage(imageInfo interface{}, image []byte) (string, error) {
	imagePath := i.formImagePath(imageInfo.(string))
	file, err := os.Create(imagePath)
	defer file.Close()
	if err != nil {
		return "", errors.Wrap(err, diskDbSourceError)
	}
	_, err = file.Write(image)
	return imageInfo.(string), err

}

func (i ImageRepository) formImagePath(imageName string) string {
	var imagePath string
	if imageName[0:1] != "/" {
		imagePath = i.workDir + "/" + imageName + imageExtend
	} else {
		imagePath = i.workDir + imageName + imageExtend
	}
	return imagePath
}
