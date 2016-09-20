import ProtoBuf from 'protobufjs';

const KEY_UP = 38
const KEY_LEFT = 37
const KEY_RIGHT = 39

// Create the canvas
var canvas = document.createElement("canvas");
var ctx = canvas.getContext("2d");
canvas.width = 800;
canvas.height = 600;
document.body.appendChild(canvas);

var builder = ProtoBuf.loadJsonFile("js/superstellar_proto.json");
var Space = builder.build("superstellar.Space");
var UserInput = builder.build("superstellar.UserInput")

// Background image
var shipReady = false;
var shipThrustReady = false;

var shipImage = new Image();
shipImage.onload = function () {
  shipReady = true;
};
shipImage.src = "images/ship.png";

var shipThrustImage = new Image();
shipThrustImage.onload = function () {
	shipThrustReady = true;
};
shipThrustImage.src = "images/ship_thrust.png";

var ships = [];
var myID = 0;

var viewport = {vx: 0, vy: 0, width: 800, height: 600}

// TODO: Use config for this
var ws = new WebSocket("ws://127.0.0.1:8080/superstellar");

var frameCounter = 0;
var lastTime = Date.now();
var fps = 0;

ws.onmessage = function(e) {
	var fileReader = new FileReader();
	fileReader.onload = function() {
    ships = Space.decode(this.result).spaceships

    if (myID == 0) {
      for (var i in ships) {
        if (ships[i].id > myID) {
          myID = ships[i].id
        }
      }
    }
  };

	fileReader.readAsArrayBuffer(e.data);
};

// Handle keyboard controls
var keysDown = {};

addEventListener("keydown", function (e) {
	keysDown[e.keyCode] = true;
}, false);

addEventListener("keyup", function (e) {
	delete keysDown[e.keyCode];
}, false);

var translateToViewport = function (x, y, viewport) {
  var newX = x - viewport.vx + viewport.width / 2;
	var newY = -y + viewport.vy + viewport.height / 2;
	return {x: newX, y: newY}
}

var sendInput = function() {
	var thrust = KEY_UP in keysDown;

	var direction = "NONE";
	if (KEY_LEFT in keysDown) {
		direction = "LEFT";
	} else if (KEY_RIGHT in keysDown) {
		direction = "RIGHT";
	}

	var userInput = new UserInput(thrust, direction);
	var buffer = userInput.encode();

	if (ws.readyState == WebSocket.OPEN) {
			ws.send(buffer.toArrayBuffer());
	}

}

// Draw everything
var render = function () {
	ctx.beginPath();

	var myShip;

	if (ships.length > 0) {
    myShip = ships.find(function(ship) { return ship.id == myID })
    var ownPosition = {x: myShip.position.x/100, y: myShip.position.y/100};
    viewport = {vx: ownPosition.x, vy: ownPosition.y, width: 800, height: 600};
  }

	ctx.rect(0, 0, 800, 600);
	ctx.fillStyle = "black";
	ctx.fill();

	if (shipReady) {
		for (var shipID in ships) {
			var ship = ships[shipID]

			let image = ship.input_thrust ? shipThrustImage : shipImage

			var translatedPosition = translateToViewport(ship.position.x/100, ship.position.y/100, viewport)

			ctx.translate(translatedPosition.x, translatedPosition.y);
			ctx.fillStyle = "rgb(250, 250, 250)";
			ctx.font = "18px Helvetica";
			ctx.fillText(ship.id, -35, -60);

			ctx.rotate(ship.facing);
			ctx.drawImage(image, -30, -22);
			ctx.rotate(-ship.facing);
			ctx.translate(-translatedPosition.x, -translatedPosition.y);
		}
	}

	frameCounter++;

	if (frameCounter === 100) {
		frameCounter = 0;
		var now = Date.now();
		var delta = (now - lastTime) / 1000;
		fps = (100 / delta).toFixed(1);
		lastTime = now;
	}

	// Score
	ctx.fillStyle = "rgb(250, 250, 250)";
	ctx.font = "24px Helvetica";
	ctx.textAlign = "left";
	ctx.textBaseline = "top";
	ctx.fillText("Ships: " + ships.length, 580, 10);
	ctx.fillText("FPS: " + fps, 580, 40);
	if (undefined != myShip) {
    ctx.fillText("X: " + Math.floor(myShip.position.x / 100), 580, 70);
    ctx.fillText("Y: " + Math.floor(myShip.position.y / 100), 580, 100);
	}

	sendInput()
};

// The main game loop
var main = function () {
	render();
	// Request to do this again ASAP
	requestAnimationFrame(main);
};

// Let's play this game!
var then = Date.now();
main();
