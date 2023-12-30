class OfflineChess {
    constructor(ui, playerName, standardBoard) {
        this.player = { name: playerName, color: null };
        this.status = "started";
        this.opponent = { name: "", color: null };
        this.board = standardBoard;
        this.ui = ui;
    }

    async start() {
        if (Math.random() > 0.5) {
            this.player.color = "white";
            this.opponent.color = "black";
        } else {
            this.player.color = "black";
            this.opponent.color = "white";
        }
        this.engine = new ChessEngine(this.board, this.player.color);
        this.ui.setUp(this.player, this.opponent, this);
        this.status = "started";
        setGameState(this.status);
    }
    isMyPiece(location) {
        if (this.board[location.row][location.col]) {
            if (this.board[location.row][location.col].color === this.engine.turn) {
                return true;
            }
        }
        return false;
    }
    hasWinner() {
        return this.engine.winner;
    }
    play(from, to) {
        if (!this.engine.movePiece(from, to)) {
            console.log("not your turn");
            return false
        }

        if (this.engine.winner) {
            this.gameOver();
        }
        return true
    }

    gameOver() {
        this.status = "game over";
        setGameState(this.status);
    }
    exit() {
        this.status = "exited";
        setGameState(this.status);
    }
}
class OnlineChess {
    constructor(ui, playerName, standardBoard, serverConnection) {
        this.player = { name: playerName, color: null };
        this.status = "waiting";
        this.opponent = {};
        this.board = standardBoard;
        this.server = serverConnection;
        this.ui = ui;
        this.winner = null;

        this.server.addMessageHandler(async (event) => {
            await this.update(JSON.parse(event.data))
        })
    }

    onstart(msg) {
        this.player.color = msg.color;
        this.opponent = {
            name: msg.name,
            color: oppositeColor(this.player.color),
        };

        this.engine = new ChessEngine(this.board, this.player.color);
        this.ui.setUp(this.player, this.opponent, this);
        this.status = "started";
        setGameState(this.status);
    }

    isMyPiece(location) {
        if (this.board[location.row][location.col]) {
            if (this.board[location.row][location.col].color === this.player.color) {
                return true;
            }
        }
        return false;
    }
    play(from, to) {
        if (this.status !== "started") return false;

        if (!this.engine.movePiece(from, to)) {
            return false;
        }
        this.server.send({
            "type": "play",
            "payload": {
                "move": { from: from, to: to },
            }
        });
        return true
    }

    async update(msg) {
        switch (msg.type) {
            case "started":
                this.onstart({ name: msg.payload.opponent, color: msg.payload.tile });
                setGameState("started");
                break;
            case "played":
                this.engine.movePiece(
                    msg.payload.move.from,
                    msg.payload.move.to,
                )
                await this.ui.render();
                break;
            case "ended":
                this.winner = msg.winner;
                if (msg.payload.winner === this.player.color) {
                    document.getElementById('gameWinner').innerText = "You won!";
                } else {
                    document.getElementById('gameWinner').innerText = "You lost";
                }
                this.status = "finished";
                this.server.send({
                    "type": "exit",
                    "payload": ""
                });
                break;
            default:
                break;
        }
        document.getElementById("gameStatus").innerText = this.status;
    }

    hasWinner() {
        return this.winner;
    }
    start() {
        this.server.ws.addEventListener("open", (event) => {
            this.server.send({
                "type": "start",
                "payload": {
                    "name": this.player.name,
                }
            });
        })
    }
    exit() {
        this.status = "exited";
        setGameState(this.status)
        this.server.send({
            "type": "exit",
            "payload": ""
        });
    }
}


