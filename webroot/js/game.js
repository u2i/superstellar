import * as PIXI from "pixi.js";
import Assets from './assets';
import * as Constants from './constants';
import * as Utils from './utils';
import Spaceship from './spaceship';
import { renderer, stage, globalState } from './globals';
import { initializeConnection, sendMessage, registerMessageHandler, UserMessage } from './communicationLayer';

// TODO: Use config for this
const HOST = window.location.hostname;
const PORT = '8080';
const PATH = '/superstellar';

let overlay;

function AnnulusFilter() {
  PIXI.Filter.call(this, null, shaderContent);
  this.uniforms.worldCoordinates = new Float32Array([0.0, 0.0]);
  this.uniforms.worldSize = new Float32Array([1000.0, 1400.0]);
  this.uniforms.magicMatrix = new PIXI.Matrix;
}

AnnulusFilter.prototype = Object.create(PIXI.Filter.prototype);
AnnulusFilter.prototype.constructor = AnnulusFilter;

Object.defineProperties(AnnulusFilter.prototype,
{
  worldCoordinates: {
    get: function () {return this.uniforms.worldCoordinates;},
    set: function (value) {this.uniforms.worldCoordinates = value;}
  },
  worldSize: {
    get: function () {return this.uniforms.worldSize;},
    set: function (value) {this.uniforms.worldSize = value;}
  }
});

AnnulusFilter.prototype.apply = function (filterManager, input, output)
{
    filterManager.calculateNormalizedScreenSpaceMatrix(this.uniforms.magicMatrix);
    filterManager.applyFilter(this, input, output);
};

let shaderContent = require('raw!../shaders/annulus_fog.frag');
let fogShader = new AnnulusFilter();

document.body.appendChild(renderer.view);

const loadProgressHandler = (loader, resource) => {
  console.log(`progress: ${loader.progress}%`);
};

const spaceMessageHandler = (space) => {
  globalState.physicsFrameID = space.physicsFrameID;
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

const shotHandler = (message) => {
  let frameId = message.frameId;
  let origin = message.origin;
  let facing = message.facing;
  let range = message.range;

  console.log(frameId + " " + origin.x + " " + origin.y + " " + facing + " " + range);
};

registerMessageHandler(Constants.HELLO_MESSAGE,       helloMessageHandler);
registerMessageHandler(Constants.SPACE_MESSAGE,       spaceMessageHandler);
registerMessageHandler(Constants.PLAYER_LEFT_MESSAGE, playerLeftHandler);
registerMessageHandler(Constants.SHOT_MESSAGE,        shotHandler);

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

  overlay = new PIXI.Graphics();
  overlay.drawRect(0, 0, 10, 10);
  overlay.filterArea = new PIXI.Rectangle(0, 0, 800, 600);
  overlay.filters = [fogShader];
  stage.addChild(overlay);

  hudText = new PIXI.Text('', hudTextStyle);
  hudText.x = 580;
  hudText.y = 0;
  stage.addChild(hudText);

  main();
}

var viewport = {vx: 0, vy: 0, width: 800, height: 600}

var frameCounter = 0;
var lastTime = Date.now();
var fps = 0;

// Handle keyboard controls
const KEY_SPACE = 32;
const KEY_UP = 38;
const KEY_LEFT = 37;
const KEY_RIGHT = 39;

const keysDown = new Map();

keysDown.set(KEY_SPACE,    false);
keysDown.set(KEY_UP,    false);
keysDown.set(KEY_LEFT,  false);
keysDown.set(KEY_RIGHT, false);

addEventListener("keydown", function (e) {
  updateKeysState(e.keyCode, true);
}, false);

addEventListener("keyup", function (e) {
  updateKeysState(e.keyCode, false);
}, false);

const updateKeysState = (keyCode, isPressed) => {
  const lastState = keysDown.get(keyCode);
  if (lastState == undefined) {
    return;
  }

  if (lastState != isPressed) {
    keysDown.set(keyCode, isPressed);
    sendInput(keyCode, isPressed);
  }
}

var sendInput = (keyCode, isPressed) => {
  let userInput = "CENTER"

  switch(keyCode) {
    case KEY_UP:
      userInput = isPressed ? "THRUST_ON" : "THRUST_OFF";
      break;
    case KEY_LEFT:
      userInput = isPressed ? "LEFT" : "CENTER"
      break;
    case KEY_RIGHT:
      userInput = isPressed ? "RIGHT" : "CENTER"
      break;
    case KEY_SPACE:
      userInput = isPressed ? "FIRE_START" : "FIRE_STOP"
      break;
  }

  let userMessage = new UserMessage(userInput);

  sendMessage(userMessage);
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

  fogShader.worldCoordinates[0] = x;
  fogShader.worldCoordinates[1] = y;

  renderer.render(stage);
};

// The main game loop
var main = function () {
  render();
  // Request to do this again ASAP
  requestAnimationFrame(main);
};
