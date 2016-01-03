var ws = new WebSocket(wsURL("ws"))
var wsSend = wsSender(ws)

var btn1 = $('#button1').asEventStream("click").map("1")
var btn2 = $('#button2').asEventStream("click").map("2")
var btn3 = $('#button3').asEventStream("click").map("3")
var btn4 = $('#button4').asEventStream("click").map("4")
var btn5 = $('#button5').asEventStream("click").map("5")

var btns = Bacon.mergeAll(btn1, btn2, btn3, btn4, btn5)

btns
    .flatMap(btn => wsSend(btn))
    .onError(err => alert(err, "alert-danger"))

Bacon
    .fromEvent(ws, "open")
    .onValue(() => alert("WebSocket opened", "alert-success", 2000))

Bacon
    .fromEvent(ws, "message")
    .onValue(v => {
        $("#time").html('<h1>'+v.data+'</h1>')
    })

Bacon
    .fromEvent(ws, "error")
    .onValue(() => alert("WebSocket error", "alert-danger"))

Bacon
    .fromEvent(ws, "close")
    .onValue(closeEvt => alert("WebSocket closed (code=" + closeEvt.code + ")", "alert-warning"))

function wsSender(ws) {
    return msg => {
        return Bacon.fromBinder(sink => {
            if (ws.readyState !== WebSocket.OPEN) {
                sink(Bacon.Error("invalid WebSocket state: " + ws.readyState))
                sink(Bacon.End())
                return
            }

            try {
                ws.send(msg)
            } catch (err) {
                sink(Bacon.Error(err))
            }

            sink(Bacon.End())
        })
    }
}

function wsURL(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + l.pathname + s;
}

function alert(msg, cls, delay) {
    var t = document.querySelector("#alertTemplate")
    var alrtTpl = t.content.querySelector("div>div>div")
    alrtTpl.innerHTML = msg
    alrtTpl.className = "alert " + cls

    $(document.importNode(t.content, true)).appendTo($("#container"))

    delay = delay || 0

    if (delay > 0) {
        var alrt = $("#container").find(".row:last>div>div")
        var row = $("#container").find(".row:last")
        $(alrt)
            .fadeTo(delay, 500)
            .slideUp(500, () => {
                $(alrt).alert("close")
                $(row).remove()
            })
    }
}
