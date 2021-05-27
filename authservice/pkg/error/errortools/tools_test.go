package errortools_test

import (
	errTools "authorizationservice/pkg/error/errortools"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestTools(t *testing.T) {
	a, b := errTools.ErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}{"message": "Ooops. Something went wrong!!! :("})
	require.Equal(t, b, http.StatusInternalServerError)
}
