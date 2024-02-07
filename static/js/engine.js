const WHITE = "white";
const BLACK = "black";

const KING = "king";
const QUEEN = "queen";
const BISHOP = "bishop";
const PAWN = "pawn";
const KNIGHT = "knight";
const ROOK = "rook";

const oppositeColor = (color) => (color === "black") ? "white": "black";

class BoardLocation {
    constructor(row, col) {
        this.valid = row < 8 && row >= 0 && col < 8 && col >= 0
        this.row = row
        this.col = col
    }

    compare(location) {
        return (this.row === location.row) && (this.col === location.col)
    }

    toString() {
        return `${this.row},${this.col}`
    }
}

function locationFromString(str) {
    const strList = str.split(',')
    const row = Number.parseInt(strList[0])
    const col = Number.parseInt(strList[1])

    return new BoardLocation(row, col)
}

class Rollback {
    constructor(engine) {
        this.engine = engine;
        this.capturedPiece = null;
        this.castled = false;
        this.promotedPawn = false;
    }

    do(src, dst) {
        this.move = {
            'from': src,
            'to': dst,
        };

        if (this.isCastling(src, dst)) {
            this.doCastling(src, dst);
            this.castled = true;
            return;
        }

        if (this.engine.board[dst.row][dst.col]) {
            this.capturedPiece = this.engine.board[dst.row][dst.col];
        }
        this.checkPawnPromotion(src, dst);

        this.engine.board[dst.row][dst.col] = this.engine.board[src.row][src.col]
        this.engine.board[src.row][src.col] = null
    }

    rollback() {
        if(this.castled) {
            return this.rollbackCastling();
        }
        let src = this.move.from;
        let dst = this.move.to;
        this.engine.board[src.row][src.col] = this.engine.board[dst.row][dst.col];
        this.engine.board[dst.row][dst.col] = this.capturedPiece;

        if (this.promotedPawn) {
            this.engine.board[src.row][src.col].name = PAWN; 
        }
    }

    checkPawnPromotion(src, dst) {
        let piece = this.engine.board[src.row][src.col]
        if(piece.name === PAWN) {
            if (piece.color === BLACK && dst.row === 0) {
                this.engine.board[src.row][src.col].name = QUEEN; 
                this.promotedPawn = true;
            } else if (dst.row === 7) {
                this.engine.board[src.row][src.col].name = QUEEN;
                this.promotedPawn = true;
            }
        }
    }

    isCastling(src, dst) {
        return this.engine.board[src.row][src.col].name === KING && Math.abs(src.col - dst.col) === 2
    }

    doCastling(src, dst) {
        if(src.col > dst.col) {
            this.doLeftCasting(src, dst);
        } else {
            this.doRightCastling(src, dst);
        }
    }
    
    doLeftCasting(src, dst) {
        this.engine.board[dst.row][dst.col] = this.engine.board[src.row][src.col];
        this.engine.board[src.row][src.col] = null;
        this.engine.board[dst.row][dst.col + 1] = this.engine.board[src.row][0];
        this.engine.board[src.row][0] = null;
    }
    doRightCastling(src, dst) {
        this.engine.board[dst.row][dst.col] = this.engine.board[src.row][src.col];
        this.engine.board[src.row][src.col] = null;
        this.engine.board[dst.row][dst.col - 1] = this.engine.board[src.row][7];
        this.engine.board[src.row][7] = null;
    }

    rollbackCastling() {
        let src = this.move.from;
        let dst = this.move.to;
        if(src.col > dst.col) {
            this.rollbackLeftCastling(src, dst);
        } else {
            this.rollbackRightCastling(src, dst);
        }
    }

