package engine

// import (
// 	"math"

// 	"github.com/sina-am/chess/types"
// )

// func (b types.Board) isValidKingMove(src, dst types.Location) bool {
// 	return (src.Col == dst.Col &&
// 		math.Abs(float64(src.Row)-float64(dst.Row)) == 1) ||
// 		(src.Row == dst.Row &&
// 			math.Abs(float64(src.Col)-float64(dst.Col)) == 1)
// }

// func (b types.Board) isValidRookMove(src, dst types.Location) bool {
// 	if src.Col == dst.Col && src.Row < dst.Row {
// 		// Move up
// 		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
// 			if b[rowStep][dst.Col] != nil {
// 				return false
// 			}
// 		}
// 		return true
// 	} else if src.Col == dst.Col && src.Row > dst.Row {
// 		// Move down
// 		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
// 			if b[rowStep][dst.Col] != nil {
// 				return false
// 			}
// 		}
// 	} else if src.Row == dst.Row && src.Col < dst.Col {
// 		// Move right
// 		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
// 			if b[src.Row][colStep] != nil {
// 				return false
// 			}
// 		}
// 		return true
// 	} else if src.Row == dst.Row && src.Col > dst.Col {
// 		// Move left
// 		for colStep := src.Col - 1; colStep > dst.Col; colStep++ {
// 			if b[src.Col][colStep] != nil {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func (b types.Board) isValidBishopMove(src, dst types.Location) bool {
// 	if (dst.Col-src.Col) == (dst.Row-src.Row) && (dst.Col-src.Col) > 0 {
// 		// Move up-right
// 		colStep := src.Col + 1
// 		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
// 			if b[rowStep][colStep] != nil {
// 				return false
// 			}
// 			colStep++
// 		}
// 		return true
// 	} else if (dst.Col-src.Col) == (dst.Row-src.Row) && (dst.Col-src.Col) < 0 {
// 		// Move down-left
// 		colStep := src.Col - 1
// 		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
// 			if b[rowStep][colStep] != nil {
// 				return false
// 			}
// 			colStep--
// 		}
// 	} else if (dst.Col-src.Col) == (dst.Row-src.Row)*-1 && (dst.Col-src.Col) < 0 {
// 		// Move up-left
// 		rowStep := src.Row + 1
// 		for colStep := src.Col - 1; colStep > dst.Col; colStep-- {
// 			if b[rowStep][colStep] != nil {
// 				return false
// 			}
// 			rowStep++
// 		}
// 		return true
// 	} else if (dst.Col-src.Col) == (dst.Row-src.Row)*-1 && (dst.Col-src.Col) > 0 {
// 		// Move down-right
// 		rowStep := src.Row - 1
// 		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
// 			if b[src.Col][colStep] != nil {
// 				return false
// 			}
// 			rowStep--
// 		}
// 		return true
// 	}
// 	return false
// }

// func (b types.Board) isValidKnightMove(src, dst types.Location) bool {
// 	return dst.Row == src.Row+2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
// 		dst.Row == src.Row-2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
// 		dst.Col == src.Col+2 && math.Abs(float64(dst.Row-src.Row)) == 1 ||
// 		dst.Col == src.Col-2 && math.Abs(float64(dst.Row-src.Row)) == 1
// }

// func (b types.Board) isValidPawnMove(piece *types.Piece, src, dst types.Location) bool {
// 	if piece.Color == types.White {
// 		if src.Col == dst.Col && dst.Row == src.Row+1 && b[dst.Row][dst.Col] == nil {
// 			return true
// 		} else if (dst.Col == src.Col+1 || dst.Col == src.Col-1) && dst.Row == src.Row+1 && b[dst.Row][dst.Col] != nil {
// 			return true
// 		}
// 	} else {
// 		if src.Col == dst.Col && dst.Row == src.Row-1 && b[dst.Row][dst.Col] == nil {
// 			return true
// 		} else if (dst.Col == src.Col-1 || dst.Col == src.Col+1) && dst.Row == src.Row-1 && b[dst.Row][dst.Col] != nil {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (b types.Board) isValidMove(piece *types.Piece, src types.Location, dst types.Location) bool {
// 	switch piece.Type {
// 	case types.King:
// 		return g.isValidKingMove(src, dst)
// 	case types.Rook:
// 		return g.isValidRookMove(src, dst)
// 	case types.Pawn:
// 		return g.isValidPawnMove(piece, src, dst)
// 	case types.Bishop:
// 		return g.isValidBishopMove(src, dst)
// 	case types.Queen:
// 		return g.isValidBishopMove(src, dst) || g.isValidRookMove(src, dst)
// 	case types.Knight:
// 		return g.isValidKnightMove(src, dst)
// 	default:
// 		return false
// 	}
// }
