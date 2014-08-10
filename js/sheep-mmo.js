// sheep-mmo - feed a sheep, rule the world.
// TODO: activeSheep hash for accessing sheep by id.
var offset = 16;
var board_height = 32;
var board_width = 48;

var scheme = "ws://";
var host = "localhost";
var port = "4977";
var ip = scheme + host + ":" + port;

var statusbox = $("#status-box");
var statusstatic = $("#status-static")
    .append($("<input type=button value='Gen Sheep'>")
            .click(function() { generateSheep() }));

var canvas = $('#screen')
canvas[0].width = board_width * offset;
canvas[0].height = board_height * offset;
var ctx = canvas[0].getContext("2d");

// All sheep currently active
var activeSheep = {};
// All other players' mouse positions
var activeMice = [];
// Global message to show on the canvas
var serverMessages = [];

var sheepDiffAttributes = ["name"];

var statusDisplayed = false;

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

    // conn is a global :D javascript, yay
    conn = new Conn(ip, $(canvas)[0]);
    $(canvas).on("mousemove", function(evt) {
        sendMouseMove(evt, ctx, conn);
    });
    $(canvas).click(function(evt) {
        processMouseClick(evt, ctx, conn);
    });
    var tick = 10;
    window.setInterval(loop, tick);
}

function processMouseClick(evt, ctx, conn) {
    //console.log("processMouseClick fired");
    var pos = getMousePos(ctx.canvas, evt);
    // just in case activeSheep gets clobbered in processMessages
    // TODO: are functions atomic? Ask Mary.
    var found = false;
    foundSheep = undefined;
    for (var i = 0; i < activeSheep.keys.length; i++) {
        var id = activeSheep.keys[i];
        if (activeSheep[id].type == "sheep") {
            var sheep = activeSheep[id];
            if (sheep === undefined) {
                continue;
            }
            if (pos.x >= sheep.x && pos.x < sheep.x + sheep.width &&
                pos.y >= sheep.y && pos.y < sheep.y + sheep.height) {
                found = true;
                // Pick the forward most sheep. If two sheep have the same
                // y coordinate, pick the left most one.
                if (typeof(foundSheep) == "undefined") {
                    foundSheep = sheep.id;
                } else if (activeSheep[foundSheep].Y < sheep.y) {
                    foundSheep = sheep.id;
                } else if (activeSheep[foundSheep].y == sheep.y &&
                           activeSheep[foundSheep].x > sheep.x) {
                    foundSheep = sheep.id;
                }
            }
        }
    }
    if (found) {
        displaySheepStatus();
    } else {
        clearMessage();
    }
}

// Returns true if they're different, false if they're the same.
function sheepDiff(sheep0, sheep1) {
    if (typeof(sheep0) === "undefined" ||
        typeof(sheep1) === "undefined") {
        return false;
    }
    for (var i = 0; i < sheepDiffAttributes.length; i++) {
        var attr = sheepDiffAttributes[i];
        if (sheep0[attr] !== sheep1[attr]) {
            return true;
        }
    }
    return false;
}

function clearMessage() {
    statusDisplayed = false;
    if (typeof(foundSheep) == "undefined") {
        statusbox.html("");
    } else {
        displaySheepStatus();
    }
}

function updateDisplay() {
    if (statusDisplayed === true) {
        displaySheepStatus();
    }
}

function displaySheepStatus() {
    statusDisplayed = true;
    var button = $("<input type='button' value='rename'>")
        .click(function () {displayRename(activeSheep[foundSheep].name) });
    statusbox.text("Name: " + activeSheep[foundSheep].name).append(button);
}

function displayRename(str) {
    statusDisplayed = false;
    statusbox.text("");
    var renamebutton = $("<input type='button' value='rename'>")
        .click(function() {
            sendRename($("#rename").val())
        });
    var cancelbutton = $("<input type='button' value='cancel'>")
        .click(function() {clearMessage()});
    statusbox.append("<input id='rename' type='text' value='" + str + "'>")
        .append(renamebutton).append(cancelbutton);
}

function sendRename(str) {
    conn.sendCheck("(rename " + activeSheep[foundSheep].id + " \"" + str + "\")");
    clearMessage();
}

function generateSheep() {
    conn.sendCheck("(gen-sheep)");
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
                if (stringbuffer != "") {
                    databuffer.push(stringbuffer);
                }
                result.push(databuffer);
                state = "out";
                continue;
            }
            if (data[i] == " ") {
                if (stringbuffer != "") {
                    databuffer.push(stringbuffer);
                }
                stringbuffer = "";
                continue;
            }
            if (data[i] == "\"") {
                if (stringbuffer.length != 0) {
                    state = "error";
                    continue;
                }
                state = "in-dquote";
                continue;
            }
            stringbuffer += data[i]
            break;
        case "in-dquote":
            if (data[i] == "\"") {
                databuffer.push(stringbuffer);
                stringbuffer = "";
                state = "in-message";
                continue;
            }
            stringbuffer += data[i];
            break;
        case "error":
            return undefined;
        }
    }
    if (state != "out") {
        return undefined;
    }
    return result
}

function newSheep(sheepmsg) {
    if (sheepmsg[0] !== "sheep" || sheepmsg.length !== 5) {
        console.log("bad sheep message " + sheepmsg);
        return undefined;
    }
    return {
        type: "sheep",
        id: sheepmsg[1],
        x: parseInt(sheepmsg[2]),
        y: parseInt(sheepmsg[3]),
        name: sheepmsg[4],
        width: images["sheep"].width,
        height: images["sheep"].height,
    };
}

function processMessages(msgs) {
    if (msgs.length === 0) {
        return;
    }
    ctx.clearRect(0, 0, $(canvas)[0].width, $(canvas)[0].height);
    var updated = false;
    activeSheep = {};
    activeSheep.keys = [];
    for (var i = 0; i < msgs.length; i++) {
        msg = msgs[i];
        switch (msg[0]) {
        case "sheep":
            var sheep = newSheep(msg);
            if (sheepDiff(sheep, activeSheep[sheep.id])) {
                updated = true;
            }
            activeSheep[sheep.id] = sheep;
            activeSheep.keys.push(sheep.id);
            ctx.drawImage(images["sheep"], sheep.x, sheep.y);
            break;
        case "mouse":
            ctx.drawImage(images["mouse"], parseInt(msg[2]), parseInt(msg[3]));
            break;
        }
    }
    if (updated === true) {
        updateDisplay();
    }
}

function loop() {
    processMessages(serverMessages);
    // Preserve last message
    serverMessages = getLast(serverMessages)
}

// Does not modify msgs.
function getLast(msgs) {
    var clone = msgs.slice(0);
    var result = [];
    var check = {};
    for (var i = clone.length - 1; i >= 0; i--) {
        // check to see if the message is a tick
        if (clone[i][0] == "tick") {
            return result;
        }
        var id = clone[i][0].concat(clone[i][1]);
        if (check[id] !== true) {
            check[id] = true;
            result.push(clone[i]);
        }
    }
    return result
}

function Conn(ip) {
    c = new WebSocket(ip);
    c.onclose = function() {
        console.log("connection closed");
    }
    c.onmessage = function(evt) {
        var msgs = decode(evt.data)
        if (msgs === undefined) {
            console.log("bad decode! " + evt.data);
            return;
        }
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
        x: Math.floor(evt.clientX - rect.left),
        y: Math.floor(evt.clientY - rect.top)
    };
}

function sendMouseMove(evt, ctx, conn) {
    var pos = getMousePos(ctx.canvas, evt);
    conn.sendCheck("(mouse " + pos.x + " " + pos.y + ")");
}

