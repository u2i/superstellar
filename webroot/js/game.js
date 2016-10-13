import * as PIXI from "pixi.js";
import Assets from './assets';
import * as Constants from './constants';
import * as Utils from './utils';
import { renderer, stage, globalState } from './globals';
import { initializeConnection, sendMessage, UserMessage } from './communicationLayer';
import { initializeHandlers } from './messageHandlers';
import { initializeControls } from './controls';
import Hud from './hud';

const HOST = BACKEND_HOST;
const PORT = BACKEND_PORT;
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

PIXI.loader.
  add([Constants.SHIP_TEXTURE, Constants.BACKGROUND_TEXTURE, Constants.FLAME_SPRITESHEET, Constants.PROJECTILE_SPRITESHEET]).
  on("progress", loadProgressHandler).
  load(setup);

let tilingSprite;

let hud;

function setup() {
  initializeHandlers();
  initializeConnection(HOST, PORT, PATH);
  initializeControls();

  const bgTexture = Assets.getTexture(Constants.BACKGROUND_TEXTURE);

  tilingSprite = new PIXI.extras.TilingSprite(bgTexture, renderer.width, renderer.height);
  stage.addChild(tilingSprite);

  overlay = new PIXI.Graphics();
  overlay.drawRect(0, 0, 10, 10);
  overlay.filterArea = new PIXI.Rectangle(0, 0, renderer.width, renderer.width);
  overlay.filters = [fogShader];
  stage.addChild(overlay);

  hud = new Hud();

  main();
}

window.addEventListener("resize", () => {
  Utils.getCurrentWindowSize((width, height) => {
    renderer.resize(width, height);
    tilingSprite.width = width;
    tilingSprite.height = height;
    overlay.filterArea.width = width;
    overlay.filterArea.height = height;
  });
});

var viewport = {vx: 0, vy: 0, width: 800, height: 600}

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
  globalState.projectiles.forEach((projectile) => projectile.update(viewport));

  hud.update();

  let x = myShip ? Math.floor(myShip.position.x / 100) : '?';
  let y = myShip ? Math.floor(myShip.position.y / 100) : '?';

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
