package imgcreator

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock imageservice/internal/app/imgcreator ImgCreatorI
type ImgCreatorI interface {
	CreateImg(imgBase64 string) (string, error)
	CropImg(imgURL string) (string, error)
}
