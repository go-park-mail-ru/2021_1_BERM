package sessiontools_test

import (
	"authorizationservice/internal/app/models"
	"authorizationservice/internal/app/session/tools/sessiontools"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateSession(t *testing.T) {
	session := models.Session{
		UserId:   1,
		Executor: true,
	}

	sesTools := sessiontools.SessionTools{}
	_, err := sesTools.BeforeCreate(session)
	require.NoError(t, err)
}
