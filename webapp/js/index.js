var ws = new WebSocket("ws://localhost:8080/ws")

function wsMessage(ws) {
    return Bacon.fromEvent(ws, "message")
}

function wsError(ws) {
    return Bacon.fromEvent(ws, "error")
}

wsMessage(ws)
    .map(msg => "msg " + msg.data)
    .log()

wsError(ws)
    .map(msg => "error " + msg)
    .log()

Bacon
    .repeatedly(1000, [1])
    .scan(0, (x, y) => x + y)
    .take(11)
    .onValue(function(v) {
        ws.send("msg" + v)
    })

