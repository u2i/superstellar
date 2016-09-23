import ProtoBuf from 'protobufjs';
import * as PIXI from "pixi.js";

const renderer = new PIXI.WebGLRenderer(800, 600);
const stage = new PIXI.Container();

// TODO: Use config for this
const ws = new WebSocket("ws://127.0.0.1:8080/superstellar");

const webSocketMessageReceived = (e) => {
  var fileReader = new FileReader();
  fileReader.onload = function() {
    ships = Space.decode(this.result).spaceships
    
    for (var i in ships) {
      let shipId = ships[i].id;
      if (!(shipId in sprites)) {
	sprites[shipId] = new PIXI.Sprite(shipTexture);
	stage.addChild(sprites[shipId]);
      }
    }

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


const KEY_UP = 38;
const KEY_LEFT = 37;
const KEY_RIGHT = 39;

document.body.appendChild(renderer.view);

var builder = ProtoBuf.loadJsonFile("js/superstellar_proto.json");
var Space = builder.build("superstellar.Space");
var UserInput = builder.build("superstellar.UserInput")

const loadProgressHandler = (loader, resource) => {
  console.log(`progress: ${loader.progress}%`);
};

PIXI.loader.
  add(["images/ship.png", "images/ship_thrust.png"]).
  on("progress", loadProgressHandler).
  load(setup);

let shipTexture;
let shipThrustTexture;
function setup() {
  shipTexture = PIXI.loader.resources["images/ship.png"].texture;
  shipThrustTexture = PIXI.loader.resources["images/ship_thrust.png"].texture;

  ws.onmessage = webSocketMessageReceived;

  // Let's play this game!
  var then = Date.now();
  main();
}

let sprites = {};
var ships = [];
var myID = 0;

var viewport = {vx: 0, vy: 0, width: 800, height: 600}

var frameCounter = 0;
var lastTime = Date.now();
var fps = 0;

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

const hudTextStyle = {
  fontFamily: 'Helvetica',
  fontSize: '24px',
  fill: '#FFFFFF',
  align: 'left',
  textBaseline: 'top'
};
let hudText = new PIXI.Text('', hudTextStyle);
hudText.x = 580;
hudText.y = 0;
stage.addChild(hudText);

const buildHudText = (shipCount, fps, x, y) => {
  let text = "Ships: " + shipCount + "\n";
  text += "FPS: " + fps + "\n";
  text += "X: " + x + "\n";
  text += "Y: " + y + "\n";

  return text;
}

// Draw everything
var render = function () {
  var myShip;

  if (ships.length > 0) {
    myShip = ships.find(function(ship) { return ship.id == myID })
    var ownPosition = {x: myShip.position.x/100, y: myShip.position.y/100};
    viewport = {vx: ownPosition.x, vy: ownPosition.y, width: 800, height: 600};
  }

  for (var idx in ships) {
    let ship = ships[idx]
    let sprite = sprites[ship.id];

    var translatedPosition = translateToViewport(ship.position.x/100, ship.position.y/100, viewport)

    sprite.position.set(translatedPosition.x, translatedPosition.y);
    sprite.pivot.set(sprite.width / 2, sprite.height / 2);
    sprite.rotation = ship.facing;
  }

  frameCounter++;

  if (frameCounter === 100) {
    frameCounter = 0;
    var now = Date.now();
    var delta = (now - lastTime) / 1000;
    fps = (100 / delta).toFixed(1);
    lastTime = now;
  }

  let shipCount = ships.length;

  let x = myShip ? Math.floor(myShip.position.x / 100) : '?';
  let y = myShip ? Math.floor(myShip.position.y / 100) : '?';

  hudText.text = buildHudText(shipCount, fps, x, y);
  renderer.render(stage);
  sendInput()
};

// The main game loop
var main = function () {
  render();
  // Request to do this again ASAP
  requestAnimationFrame(main);
};