    rollbackLeftCastling(src, dst) {
        this.engine.board[src.row][src.col] = this.engine.board[dst.row][dst.col];
        this.engine.board[dst.row][dst.col] = null;
        this.engine.board[src.row][0] = this.engine.board[dst.row][dst.col + 1];
        this.engine.board[dst.row][dst.col + 1] = null;
    }
    rollbackRightCastling(src, dst) {
        this.engine.board[src.row][src.col] = this.engine.board[dst.row][dst.col];
        this.engine.board[dst.row][dst.col] = null;
        this.engine.board[src.row][7] = this.engine.board[dst.row][dst.col - 1];
        this.engine.board[dst.row][dst.col - 1] = null;
    }
}
class ChessEngine {
    constructor(board, color) {
        this.board = board;
        this.myColor = color 
        this.turn = WHITE; 
        this.winner = null;
        this.castlingRight = {
            white: {
                left: true,
                right: true,
            },
            black: {
                left: true,
                right: true,
            },
        }
        this.findAllPossibleMoves()
    }

    isMyPiece(loc) {
        if(this.board[loc.row][loc.col]) {
            if(this.board[loc.row][loc.col].color === this.myColor) {
                return true;
            }
        }
        return false;
    }
    getKing(color) {
        for(let i = 0; i < 8; i++){ 
            for(let j = 0; j < 8; j++) {
                if(this.board[i][j] && this.board[i][j].color === color && this.board[i][j].name === KING) {
                    return new BoardLocation(i, j)
                }
            }
        }
    }

    switchTurn() {
        if (this.turn === WHITE) {
            this.turn = BLACK
        } else {
            this.turn = WHITE
        }
        this.findAllPossibleMoves()
    }

    isOpponentPiece(piece, dst) {
        return (piece.color === WHITE && this.board[dst.row][dst.col]?.color === BLACK) ||
            (piece.color === BLACK && this.board[dst.row][dst.col]?.color === WHITE)
    }

    /* A square is occupiable if it's a valid location, empty or occupied by opponent piece. */
    isOccupiable(piece, dst) {
        return dst.valid && (this.isOpponentPiece(piece, dst) || !this.board[dst.row][dst.col])
    }

    getPiecePossibleMoves(location) {
        if (!(location.toString() in this.possibleMoves)) {
            return []
        }
        return this.possibleMoves[location.toString()]
    }

    isCheckable(kingLocation) {
        const kingPiece = this.board[kingLocation.row][kingLocation.col];
        for(let i = 0; i < 8; i++) {
            for(let j = 0; j < 8; j++) {
                if(this.board[i][j] && this.isOpponentPiece(kingPiece, new BoardLocation(i, j))) {
                    if(this.validMove(new BoardLocation(i, j), kingLocation)) {
                        return true
                    }
                }
            }
        }
    }
    // Entry
    movePiece(from, to) {
        let src = new BoardLocation(from.row, from.col);
        let dst = new BoardLocation(to.row, to.col);

        if (this.turn !== this.board[src.row][src.col]?.color) {
            return false
        }
        if (this.winner) {
            throw "game is over";
        }
        if (this.possibleMoves[src.toString()].find(loc => loc.compare(dst))) {
            if (this.move(src, dst)) {
                this.switchTurn()
                if(!this.hasMove()) {
                    console.log("no move remains"); 
                    this.winner = oppositeColor(this.turn); 
                }
                return true
            }
        }
        return false
    }

