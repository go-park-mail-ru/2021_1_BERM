package errortools

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestTools(t *testing.T) {
	a, b := ErrorHandle(errors.New("kek"))
	require.Equal(t, a, map[string]interface{}(map[string]interface{}{"message": "Ooops. Something went wrong!!! :("}))
	require.Equal(t, b, http.StatusInternalServerError)
}
