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
        this.visitedSquire = null;
        this.last_click = 0;

        this.playerElem = playerElem;
        this.opponentElem = opponentElem;
        this.gameStatusElem = gameStatusElem
    }
    async setUp(player, opponent, game) {
        this.game = game;
        await this.setUpBoard();
        this.setUpHandlers();
        await this.setUpBar(player, opponent);
    }

    async setUpBar(player, opponent) {
        this.playerElem.innerText = player.name + " " + player.color;
        this.opponentElem.innerText = opponent.name + " " + opponent.color;
    }

    async render() {
        this.ctx.reset()
        for (let x = 0; x < 8; x++) {
            for (let y = 0; y < 8; y++) {
                this.drawSquire(x, y, getSquireColor(x, y));
                let piece = this.board[y][x]
                if (piece) {
                    await this.drawPiece(x, y, piece);
                }
            }
        }
        this.isFirstRender = false
    }
    async setUpBoard() {
        await this.render();
    }

    drawSquire(x, y, color) {
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

            ctx.drawImage(piece.image, x * SQUIRE_SIZE, y * SQUIRE_SIZE, piece.image.width, piece.image.height); 
        }
    }
    async drawPiece(x, y, piece) {
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
        ctx.drawImage(piece.image, x * SQUIRE_SIZE, y * SQUIRE_SIZE, piece.image.width, piece.image.height); 
    }
    async changeBackground(x, y, color) {
        this.drawSquire(x, y, color);
        await this.drawPiece(x, y, this.board[y][x])
    }
    gameOver(winner) {
        this.gameStatusElem.innerText = `winner is ${winner}`;
    }

    isClicked() {
        if((Date.now() - this.last_click) < 10) {
            return false
        }
        this.last_click = Date.now()
        return true;
    }

    setUpHandlers() {
        this.canvas.addEventListener("click", async (event) => {
            if(!this.isClicked()) {
                return;
            }

            if (this.pickedPiece === null) {
                const rect = this.canvas.getBoundingClientRect();
                let x = Math.floor((event.clientX - rect.x) / SQUIRE_SIZE)
                let y = Math.floor((event.clientY - rect.y) / SQUIRE_SIZE)
                if(x < 0 || x > 7 || y < 0 || y > 7) {
                    return ;
                }
                if(this.game.isMyPiece({row: y, col: x})) {
                    await this.changeBackground(x, y, SELECTED_SQUIRE_COLOR)
                    this.pickedPiece = {
                        "piece": this.board[y][x],
                        "location": {x: x, y: y},
                    }
                } 
            } else {
                await this.changeBackground(
                    this.pickedPiece.location.x, 
                    this.pickedPiece.location.y,
                    getSquireColor(
                        this.pickedPiece.location.x,
                        this.pickedPiece.location.y,
                    ),
                )
                const rect = this.canvas.getBoundingClientRect();
                let x = Math.floor((event.clientX - rect.x) / SQUIRE_SIZE)
                let y = Math.floor((event.clientY - rect.y) / SQUIRE_SIZE)
                if(x < 0 || x > 7 || y < 0 || y > 7) {
                    this.pickedPiece = null;
                    return ;
                }
                const played = this.game.play(
                    { row: this.pickedPiece.location.y, col: this.pickedPiece.location.x },
                    { row: y, col: x },
                )
                if(played) {
                    await this.render();
                } else {
                    console.log("invalid play")
                }

                this.pickedPiece = null;
            }
        });
    }
}