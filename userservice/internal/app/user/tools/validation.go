package tools

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"user/Error"
	"user/internal/app/models"
)

const (
	minPasswordLength = 5
	maxPasswordLength = 30
)

func ValidationCreateUser(user *models.NewUser) error {
	err := validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(minPasswordLength, maxPasswordLength)),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.NameSurname, validation.Required),
	)
	if err != nil {
		extendError := &Error.Error{
			InternalError: false,
			Err:           errors.New(err.Error()),
		}
		errDescription, err := validationErrorToMap(err)
		if err != nil {
			return err
		}
		extendError.ErrorDescription = *errDescription
		return extendError
	}
	return nil
}

func ValidationChangeUser(user *models.ChangeUser) error {

	err := validation.ValidateStruct(
		user,
		validation.Field(&user.Email, is.Email),
		validation.Field(&user.Password, validation.Length(minPasswordLength, maxPasswordLength)),
	)
	if err != nil {
		extendError := &Error.Error{
			InternalError: false,
			Err:           errors.New(err.Error()),
		}
		errDescription, err := validationErrorToMap(err)
		if err != nil {
			return err
		}
		extendError.ErrorDescription = *errDescription
		return extendError
	}
	return nil
}

func validationErrorToMap(err error) (*map[string]interface{}, error) {
	validationError := &validation.Errors{}
	if !errors.Is(err, validationError) {
		return nil, &Error.Error{
			InternalError:    true,
			ErrorDescription: Error.InternalServerErrorDescription,
			Err:              errors.New("Is not validation error"),
		}
	}
	errorData, err := validationError.MarshalJSON()
	if err != nil {
		return nil, &Error.Error{
			InternalError:    true,
			ErrorDescription: Error.InternalServerErrorDescription,
			Err:              err,
		}
	}
	errorDescription := &map[string]interface{}{}
	err = json.Unmarshal(errorData, errorDescription)
	if err != nil {
		return nil, &Error.Error{
			InternalError:    true,
			ErrorDescription: Error.InternalServerErrorDescription,
			Err:              err,
		}
	}
	return errorDescription, nil
}
