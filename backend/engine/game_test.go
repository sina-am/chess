package engine

import (
	"context"
	"sync"
	"testing"

	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetPlayer(t *testing.T) {
	ctx := context.Background()
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	userId := primitive.NewObjectID()
	game.Join(ctx, userId)

	player, err := game.getPlayer(userId)
	assert.Nil(t, err)
	assert.Equal(t, userId, player.UserId)
}

func TestPlayersColor(t *testing.T) {
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	_, err = game.Join(context.TODO(), primitive.NewObjectID())
	assert.Nil(t, err)

	_, err = game.Join(context.TODO(), primitive.NewObjectID())
	assert.Nil(t, err)

	assert.NotEqual(t, game.game.Players[0].Color, game.game.Players[1])
}
func TestPlayersTurn(t *testing.T) {
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	_, err = game.Join(context.TODO(), primitive.NewObjectID())
	assert.Nil(t, err)
	_, err = game.Join(context.TODO(), primitive.NewObjectID())
	assert.Nil(t, err)

	for _, player := range game.game.Players {
		if player.Color == types.White {
			assert.True(t, player.Turn)
		} else {
			assert.False(t, player.Turn)
		}
	}
}
func TestJoin(t *testing.T) {
	ctx := context.Background()
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	t.Run("Join one player", func(t *testing.T) {
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)

		assert.Nil(t, err)
		assert.Equal(t, userId, player.id)
		assert.Equal(t, game.players[0].id, userId)
		assert.False(t, game.started)
	})

	t.Run("Join next player", func(t *testing.T) {
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)

		assert.Nil(t, err)
		assert.Equal(t, userId, player.id)
		assert.Equal(t, game.players[1].id, userId)
		assert.True(t, game.started)
	})
	// t.Run("No more than to player", func(t *testing.T) {
	// 	userId := primitive.NewObjectID()
	// 	_, err := game.Join(ctx, userId)
	// 	assert.NotNil(t, err)
	// })
}

func TestStart(t *testing.T) {
	ctx := context.Background()
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)
		assert.Nil(t, err)

		player.WaitForStart(ctx)
	}()
	go func() {
		defer wg.Done()
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)
		assert.Nil(t, err)

		player.WaitForStart(ctx)
	}()

	wg.Wait()
	assert.True(t, game.started)
}

func TestPlay(t *testing.T) {
	ctx := context.Background()
	game, err := NewOnlineGame(10)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)
		assert.Nil(t, err)

		assert.Nil(t, player.WaitForStart(ctx))

		assert.Nil(t, player.WaitForMyTurn(ctx))
		err = game.Play(ctx, player.id, types.Location{Row: 6, Col: 1}, types.Location{Row: 5, Col: 1})
		assert.Nil(t, err)

	}()
	go func() {
		defer wg.Done()
		userId := primitive.NewObjectID()
		player, err := game.Join(ctx, userId)
		assert.Nil(t, err)

		assert.Nil(t, player.WaitForStart(ctx))

		assert.Nil(t, player.WaitForMyTurn(ctx))
		err = game.Play(ctx, player.id, types.Location{Row: 1, Col: 1}, types.Location{Row: 2, Col: 1})
		assert.Nil(t, err)
	}()

	wg.Wait()
}
