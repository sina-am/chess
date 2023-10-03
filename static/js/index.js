// state: "waiting" | "started" | "exited"
let isOnline = false;
const ws = new WebSocket("ws://localhost:8080/ws");


function displayElementById(id, display) {
    if (display) {
        document.getElementById(id).classList.remove('d-none');
    } else {
        document.getElementById(id).classList.add('d-none');
    }
}

function setGameState(state) {
    switch (state) {
        case "exited":
            displayElementById("game", false);
            displayElementById("waiting", false);
            displayElementById("startButton", true);
            displayElementById("startStage", true);
            break;
        case "waiting":
            displayElementById("game", false);
            displayElementById("startStage", true);
            displayElementById("waiting", true);
            displayElementById("startButton", false);
            break;
        case "started":
            displayElementById("startStage", false);
            displayElementById("game", true);
            break;
    }
}

function switchStage(stage) {
    switch (stage) {
        case "lobby":
            displayElementById("game", false);
            displayElementById("waiting", false);
            displayElementById("startButton", true);
            displayElementById("startStage", true);
            break;
        case "optionMenu":
            displayElementById("optionMenu", true);
            displayElementById("game", false);
            displayElementById("waiting", false);
            displayElementById("startButton", false);
            displayElementById("startStage", false);
            break;
        case "gameModeSelected":
            break;
        case "started":
            displayElementById("optionMenu", false);
            displayElementById("startStage", false);
            displayElementById("game", true);
            break;
        case "waiting":
            displayElementById("optionMenu", false);
            displayElementById("game", false);
            displayElementById("startStage", true);
            displayElementById("waiting", true);
            displayElementById("startButton", false);
            break;
    }
}

async function onGameStart(event) {
    const gameMode = document.getElementById('optionMenuSelect').value;

    let chess = null;
    if(gameMode === "online") {
        const playerName = document.getElementById("playerNameInput").value;
        if (playerName === "") {
            let myModal = new bootstrap.Modal(document.getElementById('playerNameModal'), {});
            myModal.toggle();
            return;
        }
        chess = new OnlineChess(ws, playerName);
        switchStage("waiting");
    } else {
        chess = new OfflineChess("");
        switchStage("started");
    }

    await chess.start();
    document.getElementById("exitBtn").onclick = (event) => {
        chess.exit();
    };
}

window.onload = async () => {
    ws.addEventListener("message", async (event) => {
        console.log("websocket new message: ", event.data);
    });

    ws.addEventListener("open", async (event) => {
        console.log("websocket connected");
        isOnline = true;
    });
    ws.addEventListener("close", async (event) => {
        console.log("websocket connection closed");
        isOnline = false;
    });
    ws.addEventListener("error", (event) => {
        console.log("websocket error: ", event);
        isOnline = false;
    });
}