    haveCastleRight(color) {
        return this.castlingRight[color].left || this.castlingRight[color].right;
    }
    updateCastleRights(src, movedPiece) {
        if(movedPiece.name === KING)  {
            this.castlingRight[movedPiece.color].left = false
            this.castlingRight[movedPiece.color].right = false
        } else if (movedPiece.name === ROOK) {
            if(src.col === 0) {
                // First time moving left side rook
                this.castlingRight[movedPiece.color].left = false; 
            } else if(src.col === 7) {
                // First time moving right side rook
                this.castlingRight[movedPiece.color].right = false; 
            }
        }
    }
    move(src, dst) {
        const rollback = new Rollback(this);
        rollback.do(src, dst);

        let movedPiece = this.board[dst.row][dst.col];
        if (this.haveCastleRight(movedPiece.color)) {
            this.updateCastleRights(src, movedPiece);
        } 
        
        return true
    }
    /* Reevaluates pieces and their possible moves */
    findAllPossibleMoves() {
        this.possibleMoves = {}
        for (let i = 0; i < 8; i++) {
            for (let j = 0; j < 8; j++) {
                if (this.board[i][j] && this.board[i][j].color === this.turn) {
                    const location = new BoardLocation(i, j)
                    this.possibleMoves[location.toString()] = this.findPiecePossibleMoves(location)
                }
            }
        }

        const kingLocation = this.getKing(this.turn)
        const newPossibleMoves = {}
        
        Object.keys(this.possibleMoves).forEach(srcKey => {
            const src = locationFromString(srcKey)
            newPossibleMoves[srcKey] = []

            this.possibleMoves[srcKey].forEach(dst => {
                const rollback = new Rollback(this);
                rollback.do(src, dst)
                if(this.board[dst.row][dst.col].name === KING) {
                    if(!this.isCheckable(new BoardLocation(dst.row, dst.col))) {
                        newPossibleMoves[srcKey].push(dst)
                    }
                }
                else if(!this.isCheckable(kingLocation)){ 
                    newPossibleMoves[srcKey].push(dst)
                }
                rollback.rollback()
            });
        });

        this.possibleMoves = newPossibleMoves 
    }


    hasMove() {
        let moves = false;
        Object.keys(this.possibleMoves).forEach(srcKey => {
            if(this.possibleMoves[srcKey].length > 0) {
                moves = true;
                return;
            }
        });
        return moves;
    }
    findPiecePossibleMoves(location) {
        const piece = this.board[location.row][location.col];
        if (!piece) {
            return []
        }
        switch (piece.name) {
            case PAWN:
                return this.findPawnPossibleMoves(piece, location)
            case ROOK:
                return this.findRookPossibleMoves(piece, location)
            case BISHOP:
                return this.findBishopPossibleMoves(piece, location)
            case KNIGHT:
                return this.findKnightPossibleMoves(piece, location)
            case QUEEN:
                return this.findQueenPossibleMoves(piece, location)
            case KING:
                return this.findKingPossibleMoves(piece, location)
            default:
                return []
        }
    }

    findKingPossibleMoves(piece, location) {
        let possibleMoves = []
        const destinations = [
            new BoardLocation(location.row + 1, location.col), // Down
            new BoardLocation(location.row - 1, location.col), // Up
            new BoardLocation(location.row, location.col + 1), // Right 
            new BoardLocation(location.row, location.col - 1), // Left

            new BoardLocation(location.row + 1, location.col + 1), // Down-right
            new BoardLocation(location.row + 1, location.col - 1), // Down-left
            new BoardLocation(location.row - 1, location.col + 1), // Up-right
            new BoardLocation(location.row - 1, location.col - 1), // Up-left
        ]
        destinations.forEach((dst) => {
            if (this.isOccupiable(piece, dst)) {
                possibleMoves.push(dst)
            }
        })

        if(this.castlingRight[piece.color].left) {
            if(!this.board[location.row][3] && !this.board[location.row][2] && !this.board[location.row][1]) {
                possibleMoves.push(new BoardLocation(location.row, location.col - 2))
            }
        }
        if(this.castlingRight[piece.color].right) {
            if(!this.board[location.row][5] && !this.board[location.row][6]) {
                possibleMoves.push(new BoardLocation(location.row, location.col + 2))
            }
        }
        return possibleMoves
    }

