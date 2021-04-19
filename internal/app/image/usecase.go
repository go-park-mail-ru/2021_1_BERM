package image

type ImageUseCase interface {
	GetImage(imageInfo interface{}) (*models.User, error)
	SetImage(imageInfo interface{}, image []byte) (*models.User, error)
}
