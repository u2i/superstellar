// Create the canvas
var canvas = document.createElement("canvas");
var ctx = canvas.getContext("2d");
canvas.width = 800;
canvas.height = 600;
document.body.appendChild(canvas);

var ProtoBuf = dcodeIO.ProtoBuf;
var builder = ProtoBuf.loadJsonFile("js/superstellar_proto.json");
var Space = builder.build("superstellar.Space");

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

var ships = {};

var viewport = {vx: 0, vy: 0, width: 800, height: 600}

var ws = new WebSocket("ws://" + window.location.host + "/superstellar");

var frameCounter = 0;
var lastTime = Date.now();
var fps = 0;

ws.onmessage = function(e) {
	var fileReader = new FileReader();
	fileReader.onload = function() {
			ships = Space.decode(this.result).spaceships
	};

	fileReader.readAsArrayBuffer(e.data);
};

// Handle keyboard controls
var keysDown = {};

addEventListener("keydown", function (e) {
	var direction
	switch (e.keyCode) {
		case 38:
			direction = "up"
			break;
		case 40:
			direction = "down"
			break;
		case 37:
			direction = "left"
			break;
		case 39:
			direction = "right"
			break;
	}

	ws.send($.toJSON({client_id: 1, direction: direction}))
}, false);

var translateToViewport = function (x, y, viewport) {
  var newX = x + viewport.vx + viewport.width / 2;
	var newY = -y - viewport.vy + viewport.height / 2;
	return {x: newX, y: newY}
}

// Draw everything
var render = function () {
	ctx.beginPath();
	ctx.rect(0, 0, 800, 600);
	ctx.fillStyle = "black";
	ctx.fill();

	if (shipReady) {
		for (var shipID in ships) {
			var ship = ships[shipID]

			image = ship.input_thrust ? shipThrustImage : shipImage

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