    findKnightPossibleMoves(piece, location) {
        let possibleMoves = []

        const destinations = [
            new BoardLocation(location.row + 2, location.col + 1),
            new BoardLocation(location.row + 2, location.col - 1),
            new BoardLocation(location.row - 2, location.col + 1),
            new BoardLocation(location.row - 2, location.col - 1),
            new BoardLocation(location.row + 1, location.col + 2),
            new BoardLocation(location.row + 1, location.col - 2),
            new BoardLocation(location.row - 1, location.col + 2),
            new BoardLocation(location.row - 1, location.col - 2),
        ]

        destinations.forEach((dst) => {
            if (this.isOccupiable(piece, dst)) {
                possibleMoves.push(dst)
            }
        })

        return possibleMoves
    }

    findPawnPossibleMoves(piece, location) {
        let possibleMoves = []

        if (piece.color === WHITE) {
            if (!this.board[location.row + 1][location.col]) {
                possibleMoves.push(new BoardLocation(location.row + 1, location.col))
                // First move
                if (location.row === 1 && !this.board[location.row + 2][location.col]) {
                    possibleMoves.push(new BoardLocation(location.row + 2, location.col))
                }
            }

            let destinations = [
                new BoardLocation(location.row + 1, location.col + 1),
                new BoardLocation(location.row + 1, location.col - 1)
            ]
            destinations.forEach((dst) => {
                if (this.isOpponentPiece(piece, dst)) {
                    possibleMoves.push(dst)
                }
            })
        } else {
            if (!this.board[location.row - 1][location.col]) {
                possibleMoves.push(new BoardLocation(location.row - 1, location.col))
                // First move
                if (location.row === 6 && !this.board[location.row - 2][location.col]) {
                    possibleMoves.push(new BoardLocation(location.row - 2, location.col))
                }
            }
            let destinations = [
                new BoardLocation(location.row - 1, location.col + 1),
                new BoardLocation(location.row - 1, location.col - 1)
            ]
            destinations.forEach((dst) => {
                if (this.isOpponentPiece(piece, dst)) {
                    possibleMoves.push(dst)
                }
            })
        }
        return possibleMoves
    }

    findBishopPossibleMoves(piece, location) {
        let possibleMoves = []

        // checks for possible actions
        const shortcutFunc = (loc) => {
            // return false means don't check further squares in that direction
            if (!this.isOccupiable(piece, loc)) {
                return false
            }
            if (this.isOpponentPiece(piece, loc)) {
                possibleMoves.push(loc)
                return false
            }
            possibleMoves.push(loc)
            return true
        }

        let j = location.col + 1
        let i = location.row + 1

        // Down-right
        for (i = location.row + 1; i < 8; i++) {
            if (!shortcutFunc(new BoardLocation(i, j))) break
            j++
        }
        j = location.col - 1;
        for (i = location.row + 1; i < 8; i++) {
            // Down-left
            if (!shortcutFunc(new BoardLocation(i, j))) break
            j--
        }
        j = location.col + 1;
        for (i = location.row - 1; i >= 0; i--) {
            // Up-right
            if (!shortcutFunc(new BoardLocation(i, j))) break
            j++
        }
        j = location.col - 1;
        for (i = location.row - 1; i >= 0; i--) {
            // Up-left
            if (!shortcutFunc(new BoardLocation(i, j))) break
            j--;
        }
        return possibleMoves
    }
    findRookPossibleMoves(piece, location) {
        let possibleMoves = []
        let i, j = 0

        // checks for possible actions
        const shortcutFunc = (loc) => {
            // return false means don't check further squares in that direction
            if (!this.isOccupiable(piece, loc)) {
                return false
            }
            if (this.isOpponentPiece(piece, loc)) {
                possibleMoves.push(loc)
                return false
            }
            possibleMoves.push(loc)
            return true
        }

        for (i = location.row + 1; i < 8; i++) {
            if (!shortcutFunc(new BoardLocation(i, location.col))) break
        }
        for (i = location.row - 1; i >= 0; i--) {
            if (!shortcutFunc(new BoardLocation(i, location.col))) break
        }
        for (j = location.col + 1; j < 8; j++) {
            if (!shortcutFunc(new BoardLocation(location.row, j))) break
        }
        for (j = location.col - 1; j >= 0; j--) {
            if (!shortcutFunc(new BoardLocation(location.row, j))) break
        }
        return possibleMoves
    }

