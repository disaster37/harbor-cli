package harbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {

	cfg := Config{
		Address:  "http://localhost",
		Username: "test",
		Password: "test",
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
