package image

type ImageRepository interface {
	GetImage(imageInfo interface{}) ([]byte, error)
	SetImage(imageInfo interface{}, image []byte) (string, error)
}

