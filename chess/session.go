package chess

import (
	"time"
)

type chessSession struct {
	*ChessEngine
	lastTimePlayed map[Color]time.Time
	remainingTimes map[Color]time.Duration
	tickers        map[Color]*time.Ticker
}

func NewSession(duration time.Duration) *chessSession {
	engine := NewEngine()
	whiteTimer := time.NewTicker(10 * time.Minute)

	session := &chessSession{
		ChessEngine: engine,
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
	g.finish(Timeout, playerColor.OppositeColor())
}

func (g *chessSession) Play(playerColor Color, m Move) error {
	if g.finished {
		return ErrGameEnd
	}

	err := g.ChessEngine.Play(playerColor, m)
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

	return nil
}

func (g *chessSession) Exit() {
	for _, ticker := range g.tickers {
		if ticker != nil {
			ticker.Stop()
		}
	}
	g.finished = true
}
