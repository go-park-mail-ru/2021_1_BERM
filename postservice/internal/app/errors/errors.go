package errors

import (
	"FL_2/store/postgresstore"
	"FL_2/store/tarantoolcache"
	"FL_2/usecase/implementation"
	"database/sql"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"net/http"
)
const (
	duplicateErrorCode = "23505"
	sqlDbSourceError   = "SQL sb source error"
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
				"error": sql.ErrNoRows,
			},
		}
	}
	//TODO: Ошибки
	dup := &postgresstore.DuplicateSourceErr{}
	if errors.As(err, &dup) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusBadRequest,
			Field: map[string]interface{}{
				"error": "Duplicate error",
			},
		}
	}
	valid := &validation.Errors{}
	if errors.As(err, valid) {
		j, errJ := valid.MarshalJSON()
		if errJ != nil {
			return &Error{
				Err:  errors.Wrap(errJ, "Error json marshal on create htttp error"),
				Type: TypeInternal,
				Code: http.StatusInternalServerError,
				Field: map[string]interface{}{
					"error": "Intertnal server error",
				},
			}
		}
		field := make(map[string]interface{})
		errJ = json.Unmarshal(j, &field)
		if errJ != nil {
			return &Error{
				Err:  errors.Wrap(errJ, "Error json marshal on create htttp error"),
				Type: TypeInternal,
				Code: http.StatusInternalServerError,
				Field: map[string]interface{}{
					"error": "Intertnal server error",
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
	//TODO: Ошибки
	if errors.Is(err, tarantoolcache.NotAuthorized) {
		return &Error{
			Err:  err,
			Type: TypeExternal,
			Code: http.StatusUnauthorized,
			Field: map[string]interface{}{
				"Error": "Not authorized",
			},
		}
	}
	//TODO: Ошибки
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
			"error": "Intertnal server error",
		},
	}
}



type DuplicateSourceErr struct {
	Err error
}

func (e *DuplicateSourceErr) Error() string {
	return e.Err.Error()
}

