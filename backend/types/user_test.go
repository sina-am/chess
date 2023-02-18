package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrationRequest(t *testing.T) {
	NewValidator()
	t.Run("Valid email", func(t *testing.T) {
		req := RegistrationRequest{
			Email:    "test@gmail.com",
			Password: "strongpassword",
		}

		assert.Nil(t, req.Validate())
	})
	t.Run("Invalid email", func(t *testing.T) {
		req := RegistrationRequest{
			Email:    "invalidgmail.com",
			Password: "strongpassword",
		}

		assert.NotNil(t, req.Validate())
	})
	t.Run("Short password", func(t *testing.T) {
		req := RegistrationRequest{
			Email:    "test@gmail.com",
			Password: "test",
		}

		assert.NotNil(t, req.Validate())
	})
}
