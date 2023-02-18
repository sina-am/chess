package service

import (
	"context"
	"sync"
	"testing"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
)

func TestStartGame(t *testing.T) {
	ctx := context.Background()
	db := database.NewMemoryDatabase(ctx)
	user1 := types.NewUser(
		"test1@gmail.com",
		"test1",
	)
	user2 := types.NewUser(
		"test2@gmail.com",
		"test2",
	)

	db.InsertUser(ctx, user1)
	db.InsertUser(ctx, user2)

	gameSrv, err := NewGameService(db)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := gameSrv.StartGame(ctx, user1)
		assert.Nil(t, err)

	}()
	go func() {
		defer wg.Done()
		_, err := gameSrv.StartGame(ctx, user2)
		assert.Nil(t, err)
	}()

	wg.Wait()

}