    findQueenPossibleMoves(piece, location) {
        return [
            ...this.findBishopPossibleMoves(piece, location),
            ...this.findRookPossibleMoves(piece, location)
        ]
    }

    validMove(src, dst) {
        switch(this.board[src.row][src.col].name) {
            case QUEEN:
                return this.validQueenMove(src, dst) 
            case BISHOP:
                return this.validBishopMove(src, dst)
            case ROOK:
                return this.validRookMove(src, dst)
            case PAWN:
                return this.validPawnMove(src, dst)
            case KNIGHT:
                return this.validKnightMove(src, dst)
            case KING:
                return this.validKingMove(src, dst)

        }
        return false
    }
    validQueenMove(src, dst) {
        return (this.validBishopMove(src, dst) || this.validRookMove(src, dst))
    }

    validKingMove(src, dst) {
        return (src.col == dst.col &&
            Math.abs(src.row - dst.row) === 1) ||
            (src.row === dst.row &&
                Math.abs(src.col - dst.col) === 1) ||
            (src.col === dst.col+1) && (src.row === dst.row+1) ||
            (src.col === dst.col+1) && (src.row === dst.row-1) ||
            (src.col === dst.col-1) && (src.row === dst.row+1) ||
            (src.col === dst.col-1) && (src.row === dst.row-1)
    }
    validKnightMove(src, dst) {
        return ((dst.row === src.row + 2 && Math.abs(dst.col - src.col) === 1) ||
            (dst.row === src.row - 2 && Math.abs(dst.col - src.col) === 1) ||
            (dst.col === src.col + 2 && Math.abs(dst.row - src.row) === 1) ||
            (dst.col === src.col - 2 && Math.abs(dst.row - src.row) === 1))
    }
    validBishopMove(src, dst) {
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
        } else if ((dst.col - src.col) === (dst.row - src.row) * (-1) && (dst.col - src.col) < 0) {
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
    validPawnMove(src, dst) {
        
        if (this.board[src.row][src.col].color === BLACK) {
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

    validRookMove(src, dst) {
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

function mapChar2Piece(char) {
    switch(char) {
        case 'r':
            return ROOK
        case 'n':
            return KNIGHT
        case 'b':
            return BISHOP 
        case 'q':
            return QUEEN
        case 'k':
            return KING
        case 'p':
            return PAWN
    }
}
function mapChar2Color(char) {
    switch(char) {
        case 'w':
            return WHITE
        case 'b':
            return BLACK
    }
}

function NewChessBoard() {
    const map = [
        ['rw', 'nw', 'bw', 'qw', 'kw', 'bw', 'nw', 'rw'],
        ['pw', 'pw', 'pw', 'pw', 'pw', 'pw', 'pw', 'pw'],
        ['--', '--', '--', '--', '--', '--', '--', '--'],
        ['--', '--', '--', '--', '--', '--', '--', '--'],
        ['--', '--', '--', '--', '--', '--', '--', '--'],
        ['--', '--', '--', '--', '--', '--', '--', '--'],
        ['pb', 'pb', 'pb', 'pb', 'pb', 'pb', 'pb', 'pb'],
        ['rb', 'nb', 'bb', 'qb', 'kb', 'bb', 'nb', 'rb'],
    ]

    const board = []
    for(let i = 0; i < 8; i++) {
        board.push([])
        for(let j = 0; j < 8; j++) {
            if(map[i][j] === '--') continue;
            board[i].push({
                name: mapChar2Piece(map[i][j][0]),
                color: mapChar2Color(map[i][j][1]),
                location: new BoardLocation(i, j),
            });
        }
    }
    return board;
}