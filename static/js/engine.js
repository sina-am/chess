const WHITE = "white";
const BLACK = "black";

const KING = "king";
const QUEEN = "queen";
const BISHOP = "bishop";
const PAWN = "pawn";
const KNIGHT = "knight";
const ROOK = "rook";

const rowNumbs = [8, 7, 6, 5, 4, 3, 2, 1]
const colNumbs = ["a", "b", "c", "d", "e", "f", "g", "h"]

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

class DummyEngine {
    constructor(board) {
        this.board = board;
    }
    move(from, to) {
        this.board[from.row][from.col] = this.board[to.row][to.col]; // for test
    }
    isPieceMine(loc) {
        if(this.board[loc.row][loc.col]) {
            if(this.board[loc.row][loc.col].color === "white") {
                return true;
            }
        }
        return false;
    }
}

function getEmptyBoard() {
    const board = [];
    for(let i = 0; i < 8; i++){
        board.push([]);
        for(let j = 0; j < 8; j++) {
            board[i].push(null);
        }
    }
    board[0][0] = chessSetup[0][4]
    board[0][6] = chessSetup[0][3]
    board[7][7] = chessSetup[7][4]
    return board;
}

class ChessEngine {
    constructor(board, color) {
        this.board = board;
        this.myColor = color 
        this.turn = WHITE; 
        this.winner = null;
        this.castlingRight = new Map([[WHITE, true], [BLACK, true]])
        this.capturedPieces = new Map([[WHITE, []], [BLACK, []]])

        this.lastMove = {
            "piece": null,
            "from": null,
            "to": null,
            "capturedPiece": null
        }

        this.possibleMoves = {}
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

    redo() {
        this.board[this.lastMove.from.row][this.lastMove.from.col] = this.lastMove.piece
        if (this.lastMove.capturedPiece) {
            this.board[this.lastMove.to.row][this.lastMove.to.col] = this.lastMove.capturedPiece
        } else {
            this.board[this.lastMove.to.row][this.lastMove.to.col] = null 
        }

        this.switchTurn(); 
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
                let capturedPiece = null

                if (this.isOpponentPiece(this.board[src.row][src.col], dst)) {
                    capturedPiece = this.board[dst.row][dst.col];
                }
                this.board[dst.row][dst.col] = this.board[src.row][src.col]
                this.board[src.row][src.col] = null
            
                if(this.board[dst.row][dst.col].name === KING) {
                    if(!this.isCheckable(new BoardLocation(dst.row, dst.col))) {
                        newPossibleMoves[srcKey].push(dst)
                    }
                }
                else if(!this.isCheckable(kingLocation)){ 
                    newPossibleMoves[srcKey].push(dst)
                }

                this.board[src.row][src.col] = this.board[dst.row][dst.col]
                if (capturedPiece) {
                    this.board[dst.row][dst.col] = capturedPiece 
                } else {
                    this.board[dst.row][dst.col] = null 
                }
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

        if (this.castlingRight.get(piece.color)) {
            let leftCastling = true;
            for (let j = location.col - 1; j > 0; j--) {
                if (this.board[location.row][j]) {
                    leftCastling = false
                    break
                }
            }
            if (leftCastling) {
                if (piece.color === BLACK) {
                    possibleMoves.push(new BoardLocation(location.row, location.col - 3))
                } else {
                    possibleMoves.push(new BoardLocation(location.row, location.col - 2))
                }
            }

            let rightCastling = true;
            for (let j = location.col + 1; j < 7; j++) {
                if (this.board[location.row][j]) {
                    rightCastling = false
                    break
                }
            }
            if (rightCastling) {
                if (piece.color === WHITE) {
                    possibleMoves.push(new BoardLocation(location.row, location.col + 3))
                } else {
                    possibleMoves.push(new BoardLocation(location.row, location.col + 2))
                }
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

    move(src, dst) {
        this.lastMove.piece = this.board[src.row][src.col]
        this.lastMove.from = new BoardLocation(src.row, src.col)
        this.lastMove.to = new BoardLocation(dst.row, dst.col)

        if (this.isOpponentPiece(this.board[src.row][src.col], dst)) {
            // TODO: make this clean
            this.lastMove.capturedPiece = this.board[dst.row][dst.col]
            this.capturedPieces.set(this.board[dst.row][dst.col].color, [this.board[dst.row][dst.col], ...this.capturedPieces.get(this.board[dst.row][dst.col].color)])
            this.board[dst.row][dst.col] = null;
        }
        this.board[dst.row][dst.col] = this.board[src.row][src.col]
        this.board[src.row][src.col] = null

        if(this.board[dst.row][dst.col].name === KING)  {
            this.castlingRight.set(this.board[dst.row][dst.col].color, false);
        }
        
        return true
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

const chessSetup = [
    [
        {
            name: 'rook', 
            color: 'white',
        }, 
        {
            name: 'knight',
            color: 'white',
        },
        {
            name: 'bishop',
            color: 'white',
        },
        {
            name: 'queen', 
            color: 'white',
        }, {
            name: 'king',
            color: 'white',
        },
        {
            name: 'bishop',
            color: 'white',
        },
        {
            name: 'knight',
            color: 'white',
        },
        {
            name: 'rook', 
            color: 'white',
        }
    ],
    [
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
        {
            name: 'pawn', 
            color: 'white',
        },
    ],
    [null, null, null, null, null, null, null, null], 
    [null, null, null, null, null, null, null, null],
    [null, null, null, null, null, null, null, null], 
    [null, null, null, null, null, null, null, null],
    [
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
        {
            name: 'pawn', 
            color: 'black',
        },
    ],  
   [
        {
            name: 'rook', 
            color: 'black',
        }, 
        {
            name: 'knight',
            color: 'black',
        },
        {
            name: 'bishop',
            color: 'black',
        },
        {
            name: 'queen', 
            color: 'black',
        }, {
            name: 'king',
            color: 'black',
        },
        {
            name: 'bishop',
            color: 'black',
        },
        {
            name: 'knight',
            color: 'black',
        },
        {
            name: 'rook', 
            color: 'black',
        }
    ], 
    
]