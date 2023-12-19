package chess

import "fmt"

type chessSession struct {
	engine  *chessEngine
	players map[string]Color
}

func NewSession(playersId []string) *chessSession {
	engine := NewEngine()
	return &chessSession{
		engine: engine,
		players: map[string]Color{
			playersId[0]: White,
			playersId[1]: Black,
		},
	}
}

func (g *chessSession) GetWinner() Color {
	return g.engine.GetWinner()
}

func (g *chessSession) IsFinished() bool {
	return g.engine.IsFinished()
}

func (g *chessSession) Play(playerId string, m Move) error {
	playerColor, found := g.players[playerId]
	if !found {
		return fmt.Errorf("player with id %s not found", playerId)
	}
	return g.engine.Play(playerColor, m)
}

func (g *chessSession) Exit(playerId string) (Color, error) {
	playerColor, found := g.players[playerId]
	if !found {
		return 0, fmt.Errorf("player with id %s not found", playerId)
	}

	return playerColor.OppositeColor(), nil
}

func (g *chessSession) GetPlayers() []string {
	var ids []string
	for id := range g.players {
		ids = append(ids, id)
	}
	return ids
}

func (g *chessSession) GetPlayerByColor(color Color) string {
	for player, playerColor := range g.players {
		if playerColor == color {
			return player
		}
	}
	panic("invalid color")
}

func (g *chessSession) InGame(playerId string) bool {
	if _, found := g.players[playerId]; found {
		return true
	}
	return false
}
