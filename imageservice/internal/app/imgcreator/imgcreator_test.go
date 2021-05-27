package imgcreator_test

import (
	"bufio"
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"imageservice/internal/app/imgcreator"
	"os"
	"testing"
)

func TestImgCreator_CreateImg(t *testing.T) {
	imgCr := imgcreator.ImgCreator{}

	imgFile, _ := os.Open("image/kek.jpeg")

	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	_, _ = fReader.Read(buf)

	img := base64.StdEncoding.EncodeToString(buf)
	_, err := imgCr.CreateImg(img)
	require.NoError(t, err)
}

func TestImgCreator_CreateImgErr(t *testing.T) {
	imgCr := imgcreator.ImgCreator{}
	img := "/9j/4"
	_, err := imgCr.CreateImg(img)
	require.Error(t, err)
}

func TestImgCreator_CropImg(t *testing.T) {
	imgCr := imgcreator.ImgCreator{}
	imgURL := "kek.jpeg"
	imgPath, err := imgCr.CropImg(imgURL)
	require.NoError(t, err)
	require.Equal(t, imgURL, imgPath)
}

func TestImgCreator_CropImgErr(t *testing.T) {
	imgCr := imgcreator.ImgCreator{}
	imgURL := "kek1234.jpeg"
	_, err := imgCr.CropImg(imgURL)
	require.Error(t, err)
}
