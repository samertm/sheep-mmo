// sheep-mmo - feed a sheep, rule the world.
var offset = 16;
var board_height = 32;
var board_width = 48;

var scheme = "ws://";
var host = "localhost";
var port = "4977";
var ip = scheme + host + ":" + port;

var canvas = $('#screen')
canvas[0].width = board_width * offset;
canvas[0].height = board_height * offset;
var ctx = canvas[0].getContext("2d");

var images = {
    "sheep": function() {
        var i = new Image();
        i.src = "/res/sheep.png"
        i.height = 40;
        i.width = 38;
        return i;
    }()
}

window.onload = function() {
    if (!window["WebSocket"]) {
        console.log("your browser does not support websockets.");
        return;
    }

    var conn = new Conn(ip, $(canvas)[0]);
    $(canvas).on("mousemove", function(evt) {
        sendMouseMove(evt, ctx, conn);
    });
}

function decode(data) {
    var result = [];
    var state = "out";
    for (var i = 0; i < data.length; i++) {
        switch (state) {
        case "out":
            if (data[i] == "(") {
                var databuffer = [];
                var stringbuffer = "";
                state = "in-message";
                continue;
            }
            break;
        case "in-message":
            if (data[i] == ")") {
                databuffer.push(stringbuffer);
                result.push(databuffer);
                state = "out"
                continue;
            }
            if (data[i] == " ") {
                databuffer.push(stringbuffer);
                stringbuffer = "";
                continue;
            }
            stringbuffer += data[i]
            break;
        }
    }
    return result
}

function processMessages(msgs) {
    ctx.clearRect(0, 0, $(canvas)[0].width, $(canvas)[0].height);
    for (var i = 0; i < msgs.length; i++) {
        msg = msgs[i];
        if (msg[0] == "sheep") {
            ctx.drawImage(images["sheep"], parseInt(msg[1]), parseInt(msg[2]));
        }
    }
}

function Conn(ip) {
    c = new WebSocket(ip);
    c.onclose = function() {
        console.log("connection closed");
    }
    c.onmessage = function(evt) {
        var msgs = decode(evt.data)
        processMessages(msgs)
    }
    return c
}

function getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
        x: evt.clientX - rect.left,
        y: evt.clientY - rect.top
    };
}

function sendMouseMove(evt, ctx, conn) {
    var pos = getMousePos(ctx.canvas, evt);
    if (conn.readyState == conn.OPEN) {
        conn.send("(mouse " + pos.x + " " + pos.y + ")")
    }
}
