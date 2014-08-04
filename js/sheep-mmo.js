// sheep-mmo - feed a sheep, rule the world.
var offset = 16;
var board_height = 32;
var board_width = 48;

var scheme = "ws://";
var host = "localhost";
var port = "4977";
var ip = scheme + host + ":" + port;

var canvas = $('#screen').width(board_width * offset).height(board_height * offset);
var ctx = canvas[0].getContext("2d");

window.onload = function() {
    if (!window["WebSocket"]) {
        console.log("your browser does not support websockets.");
        return;
    }

    var conn = new Conn(ip);
    $(canvas).on("mousemove", function(evt) {
        sendMouseMove(evt, $(canvas)[0], conn);
    });
}

function Conn(ip) {
    c = new WebSocket(ip);
    c.onclose = function() {
        console.log("connection closed");
    }
    c.onmessage = function(evt) {
        var msgs = decode(evt.data)
        
        //console.log(evt.data)
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

function sendMouseMove(evt, canvas, conn) {
    var pos = getMousePos(canvas, evt);
    conn.send("(mouse " + pos.x + " " + pos.y + ")")
}
