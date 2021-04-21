package store

//go:generate mockgen -destination=mock/mock_media_store.go -package=mock FL_2/store MediaStore
type MediaStore interface {
	Image() ImageRepository
}
