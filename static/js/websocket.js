class ServerConnection {
    constructor(ws) {
        this.ws = ws;
    }

    addMessageHandler(handler) {
        this.ws.addEventListener("message", handler);
    }

    send(msg) {
        this.ws.send(JSON.stringify(msg))
    }
}