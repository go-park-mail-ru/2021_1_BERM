package errors

import (
	"database/sql"
	"encoding/json"
	"ff/usecase/implementation"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"net/http"
)
var (
	InvalidJSON = &Error{
		Err:  errors.New("Invalid json. "),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
		Field: map[string]interface{}{
			"errors": "Invalid json",
		},
	}

	InvalidCookies = &Error{
		Err:  errors.New("Invalid cookie.\n"),
		Code: http.StatusBadRequest,
		Type: TypeExternal,
	}
)
const (
	diskDbSourceError  = "Disc db sb source errors"
	DuplicateErrorCode = "23505"
	SqlDbSourceError   = "SQL sb source errors"
)
var(
	NotAuthorized = errors.New("Not authorized.")
)

type DuplicateSourceErr struct {
	Err error
}

func (e *DuplicateSourceErr) Error() string {
	return e.Err.Error()
}


const (
	TypeInternal = "Internal"
	TypeExternal = "External"
)

type Error struct {
	Err   error
	Type  string
	Field map[string]interface{}
	Code  int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func New(err error) *Error {

	if errors.Is(err, sql.ErrNoRows) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusNotFound,
			Field: map[string]interface{}{
				"errors": sql.ErrNoRows,
			},
		}
	}
	dup := DuplicateSourceErr{}
	if errors.As(err, &dup) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusBadRequest,
			Field: map[string]interface{}{
				"errors": "Duplicate errors",
			},
		}
	}
	valid := &validation.Errors{}
	if errors.As(err, valid) {
		j, errJ := valid.MarshalJSON()
		if errJ != nil {
			return &Error{
				Err:  errors.Wrap(errJ, "Error json marshal on create htttp errors"),
				Type: TypeInternal,
				Code: http.StatusInternalServerError,
				Field: map[string]interface{}{
					"errors": "Intertnal server errors",
				},
			}
		}
		field := make(map[string]interface{})
		errJ = json.Unmarshal(j, &field)
		if errJ != nil {
			return &Error{
				Err:  errors.Wrap(errJ, "Error json marshal on create htttp errors"),
				Type: TypeInternal,
				Code: http.StatusInternalServerError,
				Field: map[string]interface{}{
					"errors": "Intertnal server errors",
				},
			}
		}
		return &Error{
			Err:   err,
			Type:  TypeExternal,
			Code:  http.StatusBadRequest,
			Field: field,
		}
	}
	if errors.Is(err, NotAuthorized) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusUnauthorized,
			Field: map[string]interface{}{
				"Error": "Not authorized",
			},
		}
	}
	if errors.Is(err, implementation.ErrBadPassword) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusBadRequest,
			Field: map[string]interface{}{
				"Error": implementation.ErrBadPassword.Error(),
			},
		}
	}
	return &Error{
		Err:  err,
		Type: TypeInternal,
		Code: http.StatusInternalServerError,
		Field: map[string]interface{}{
			"errors": "Intertnal server errors",
		},
	}
}
