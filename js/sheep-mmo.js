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

// All sheep currently active
var activeSheep = [];
// All other players' mouse positions
var activeMice = [];
// Global message to show on the canvas
var message = "";
var serverMessages = [];

var images = {
    "sheep": function() {
        var i = new Image();
        i.src = "/res/sheep.png"
        i.height = 40;
        i.width = 38;
        return i;
    }(),
    "mouse": function() {
        var i = new Image();
        i.src = "/res/purplepointer.png"
        i.height = 24;
        i.width = 24;
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
    $(canvas).click(function(evt) {
        processMouseClick(evt, activeSheep, ctx, conn);
    });
    window.setInterval(loop, 30);
}

function processMouseClick(evt, activeSheep, ctx, conn) {
    //console.log("processMouseClick fired");
    var pos = getMousePos(ctx.canvas, evt);
    // just in case activeSheep gets clobbered in processMessages
    // TODO: are functions atomic? Ask Mary.
    var snapshot = activeSheep.slice(0);
    var found = false;
    for (var i = 0; i < snapshot.length; i++) {
        if (snapshot[i][0] == "sheep") {
            var sheepx = parseInt(snapshot[i][2]);
            var sheepy = parseInt(snapshot[i][3]);
            if (pos.x >= sheepx && pos.x < sheepx + images["sheep"].width &&
                pos.y >= sheepy && pos.y < sheepy + images["sheep"].height) {
                found = true;
                // Pick the forward most sheep. If two sheep have the same
                // y coordinate, pick the left most one.
                if (typeof(foundSheepY) == "undefined") {
                    var foundSheepY = sheepy;
                    var foundSheepX = sheepx;
                    var foundSheep = snapshot[i];
                } else if (foundSheepY < sheepy) {
                    foundSheepY = sheepy;
                    foundSheepX = sheepx;
                    foundSheep = snapshot[i];
                } else if (foundSheepY == sheepy && foundSheepX > sheepx) {
                    foundSheepY = sheepy;
                    foundSheepX = sheepx;
                    foundSheep = snapshot[i];
                }
            }
        }
    }
    if (found) {
        displayMessage("found one! id: " + foundSheep);
    } else {
        clearMessage();
    }
}

function clearMessage() {
    message = "";
}

function displayMessage(str) {
    message = str;
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
    activeSheep = [];
    for (var i = 0; i < msgs.length; i++) {
        msg = msgs[i];
        switch (msg[0]) {
        case "sheep":
            activeSheep.push(msg);
            ctx.drawImage(images["sheep"], parseInt(msg[2]), parseInt(msg[3]));
            break;
        // case "mouse":
        //     ctx.drawImage(images["mouse"], parseInt(msg[1]), parseInt(msg[2]));
        //     break;
        }
    }
    if (message != "") {
        ctx.fillStyle = "black";
        ctx.font = "bold 16px Arial";
        ctx.fillText(message, 10, 20);
    }
    
}

function loop() {
    processMessages(serverMessages);
    // Preserve last message
    if (serverMessages.length > 0) {
        serverMessages = [serverMessages[serverMessages.length - 1]];
    }
}

function Conn(ip) {
    c = new WebSocket(ip);
    c.onclose = function() {
        console.log("connection closed");
    }
    c.onmessage = function(evt) {
        var msgs = decode(evt.data)
        for (i = 0; i < msgs.length; i++) {
            serverMessages.push(msgs[i])
        }
    }
    // Only sends a message if the connection is open.
    c.sendCheck = function(msg) {
        if (this.readyState == this.OPEN) {
            this.send(msg)
        }
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
    conn.sendCheck("(mouse " + pos.x + " " + pos.y + ")")
}
