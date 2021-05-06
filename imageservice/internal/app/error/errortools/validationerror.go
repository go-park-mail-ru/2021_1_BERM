package errortools

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"net/http"
)

func validationErrorHandle(err error) (map[string]interface{}, int, bool) {
	validErr := validation.Errors{} //т.к это мэпа, по идее это уже указатель,& не ставлю
	if errors.As(err, &validErr) {
		return map[string]interface{}{
			"message": validErr,
		}, http.StatusBadRequest, true
	}
	return nil, 0, false
}
