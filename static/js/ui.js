class ChessUI {
    constructor(boardElem, playerElem, opponentElem, board) {
        this.boardElem = boardElem;
        this.playerElem = playerElem;
        this.opponentElem = opponentElem;
        this.pickedPiece = null;
        this.board = board;
    }

    render() {
        for(let i = 0; i < 8; i++) {
            for(let j = 0; j < 8; j++) {
                if(this.board[i][j] !== null) {
                    document.getElementById(`squire-${i}-${j}`).innerHTML = this.board[i][j].image;
                } else {
                    document.getElementById(`squire-${i}-${j}`).innerHTML = "";
                }
            }
        }
    }
    async getImages() {
        for(let i = 0; i < 8; i++) {
            for(let j = 0; j < 8; j++) {
                if(this.board[i][j] !== null) {
                    this.board[i][j].image = 
                        await (await fetch(`img/pieces/${this.board[i][j].name}-${this.board[i][j].color}.svg`)).text();
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
            this.render();
        }
    }

    async setUp(player, opponent, actor) {
        await this.setUpBoard(actor);
        await this.setUpBar(player, opponent);
    }
    async setUpBar(player, opponent) {
        this.playerElem.innerText = player.name + " " + player.color;
        this.opponentElem.innerText = opponent.name + " " + opponent.color;
    }

    async setUpBoard(actor) {
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
        await this.getImages();
        this.render();
    }
}