package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	passwd := HashPassword("testPassword")
	assert.Nil(t, VerifyPassword("testPassword", passwd))

	assert.NotNil(t, VerifyPassword("testpassword", passwd))
}
