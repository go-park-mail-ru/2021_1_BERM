package errortools_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	errTools "authorizationservice/pkg/error/errortools"
)

func TestTools(t *testing.T) {
	a, b := errTools.ErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}{"message": "Ooops. Something went wrong!!! :("})
	require.Equal(t, b, http.StatusInternalServerError)
}
