let BOARD_SIZE = 650;
if (window.screen.height < 800 || window.screen.width < 1000) {
    BOARD_SIZE = (window.screen.height - 150) < (window.screen.width - 15)? (window.screen.height - 150): (window.screen.width - 15);
}

const SQUIRE_SIZE = BOARD_SIZE / 8;
const WHITE_COLOR = '#e9edcc';
const BLACK_COLOR = '#779954';
const SELECTED_SQUIRE_COLOR = "#f4f67e";
const IMAGE_LINK = "static/img/pieces/";

function getSquireColor(x, y) {
    if (y % 2 === 0) {
        if (x % 2 === 0) {
            return BLACK_COLOR;
        } else {
            return WHITE_COLOR;
        }
    } else {
        if (x % 2 === 0) {
            return WHITE_COLOR;
        } else {
            return BLACK_COLOR;
        }
    }
}

class ChessUI {
    constructor(boardElem, playerElem, opponentElem, gameStatusElem, board) {
        console.log("width: ", window.screen.width);
        console.log("height: ", window.screen.height);
        this.canvas = document.createElement("canvas");
        this.canvas.width = BOARD_SIZE;
        this.canvas.height = BOARD_SIZE
        boardElem.appendChild(this.canvas)

        this.ctx = this.canvas.getContext("2d");
        this.board = board;
        this.pickedPiece = null;
        this.lastClick = 0;

        this.playerElem = playerElem;
        this.opponentElem = opponentElem;
        this.gameStatusElem = gameStatusElem;
    }
    setUp(player, opponent, game) {
        this.game = game;
        this.viewAs = player.color;
        this.setUpBoard();
        this.setUpHandlers();
        this.setUpBar(player, opponent);
    }

    setUpBar(player, opponent) {
        this.playerElem.innerText = player.name + " " + player.color;
        this.opponentElem.innerText = opponent.name + " " + opponent.color;
    }

    render() {
        this.ctx.reset()
        for (let x = 0; x < 8; x++) {
            for (let y = 0; y < 8; y++) {
                this.drawSquire(x, y, getSquireColor(x, y));
                let piece = this.board[y][x]
                if (piece) {
                    this.drawPiece(x, y, piece);
                }
            }
        }
        this.isFirstRender = false
    }
    setUpBoard() {
        this.render();
    }

    changeBackground(x, y, color) {
        this.drawSquire(x, y, color);
        this.drawPiece(x, y, this.board[y][x])
    }
    gameOver(winner) {
        this.gameStatusElem.innerText = `winner is ${winner}`;
    }

    isClicked() {
        if((Date.now() - this.lastClick) < 10) {
            return false
        }
        this.lastClick = Date.now()
        return true;
    }

    setUpHandlers() {
        const pickupPiece = (x, y) => {
            if(this.game.isMyPiece({row: y, col: x})) {
                this.changeBackground(x, y, SELECTED_SQUIRE_COLOR)
                this.pickedPiece = {
                    "piece": this.board[y][x],
                    "location": {x: x, y: y},
                }
            } 
        }
        const dropPiece = (x, y) => {
            this.changeBackground(
                this.pickedPiece.location.x, 
                this.pickedPiece.location.y,
                getSquireColor(
                    this.pickedPiece.location.x,
                    this.pickedPiece.location.y,
                ),
            )
            const played = this.game.play(
                { row: this.pickedPiece.location.y, col: this.pickedPiece.location.x },
                { row: y, col: x },
            )

            this.pickedPiece = null;
            if(!played) {
                console.log("invalid play")
                return;
            }

            this.render();
        }
        this.canvas.addEventListener("click", (event) => {
            if(!this.isClicked()) {
                return;
            }

            const location = this.getClickedCoordination(event)
            if(!location) {
                this.pickedPiece = null;
                return;
            }

            if (this.pickedPiece === null) {
                pickupPiece(location.x, location.y);
            } else {
                dropPiece(location.x, location.y);
            }
        });
    }


    convertToBoardCoordination(x, y) {
        const dy = (this.viewAs == "black")? 0: 7
        return x, Math.abs(y - dy)
    }

    drawSquire(x, y, color) {
        x, y = this.convertToBoardCoordination(x, y)
        this.ctx.fillStyle = color;
        this.ctx.fillRect(x * SQUIRE_SIZE, y * SQUIRE_SIZE, SQUIRE_SIZE, SQUIRE_SIZE);
    }

    loadImage(x, y, piece) {
        let ctx = this.ctx;
        piece.image = new Image();
        piece.image.src = IMAGE_LINK + `${piece.name}-${piece.color}.svg`;
        piece.image.onload = () => {
            piece.image.width = SQUIRE_SIZE;
            piece.image.height = SQUIRE_SIZE;
            piece.image.setAttribute("name", piece.name);
            piece.image.setAttribute("color", piece.color);

            x, y = this.convertToBoardCoordination(x, y)
            ctx.drawImage(piece.image, x * SQUIRE_SIZE, y * SQUIRE_SIZE, piece.image.width, piece.image.height); 
        }
    }
    drawPiece(x, y, piece) {
        this.ctx.fillStyle = piece.color === "white" ? "#fff" : "#000";
        let ctx = this.ctx;
        if(!piece?.image) {
            this.loadImage(x, y, piece);
            return;
        }
        if(piece.image.getAttribute("name") !== piece.name) {
            // Pawn promotion
            this.loadImage(x, y, piece);
            return;
        } 

        x, y = this.convertToBoardCoordination(x, y)
        ctx.drawImage(piece.image, x * SQUIRE_SIZE, y * SQUIRE_SIZE, piece.image.width, piece.image.height); 
    }

    getClickedCoordination(event) {
        const rect = this.canvas.getBoundingClientRect();
        let x = Math.floor((event.clientX - rect.x) / SQUIRE_SIZE)
        let y = Math.floor((event.clientY - rect.y) / SQUIRE_SIZE)
        if(x < 0 || x > 7 || y < 0 || y > 7) {
            return null;
        }

        x, y = this.convertToBoardCoordination(x, y)
        return {x: x, y: y}
    }
}