package toerr_test

import (
	"master-blaster/internal/toerr"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func CallLevel1() error {

	err := CallLevel2()
	if err != nil {
		return errors.Wrap(err, "CallLevel1")
	}
	return nil
}

func CallLevel2() error {
	err := CallLevel3()
	if err != nil {
		return errors.Wrap(err, "CallLevel2")
	}
	return nil
}

func CallLevel3() error {
	return toerr.NewValidationError(http.StatusBadRequest, "Bad parameters passed")
}

func TestErrorUnwrapping(t *testing.T) {
	err := CallLevel1()
	assert.NotNil(t, err)
	verr := toerr.ExtractValidationError(err)
	assert.NotNil(t, verr)

	assert.Equal(t, verr.HttpErrorCode, http.StatusBadRequest)

	assert.NotEmpty(t, verr.Message)

}

func Level1() error {
	err := Level2()
	if err != nil {
		return errors.Wrap(err, "Level 1")
	}
	return nil
}

func Level2() error {
	return errors.New("Level2 error")
}

func TestErrorUnwrappingUnknownError(t *testing.T) {
	err := Level1()
	assert.NotNil(t, err)
	verr := toerr.ExtractValidationError(err)
	assert.Nil(t, verr)

}
