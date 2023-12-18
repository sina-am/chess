class OfflineChess{
    constructor(playerName, standardBoard) {
        this.player = {name: playerName, color: null};
        this.status = "started";
        this.opponent = {name: "", color: null};
        this.board = standardBoard;

    }
    async start() {
        if(Math.random() > 0.5) {
            this.player.color = "white";
            this.opponent.color = "black";
        } else {
            this.player.color = "black";
            this.opponent.color = "white";
        }
        this.engine = new ChessEngine(this.board, this.player.color);
        this.ui = new ChessUI(
            document.getElementById("gameBoard"), 
            document.getElementById("playerName"), 
            document.getElementById("opponentName"), 
            this.board, 
        );

        this.ui.setUp(this.player, this.opponent, this);
        this.status = "started";
        setGameState(this.status);
    }
    isMyPiece(location) {
        if(this.board[location.row][location.col]) {
            if(this.board[location.row][location.col].color === this.engine.turn) {
                return true;
            }
        }
        return false;
    }
    hasWinner() {
        return this.engine.winner;
    }
    play(from, to) {
        if(!this.engine.movePiece(from, to)) {
            console.log("not your turn");
            return
        }

        if(this.engine.winner) {
            this.gameOver();
        }
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
    constructor(ws, playerName, standardBoard) {
        this.player = {name: playerName, color: null};
        this.status = "waiting";
        this.opponent = {};
        this.board = standardBoard; 
        this.ws = ws;
        this.ui = null;
        this.winner = null;
        ws.addEventListener("message", async (event) => {
            this.update(JSON.parse(event.data))
        });
    }

    onstart(msg) {
        this.player.color = msg.color;
        this.opponent = {
            name: msg.name,
            color: oppositeColor(this.player.color),
        };

        this.engine = new ChessEngine(this.board, this.player.color);

        this.ui = new ChessUI(
            document.getElementById("gameBoard"), 
            document.getElementById("playerName"), 
            document.getElementById("opponentName"), 
            this.board, 
        );

        this.ui.setUp(this.player, this.opponent, this);
        this.status = "started";
        setGameState(this.status);
    }

    isMyPiece(location) {
        if(this.board[location.row][location.col]) {
            if(this.board[location.row][location.col].color === this.player.color) {
                return true;
            }
        }
        return false;
    }
    play(from, to) {
        if (this.status !== "started") return;

        if(!this.engine.movePiece(from, to)) {
            console.log("not your turn");
            return;
        }
        this.ws.send(JSON.stringify({
            "type": "play",
            "payload": {
                "move": {from: from, to: to},
            }
        }));
    }

    update(msg) {
        switch (msg.type) {
            case "started":
                this.onstart({name: msg.payload.opponent, color: msg.payload.tile}); 
                setGameState("started");
                break;
            case "played":
                this.engine.movePiece(
                    msg.payload.move.from,
                    msg.payload.move.to,
                ) 
                this.ui.render();
                break;
            case "ended":
                this.winner = msg.winner;
                if (msg.payload.winner === this.player.color) {
                    document.getElementById('gameWinner').innerText = "You won!";
                } else {
                    document.getElementById('gameWinner').innerText = "You lost";
                }
                this.status = "finished";
                this.ws.send(JSON.stringify({
                    "type": "exit",
                    "payload": ""
                }));
                break;
            default:
                break;
        }
        document.getElementById("gameStatus").innerText = this.status;
    }
   
    hasWinner() {
        return this.winner;
    }
    async start() {
        this.ws.send(JSON.stringify({
            "type": "start",
            "payload": {
                "name": this.player.name,
            }
        }));
    }
    exit() {
        this.status = "exited";
        setGameState(this.status)
        this.ws.send(JSON.stringify({
            "type": "exit",
            "payload": ""
        }));
    }
}

