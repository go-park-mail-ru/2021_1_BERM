package store

//go:generate mockgen -destination=mock/mock_image_repo.go -package=mock FL_2/store ImageRepository
type ImageRepository interface {
	GetImage(imageInfo interface{}) ([]byte, error)
	SetImage(imageInfo interface{}, image []byte) (string, error)
}
