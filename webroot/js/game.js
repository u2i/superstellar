import ProtoBuf from 'protobufjs';
import * as PIXI from "pixi.js";
import * as Constants from './constants.js';
import * as Utils from './utils.js';
import Spaceship from './spaceship.js';
import { renderer, stage } from './globals.js';

// TODO: Use config for this
const ws = new WebSocket("ws://" + window.location.hostname + ":8080/superstellar");

const webSocketMessageReceived = (e) => {
  var fileReader = new FileReader();
  fileReader.onload = function() {
    const ships = Space.decode(this.result).spaceships;

    for (var i in ships) {
      let shipId = ships[i].id;

      if (!(shipId in shipIds)) {
	const newSpaceship = new Spaceship(shipTexture, shipThrustFrames, ships[i]);

	spaceships.push(newSpaceship);

	shipIds[shipId] = newSpaceship;
      } else {
	shipIds[shipId].updateData(ships[i]);
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

const shipIds = {};
const spaceships = [];

document.body.appendChild(renderer.view);

const builder = ProtoBuf.loadJsonFile(Constants.PROTOBUF_DEFINITION);
const Space = builder.build(Constants.SPACE_DEFINITION);
const UserInput = builder.build(Constants.USER_INPUT_DEFINITION);


const loadProgressHandler = (loader, resource) => {
  console.log(`progress: ${loader.progress}%`);
};

PIXI.loader.
  add([Constants.SHIP_TEXTURE, Constants.BACKGROUND_TEXTURE, Constants.FLAME_SPRITESHEET]).
  on("progress", loadProgressHandler).
  load(setup);

let shipTexture;
let shipThrustTexture;
let bgTexture;

let tilingSprite;

const hudTextStyle = {
  fontFamily: 'Helvetica',
  fontSize: '24px',
  fill: '#FFFFFF',
  align: 'left',
  textBaseline: 'top'
};

const buildHudText = (shipCount, fps, x, y) => {
  let text = "Ships: " + shipCount + "\n";
  text += "FPS: " + fps + "\n";
  text += "X: " + x + "\n";
  text += "Y: " + y + "\n";

  return text;
}

let hudText;
let shipThrustFrames = [];
let thrustAnim;

function setup() {
  shipTexture = PIXI.loader.resources[Constants.SHIP_TEXTURE].texture;

  Constants.FLAME_SPRITESHEET_FRAME_NAMES.forEach((frameName) =>  {
    shipThrustFrames.push(PIXI.Texture.fromFrame(frameName));
  });

  ws.onmessage = webSocketMessageReceived;

  bgTexture = PIXI.loader.resources[Constants.BACKGROUND_TEXTURE].texture;

  tilingSprite = new PIXI.extras.TilingSprite(bgTexture, renderer.width, renderer.height);
  stage.addChild(tilingSprite);

  hudText = new PIXI.Text('', hudTextStyle);
  hudText.x = 580;
  hudText.y = 0;
  stage.addChild(hudText);

  // Let's play this game!
  var then = Date.now();
  main();
}

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
  var myShip;

  let backgroundPos = Utils.translateToViewport(0, 0, viewport);
  tilingSprite.tilePosition.set(backgroundPos.x, backgroundPos.y);

  if (spaceships.length > 0) {
    myShip = spaceships.find(function(ship) { return ship.id == myID })
    viewport = myShip.viewport();
  }

  spaceships.forEach((spaceship) => spaceship.update(viewport));
  frameCounter++;

  if (frameCounter === 100) {
    frameCounter = 0;
    var now = Date.now();
    var delta = (now - lastTime) / 1000;
    fps = (100 / delta).toFixed(1);
    lastTime = now;
  }

  let shipCount = spaceships.length;

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

