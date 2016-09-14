// Create the canvas
var canvas = document.createElement("canvas");
var ctx = canvas.getContext("2d");
canvas.width = 800;
canvas.height = 600;
document.body.appendChild(canvas);


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

ws.onmessage = function(e) {
	ships = $.evalJSON(e.data).spaceships;
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

    var shipsArray = Object.keys(ships).map(function(val) { return ships[val] });

	if (shipReady) {
		for (var shipID in ships) {
			ship = ships[shipID]

			image = ship.thrust ? shipThrustImage : shipImage

			var translatedPosition = translateToViewport(ship.position.x, ship.position.y, viewport)

			ctx.translate(translatedPosition.x, translatedPosition.y);
			ctx.fillStyle = "rgb(250, 250, 250)";
			ctx.font = "18px Helvetica";
			ctx.fillText(shipID.split('-')[0], -35, -60);
			var angle = Math.atan2(-ship.facing.y, ship.facing.x);

			ctx.rotate(angle);
			ctx.drawImage(image, -30, -22);
			ctx.rotate(-angle);
			ctx.translate(-translatedPosition.x, -translatedPosition.y);
		}
	}

	// Score
	ctx.fillStyle = "rgb(250, 250, 250)";
	ctx.font = "24px Helvetica";
	ctx.textAlign = "left";
	ctx.textBaseline = "top";
	ctx.fillText("Ships: " + shipsArray.length, 580, 10);
	if (shipReady) {
		ctx.fillText("pos: " + Math.round(shipsArray[0].position.x) + ", " + Math.round(shipsArray[0].position.y), 580, 50);
		ctx.fillText("fac: " + (shipsArray[0].facing.x).toFixed(2) + ", " + (shipsArray[0].facing.y).toFixed(2), 580, 90);
	}
};

// The main game loop
var main = function () {
	var now = Date.now();
	var delta = now - then;

	render();

	then = now;

	// Request to do this again ASAP
	requestAnimationFrame(main);
};

// Let's play this game!
var then = Date.now();
main();
