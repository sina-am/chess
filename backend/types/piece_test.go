package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPieceString(t *testing.T) {
	piece := Piece{
		Color: Black,
		Type:  King,
	}
	assert.Equal(t, piece.String(), "♔")
}
