package game

import (
	"github.com/sina-am/chess/chess"
	"go.uber.org/zap"
)

type loggerHandler struct {
	logger *zap.Logger
	next   GameHandler
}

func NewLoggerHandler(logger *zap.Logger, next GameHandler) GameHandler {
	return &loggerHandler{
		logger: logger,
		next:   next,
	}
}

func (h *loggerHandler) Start() {
	h.logger.Info("server is starting")
	h.next.Start()
}

func (h *loggerHandler) Register(p *player) {
	h.logger.Info("new player connected", zap.String("playerID", p.GetId()))
	h.next.Register(p)
}

func (h *loggerHandler) UnRegister(p *player) {
	h.logger.Info("player disconnected", zap.String("playerID", p.GetId()))
	h.next.UnRegister(p)
}
func (h *loggerHandler) Play(p *player, gameId string, move chess.Move) {
	h.next.Play(p, gameId, move)
}

func (h *loggerHandler) ExitGame(gameId string, p *player) {
	h.logger.Info(
		"player exited game",
		zap.String("playerID", p.GetId()),
		zap.String("gameID", gameId),
	)
	h.next.ExitGame(gameId, p)
}

func (h *loggerHandler) AddToWaitList(p *player, gs GameSetting) {
	h.logger.Info(
		"player added to wait list",
		zap.String("playerID", p.GetId()),
	)
	h.next.AddToWaitList(p, gs)
}
func (h *loggerHandler) RemoveFromWaitList(p *player) {
	h.logger.Info(
		"player removed from wait list",
		zap.String("playerID", p.GetId()),
	)
	h.next.RemoveFromWaitList(p)
}
