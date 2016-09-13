// Create the canvas
var canvas = document.createElement("canvas");
var ctx = canvas.getContext("2d");
canvas.width = 800;
canvas.height = 600;
document.body.appendChild(canvas);


// Background image
var shipReady = false;
var shipImage = new Image();
shipImage.onload = function () {
	shipReady = true;
};
shipImage.src = "images/ship.png";

var ships = {};

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

			ctx.translate(ship.position.x, ship.position.y);
			ctx.fillStyle = "rgb(250, 250, 250)";
			ctx.font = "18px Helvetica";
			ctx.fillText(shipID.split('-')[0], -35, -60);
			var angle = Math.atan(ship.facing.y, ship.facing.x);

			ctx.rotate(angle);
			ctx.drawImage(shipImage, -30, -15);
			ctx.rotate(-angle);
			ctx.translate(-ship.position.x, -ship.position.y);
		}
	}

	// Score
	ctx.fillStyle = "rgb(250, 250, 250)";
	ctx.font = "24px Helvetica";
	ctx.textAlign = "left";
	ctx.textBaseline = "top";
	ctx.fillText("Ships: " + shipsArray.length, 680, 10);
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
