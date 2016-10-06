import * as PIXI from "pixi.js";
import Assets from './assets';
import * as Constants from './constants';
import * as Utils from './utils';
import Spaceship from './spaceship';
import { renderer, stage, globalState } from './globals';
import { initializeConnection, sendMessage, registerMessageHandler, UserInput } from './communicationLayer';

// TODO: Use config for this
const HOST = window.location.hostname;
const PORT = '8080';
const PATH = '/superstellar';

const KEY_UP = 38;
const KEY_LEFT = 37;
const KEY_RIGHT = 39;

document.body.appendChild(renderer.view);

const loadProgressHandler = (loader, resource) => {
  console.log(`progress: ${loader.progress}%`);
};

const spaceMessageHandler = (space) => {
  const ships = space.spaceships;
  const shipTexture = Assets.getTexture(Constants.SHIP_TEXTURE);

  let shipThrustFrames = [];

  Constants.FLAME_SPRITESHEET_FRAME_NAMES.forEach((frameName) =>  {
    shipThrustFrames.push(Assets.getTextureFromFrame(frameName));
  });

  for (var i in ships) {
    let shipId = ships[i].id;

    if (!globalState.spaceshipMap.has(shipId)) {
      const newSpaceship = new Spaceship(shipTexture, shipThrustFrames, ships[i]);

      globalState.spaceshipMap.set(shipId, newSpaceship);
    } else {
      globalState.spaceshipMap.get(shipId).updateData(ships[i]);
    }
  }
};

const helloMessageHandler = (message) => {
  globalState.clientId = message.myId;
};

const playerLeftHandler = (message) => {
  const playerId = message.id;

  let spaceship = globalState.spaceshipMap.get(playerId);

  spaceship.remove();

  globalState.spaceshipMap.delete(playerId);
};

registerMessageHandler(Constants.HELLO_MESSAGE,       helloMessageHandler);
registerMessageHandler(Constants.SPACE_MESSAGE,       spaceMessageHandler);
registerMessageHandler(Constants.PLAYER_LEFT_MESSAGE, playerLeftHandler);

PIXI.loader.
  add([Constants.SHIP_TEXTURE, Constants.BACKGROUND_TEXTURE, Constants.FLAME_SPRITESHEET]).
  on("progress", loadProgressHandler).
  load(setup);

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
let thrustAnim;

function setup() {
  initializeConnection(HOST, PORT, PATH);

  const bgTexture = Assets.getTexture(Constants.BACKGROUND_TEXTURE);

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
  
  let userInput = new UserInput(thrust, direction);

  sendMessage(userInput);
}

// Draw everything
var render = function () {
  if (!globalState.clientId) { return }
  let myShip;

  let backgroundPos = Utils.translateToViewport(0, 0, viewport);
  tilingSprite.tilePosition.set(backgroundPos.x, backgroundPos.y);

  if (globalState.spaceshipMap.size > 0) {
    myShip = globalState.spaceshipMap.get(globalState.clientId);
    viewport = myShip.viewport();
  }

  globalState.spaceshipMap.forEach((spaceship) => spaceship.update(viewport));
  frameCounter++;

  if (frameCounter === 100) {
    frameCounter = 0;
    var now = Date.now();
    var delta = (now - lastTime) / 1000;
    fps = (100 / delta).toFixed(1);
    lastTime = now;
  }

  let shipCount = globalState.spaceshipMap.size;

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

