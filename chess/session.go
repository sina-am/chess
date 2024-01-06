package chess

import (
	"fmt"
	"time"
)

type chessSession struct {
	engine         *chessEngine
	players        map[string]Color
	lastTimePlayed map[Color]time.Time
	remainingTimes map[Color]time.Duration
	tickers        map[Color]*time.Ticker

	finished bool
	winner   Color
}

func NewSession(playersId []string, duration time.Duration) *chessSession {

	engine := NewEngine()
	whiteTimer := time.NewTicker(10 * time.Minute)

	session := &chessSession{
		winner: Empty,
		engine: engine,
		players: map[string]Color{
			playersId[0]: White,
			playersId[1]: Black,
		},
		lastTimePlayed: map[Color]time.Time{
			White: time.Now(),
			Black: {},
		},
		remainingTimes: map[Color]time.Duration{
			White: duration,
			Black: duration,
		},
		tickers: map[Color]*time.Ticker{
			White: whiteTimer,
			Black: nil,
		},
	}

	go session.timeoutTicker(whiteTimer, White)
	return session
}

func (g *chessSession) timeoutTicker(ticker *time.Ticker, playerColor Color) {
	<-ticker.C
	fmt.Printf("%s Time out", playerColor.String())
	g.finished = true
	g.winner = playerColor.OppositeColor()
}

func (g *chessSession) GetWinner() Color {
	if g.winner != Empty {
		return g.winner
	}
	return g.engine.GetWinner()
}

func (g *chessSession) IsFinished() bool {
	return g.finished || g.engine.IsFinished()
}

func (g *chessSession) Play(playerId string, m Move) error {
	playerColor, found := g.players[playerId]
	if !found {
		return fmt.Errorf("player with id %s not found", playerId)
	}
	err := g.engine.Play(playerColor, m)
	if err != nil {
		return err
	}

	elapsed := time.Since(g.lastTimePlayed[playerColor])
	g.remainingTimes[playerColor] -= elapsed
	g.lastTimePlayed[playerColor.OppositeColor()] = time.Now()

	g.tickers[playerColor].Stop()
	if g.tickers[playerColor.OppositeColor()] != nil {
		g.tickers[playerColor.OppositeColor()].Reset(g.remainingTimes[playerColor.OppositeColor()])
	} else {
		g.tickers[playerColor.OppositeColor()] = time.NewTicker(g.remainingTimes[playerColor.OppositeColor()])
		go g.timeoutTicker(g.tickers[playerColor.OppositeColor()], playerColor.OppositeColor())
	}

	fmt.Printf("White time: %s\n", g.remainingTimes[White])
	fmt.Printf("Black time: %s\n", g.remainingTimes[Black])
	return nil
}

func (g *chessSession) Exit(playerId string) (Color, error) {
	playerColor, found := g.players[playerId]
	if !found {
		return 0, fmt.Errorf("player with id %s not found", playerId)
	}

	for _, ticker := range g.tickers {
		ticker.Stop()
	}
	fmt.Printf("Clean exit")
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
