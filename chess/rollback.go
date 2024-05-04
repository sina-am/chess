package chess

import "math"

type RollBackMovement struct {
	game          *ChessEngine
	capturedPiece *Piece
	move          Move
	castled       bool
	rolledBack    bool
	promoted      bool

	castleRightsBackup map[Color]castling
}

func NewRollBack(game *ChessEngine) *RollBackMovement {
	return &RollBackMovement{
		game: game,
		castleRightsBackup: map[Color]castling{
			White: game.castleRights[White],
			Black: game.castleRights[Black],
		},
	}
}

func (r *RollBackMovement) CheckPromotion(move Move) {
	piece := r.game.board[move.From.Row][move.From.Col]
	if piece.Type == Pawn {
		if piece.Color == Black && move.To.Row == 0 {
			piece.Type = Queen
			r.promoted = true
		} else if move.To.Row == 7 {
			piece.Type = Queen
			r.promoted = true
		}
	}
}
func (r *RollBackMovement) isCastling() bool {
	if (r.game.board[r.move.From.Row][r.move.From.Col].Type == King) && (math.Abs(float64(r.move.From.Col)-float64(r.move.To.Col)) == 2) {
		return true
	}
	return false
}

func (r *RollBackMovement) doCastling() {
	color := r.game.board[r.move.From.Row][r.move.From.Col].Color
	backRow := 0
	if color == Black {
		backRow = 7
	}

	if r.move.To.Col == 6 {
		r.game.castleRights[color] = castling{
			left:  r.game.castleRights[color].left,
			right: false,
		}
		r.game.board[backRow][5] = r.game.board[backRow][7]
		r.game.board[backRow][6] = r.game.board[backRow][4]
		r.game.board[backRow][4] = nil
		r.game.board[backRow][7] = nil

		r.game.board[backRow][5].Location.Col = 5
		r.game.board[backRow][6].Location.Col = 6
		return
	}
	if r.move.To.Col == 2 {
		r.game.castleRights[color] = castling{
			left:  false,
			right: r.game.castleRights[color].right,
		}
		r.game.board[backRow][3] = r.game.board[backRow][0]
		r.game.board[backRow][2] = r.game.board[backRow][4]
		r.game.board[backRow][4] = nil
		r.game.board[backRow][0] = nil

		r.game.board[backRow][3].Location.Col = 3
		r.game.board[backRow][2].Location.Col = 2
	}
}

func (r *RollBackMovement) Do(move Move) {
	r.move = move

	piece := r.game.board[move.From.Row][move.From.Col]
	if piece != nil && piece.Type == King {
		r.game.castleRights[piece.Color] = castling{left: false, right: false}
	}
	if r.isCastling() {
		r.doCastling()
		r.castled = true
		return
	}

	if r.game.board[move.To.Row][move.To.Col] != nil {
		r.capturedPiece = r.game.board[move.To.Row][move.To.Col]
		r.capturedPiece.Captured = true
	}

	r.CheckPromotion(move)

	r.game.board[move.To.Row][move.To.Col] = r.game.board[move.From.Row][move.From.Col]
	r.game.board[move.From.Row][move.From.Col] = nil
	r.game.board[move.To.Row][move.To.Col].Location.Col = move.To.Col
	r.game.board[move.To.Row][move.To.Col].Location.Row = move.To.Row
}

func (r *RollBackMovement) rollBackCastle() {
	color := r.game.board[r.move.To.Row][r.move.To.Col].Color

	r.game.board[r.move.To.Row][4] = r.game.board[r.move.To.Row][r.move.To.Col]
	r.game.board[r.move.To.Row][r.move.To.Col] = nil
	if r.move.To.Col == 2 {
		r.game.board[r.move.To.Row][0] = r.game.board[r.move.From.Row][3]
		r.game.board[r.move.To.Row][3] = nil
		r.game.castleRights[color] = castling{
			left:  true,
			right: r.game.castleRights[color].right,
		}
	} else {
		r.game.board[r.move.To.Row][7] = r.game.board[r.move.To.Row][5]
		r.game.board[r.move.To.Row][5] = nil
		r.game.castleRights[color] = castling{
			left:  r.game.castleRights[color].left,
			right: true,
		}
	}
}

func (r *RollBackMovement) RollBack() {
	if r.rolledBack {
		return
	}

	r.game.castleRights = r.castleRightsBackup
	if r.castled {
		r.rollBackCastle()
		r.rolledBack = true
		return
	}
	if r.capturedPiece != nil {
		r.capturedPiece.Captured = false
	}

	r.game.board[r.move.From.Row][r.move.From.Col] = r.game.board[r.move.To.Row][r.move.To.Col]
	r.game.board[r.move.To.Row][r.move.To.Col] = r.capturedPiece
	r.game.board[r.move.From.Row][r.move.From.Col].Location.Col = r.move.From.Col
	r.game.board[r.move.From.Row][r.move.From.Col].Location.Row = r.move.From.Row

	if r.promoted {
		r.game.board[r.move.From.Row][r.move.From.Col].Type = Pawn
	}
	r.rolledBack = true
}
