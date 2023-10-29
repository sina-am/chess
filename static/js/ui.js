class PieceImages {
    constructor() {
        this.images = new Map();
    }
    
    async fetch(name, color) {
        return await (await fetch(`img/pieces/${name}-${color}.svg`)).text();
    }
    async get(name, color) {
        const cached = this.images.get(`${name}-${color}`);
        if(cached)  {
            return cached; 
        }

        const image = await this.fetch(name, color);
        this.images.set(`${name}-${color}`, image);
        return image; 
    }
} 

class ChessUI {
    constructor(boardElem, playerElem, opponentElem, board) {
        this.images = new PieceImages();
        this.boardElem = boardElem;
        this.playerElem = playerElem;
        this.opponentElem = opponentElem;
        this.gameStatusElem = document.getElementById("gameStatus");
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

    async squireClick(row, col, actor) {
        if(this.pickedPiece == null) {
            if(!actor.isMyPiece({row: row, col: col})) return;
            this.pickedPiece = {row: row, col: col};
            document.getElementById(`squire-${row}-${col}`).classList.add("picked");
        }else {
            actor.play(
                {row: this.pickedPiece.row, col: this.pickedPiece.col},
                {row: row, col: col},
            );
            document.getElementById(`squire-${this.pickedPiece.row}-${this.pickedPiece.col}`).classList.remove("picked");
            this.pickedPiece = null;
            await this.render();

            const winner = actor.hasWinner();
            if(winner) {
                this.gameOver(winner);
            }
            
        }
    }

    gameOver(winner) {
        this.gameStatusElem.innerText = `winner is ${winner}`;
    }

    async setUp(player, opponent, actor) {
        await this.setUpBoard(actor);
        await this.setUpBar(player, opponent);
        await this.render();
    }
    async setUpBar(player, opponent) {
        this.playerElem.innerText = player.name + " " + player.color;
        this.opponentElem.innerText = opponent.name + " " + opponent.color;
    }

    async setUpBoard(actor) {
        if (this.boardElem.firstChild) {
            for (let i = 0; i < 8; i++) {
                for (let j = 0; j < 8; j++) { 
                    const squire = document.getElementById(`squire-${i}-${j}`);
                    squire.onclick = async (event) => {
                        this.squireClick(i, j, actor); 
                    };
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
                    this.squireClick(i, j, actor); 
                };
                
                squire.classList.add("squire");
                row.appendChild(squire);
            }
            this.boardElem.appendChild(row);
        }
    }
}
