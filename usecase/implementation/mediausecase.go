package implementation

import (
	"FL_2/model"
	"FL_2/store"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"strconv"
)

const (
	mediaUseCaseError = "Media use case error."
)

type MediaUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (s *MediaUseCase) GetImage(imageInfo interface{}) (*model.User, error) {
	return nil, nil
}

func (s *MediaUseCase) SetImage(imageInfo interface{}, image []byte) (*model.User, error) {
	u := imageInfo.(*model.User)
	sanitizer := bluemonday.UGCPolicy()
	image = sanitizer.SanitizeBytes(image)
	imageID, err := s.mediaStore.Image().SetImage(strconv.FormatUint(u.ID, 10), image)
	if err != nil {
		return nil, errors.Wrap(err, mediaUseCaseError)
	}
	u.Img = imageID
	oldUser, err := s.store.User().FindUserByID(u.ID)
	if err != nil {
		return nil, errors.Wrap(err, mediaUseCaseError)
	}
	if u.Email == "" {
		u.Email = oldUser.Email
	}

	if u.About == "" {
		u.About = oldUser.About
	}

	if u.Password == "" {
		u.Password = oldUser.Password
	}

	if u.Login == "" {
		u.Login = oldUser.Login
	}

	if u.Img == "" {
		u.Img = oldUser.Img
	}

	if u.NameSurname == "" {
		u.NameSurname = oldUser.NameSurname
	}

	if u.Rating == 0 {
		u.Rating = oldUser.Rating
	}

	u.Executor = oldUser.Executor

	for _, spec := range oldUser.Specializes {
		u.Specializes = append(u.Specializes, spec)
	}
	u, err = s.store.User().ChangeUser(*u)
	if err != nil {
		return nil, errors.Wrap(err, mediaUseCaseError)
	}
	u.Img = string(image)
	return u, nil
}
