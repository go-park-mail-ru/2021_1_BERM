package diskmediastore

import "FL_2/store"

type MediaStore struct {
	imageRepository *ImageRepository
	workDir string
}

func New(workDir string) *MediaStore {
	s := &MediaStore{
		workDir: workDir,
	}
	return s
}

func (s *MediaStore)Image() store.ImageRepository{
	if s.imageRepository == nil{
		s.imageRepository = &ImageRepository{
			workDir: s.workDir,
		}
	}
	return s.imageRepository
}