package diskmediastore

import (
	"os"
)

type ImageRepository struct {
	workDir string
}
const (
	imageExtend = ".base64"
)
func (i *ImageRepository)GetImage(imageInfo interface{}) ([]byte, error){
	imagePath := i.formImagePath(imageInfo.(string))
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil{
		return nil, err
	}
	fileStat, err := file.Stat()
	if err != nil{
		return nil, err
	}
	image := make([]byte, fileStat.Size())
	_, err = file.Read(image)
	return image, err
}

func (i *ImageRepository)SetImage(imageInfo interface{}, image []byte) (string, error){
	imagePath := i.formImagePath(imageInfo.(string) + imageExtend)
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		return "", err
	}
	_, err = file.Write(image)
	return  imagePath, err

}

func (i ImageRepository)formImagePath(imageName string) string{
	var imagePath string
	if imageName[0:1] != "/"{
		imagePath = i.workDir + "/" + imageName
	} else{
		imagePath = i.workDir + imageName
	}
	return imagePath
}
