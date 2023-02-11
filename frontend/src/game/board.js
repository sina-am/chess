import './piece'
import { BishopBlackPiece, BishopWhitePiece, KingBlackPiece, KingWhitePiece, KnightBlackPiece, KnightWhitePiece, PawnBlackPiece, PawnWhitePiece, QueenBlackPiece, QueenWhitePiece, RookBlackPiece, RookWhitePiece } from './piece';

export function getStandrardBoard() {
    return [
        [
            RookBlackPiece,
            KnightBlackPiece,
            BishopBlackPiece,
            QueenBlackPiece,
            KingBlackPiece,
            BishopBlackPiece,
            KnightBlackPiece,
            RookBlackPiece,
        ],
        [
            PawnBlackPiece, PawnBlackPiece,
            PawnBlackPiece, PawnBlackPiece,
            PawnBlackPiece, PawnBlackPiece,
            PawnBlackPiece, PawnBlackPiece
        ],
        [
            null, null, null, null, null, null, null, null
        ],
        [
            null, null, null, null, null, null, null, null
        ],
        [
            null, null, null, null, null, null, null, null
        ],
        [
            null, null, null, null, null, null, null, null
        ],
        [
            PawnWhitePiece, PawnWhitePiece,
            PawnWhitePiece, PawnWhitePiece,
            PawnWhitePiece, PawnWhitePiece,
            PawnWhitePiece, PawnWhitePiece
        ],
        [
            RookWhitePiece,
            KnightWhitePiece,
            BishopWhitePiece,
            KingWhitePiece,
            QueenWhitePiece,
            BishopWhitePiece,
            KnightWhitePiece,
            RookWhitePiece
        ],
    ]
}

export class Game {
    constructor(board, player1, player2) {
        this.board = board;
    }

    movePiece(src, dst) {
        const piece = this.board[src.row][src.col];
        
        switch (piece.type) {
            case "pawn":
                return this.validPawnMove(piece, src, dst) ? this.move(src, dst) : false
            case "rook":
                return this.validRookMove(piece, src, dst) ? this.move(src, dst) : false
            case "bishop":
                return this.validBishopMove(piece, src, dst) ? this.move(src, dst) : false
            case "knight":
                return this.validKnightMove(piece, src, dst) ? this.move(src, dst) : false
            case "queen":
                return this.validQueentMove(piece, src, dst) ? this.move(src, dst) : false
            default:
                return
        }
    }

    move(src, dst) {
        if (this.board[dst.row][dst.col]) {
            this.board[dst.row][dst.col].isDead = true;

        }
        this.board[dst.row][dst.col] = this.board[src.row][src.col]
        this.board[src.row][src.col] = null
        return true
    }

    validQueentMove(piece, src, dst) {
        return (this.validBishopMove(piece, src, dst) || this.validRookMove(piece, src, dst))
    }

    validKnightMove(piece, src, dst) {
        return ((dst.row === src.row+2 && Math.abs(dst.col-src.col) === 1) ||
            (dst.row === src.row-2 && Math.abs(dst.col-src.col) === 1) ||
            (dst.col === src.col+2 && Math.abs(dst.row-src.row) === 1) ||
            (dst.col === src.col-2 && Math.abs(dst.row-src.row) === 1))
    }
    validBishopMove(piece, src, dst) {
        if ((dst.col - src.col) === (dst.row - src.row) && (dst.col - src.col) > 0) {
            // Move down-right
            let colStep = src.col + 1
            for (let rowStep = src.row + 1; rowStep < dst.row; rowStep++) {
                if (this.board[rowStep][colStep]) {
                    return false
                }
                colStep++
            }
            return true
        } else if ((dst.col - src.col) === (dst.row - src.row) && (dst.col - src.col) < 0) {
            // Move up-left
            let colStep = src.col - 1
            for (let rowStep = src.row - 1; rowStep > dst.row; rowStep--) {
                if (this.board[rowStep][colStep]) {
                    return false
                }
                colStep--
            }
            return true
        } else if ((dst.col - src.col) === (dst.row - src.row)*(-1) && (dst.col - src.col) < 0) {
            // Move down-left
            let rowStep = src.row + 1
            for (let colStep = src.col - 1; colStep > dst.col; colStep--) {
                if (this.board[rowStep][colStep]) {
                    return false
                }
                rowStep++
            }
            return true
        } else if ((dst.col - src.col) === (dst.row - src.row) * (-1) && (dst.col - src.col) > 0) {
            // Move up-right
            let rowStep = src.row - 1
            for (let colStep = src.col + 1; colStep < dst.col; colStep++) {
                if (this.board[rowStep][colStep]) {
                    return false
                }
                rowStep--
            }
            return true
        }
        return false
    }
    validPawnMove(piece, src, dst) {
        if (piece.color === "black") {
            if (src.col === dst.col && dst.row === src.row + 1 && !this.board[dst.row][dst.col]) {
                return true
            } else if ((dst.col === src.col + 1 || dst.col === src.col - 1) && dst.row === src.row + 1 && this.board[dst.row][dst.col]) {
                return true
            }
        } else {
            if (src.col === dst.col && dst.row === src.row - 1 && !this.board[dst.row][dst.col]) {
                return true
            } else if ((dst.col === src.col - 1 || dst.col === src.col + 1) && dst.row === src.row - 1 && this.board[dst.row][dst.col]) {
                return true
            }
        }
        return false
    }

    validRookMove(piece, src, dst) {
        if (src.row === dst.row) {
            if (src.col > dst.col) {
                for (let i = src.col - 1; i > dst.col; i--) {
                    if (this.board[src.row][i]) {
                        return false
                    }
                }
                return true
            } else if (src.col < dst.col) {
                for (let i = src.col + 1; i < dst.col; i++) {
                    if (this.board[src.row][i]) {
                        return false
                    }
                }
                return true
            } else {
                return false;
            }
        } else if (src.col === dst.col) {
            if (src.row > dst.row) {
                for (let i = src.row - 1; i > dst.row; i--) {
                    if (this.board[i][src.col]) {
                        return false
                    }
                }
                return true
            } else if (src.row < dst.row) {
                for (let i = src.row + 1; i < dst.row; i++) {
                    if (this.board[i][src.col]) {
                        return false
                    }
                }
                return true
            }
        } else {
            return false
        }
    }
}
