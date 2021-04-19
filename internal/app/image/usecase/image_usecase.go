package usecase

const (
	mediaUseCaseError = "Media use case errors."
)

type ImageUseCase struct {
	store ImageRepository
}

func (s *MediaUseCase) SetImage(imageInfo interface{}, image []byte) (*models.User, error) {
	u := imageInfo.(*models.User)
	sanitizer := bluemonday.UGCPolicy()
	image = sanitizer.SanitizeBytes(image)
	imageID, err := s.mediaStore.Image().SetImage(strconv.FormatUint(u.ID, 10), image)
	if err != nil {
		return nil, errors.Wrap(err, mediaUseCaseError)
	}
	u.Img = imageID
	u, err = s.store.User().ChangeUser(*u)
	if err != nil {
		return nil, errors.Wrap(err, mediaUseCaseError)
	}
	u.Img = string(image)
	return u, nil
}
