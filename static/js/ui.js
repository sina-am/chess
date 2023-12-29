class PieceImages {
    constructor() {
        this.images = new Map();
    }

    async fetch(name, color) {
        return await (await fetch(`img/pieces/${name}-${color}.svg`)).text();
    }
    async get(name, color) {
        const cached = this.images.get(`${name}-${color}`);
        if (cached) {
            return cached;
        }

        const image = await this.fetch(name, color);
        this.images.set(`${name}-${color}`, image);
        return image;
    }
}

class ChessUI {
    constructor(boardElem, playerElem, opponentElem, gameStatusElem, board) {
        this.images = new PieceImages();
        this.boardElem = boardElem;
        this.playerElem = playerElem;
        this.opponentElem = opponentElem;
        this.gameStatusElem = gameStatusElem
        this.pickedPiece = null;
        this.board = board;
    }
    async render() {
        for(let i = 0; i < 8; i++) {
            for(let j = 0; j < 8; j++) {
                if(this.board[i][j] !== null) {
                    document.getElementById(`squire-${i}-${j}`).innerHTML = 
                        await this.images.get(this.board[i][j].name, this.board[i][j].color);
                } else {
                    document.getElementById(`squire-${i}-${j}`).innerHTML = "";
                }
            }
        }
    }

    async squireClick(row, col) {
        if (this.pickedPiece == null) {
            if (!this.game.isMyPiece({ row: row, col: col })) return;
            this.pickedPiece = { row: row, col: col };
            document.getElementById(`squire-${row}-${col}`).classList.add("picked");
        } else {
            const played = this.game.play(
                { row: this.pickedPiece.row, col: this.pickedPiece.col },
                { row: row, col: col },
            );
            console.log(played);
            if (played) {
                await this.render();
            }
            document.getElementById(`squire-${this.pickedPiece.row}-${this.pickedPiece.col}`).classList.remove("picked");
            this.pickedPiece = null;
            const winner = this.game.hasWinner();
            if (winner) {
                this.gameOver(winner);
            }

        }
    }

    gameOver(winner) {
        this.gameStatusElem.innerText = `winner is ${winner}`;
    }

    async setUp(player, opponent, game) {
        this.game = game;
        await this.setUpBoard();
        await this.setUpBar(player, opponent);
    }
    async setUpBar(player, opponent) {
        this.playerElem.innerText = player.name + " " + player.color;
        this.opponentElem.innerText = opponent.name + " " + opponent.color;
    }

    async setUpBoard() {
        if (this.boardElem.firstChild) {
            for (let i = 0; i < 8; i++) {
                for (let j = 0; j < 8; j++) {
                    const squire = document.getElementById(`squire-${i}-${j}`);
                    squire.onclick = async (event) => {
                        this.squireClick(i, j);
                    };
                    if (this.board[i][j] !== null) {
                        squire.innerHTML =
                            await this.images.get(this.board[i][j].name, this.board[i][j].color);
                    } else {
                        squire.innerHTML = "";
                    }
                }
            }
            return;
        }
        for (let i = 0; i < 8; i++) {
            const row = document.createElement("tr");
            for (let j = 0; j < 8; j++) {
                const squire = document.createElement("td");
                squire.id = `squire-${i}-${j}`;
                squire.onclick = async (event) => {
                    this.squireClick(i, j);
                };
                if (this.board[i][j] !== null) {
                    squire.innerHTML =
                        await this.images.get(this.board[i][j].name, this.board[i][j].color);
                } else {
                    squire.innerHTML = "";
                }

                squire.classList.add("squire");
                row.appendChild(squire);
            }
            this.boardElem.appendChild(row);
        }
    }
}
