require("../css/superstellar.scss");

import * as PIXI from "pixi.js";
import Assets from './assets';
import * as Constants from './constants';
import * as Utils from './utils';
import { renderer, stage, globalState, leaderboardDialog } from './globals';
import { initializeConnection } from './communicationLayer';
import { initializeHandlers } from './messageHandlers';
import Hud from './hud';
import UsernameDialog from "./dialogs/usernameDialog";
import Radar from './radar';

const HOST = window.location.hostname;
const PORT = BACKEND_PORT;
const PATH = '/superstellar';

const fogShader = globalState.worldSizeFilter;

document.getElementById('game').appendChild(renderer.view);

const loadProgressHandler = (loader) => {
  console.log(`progress: ${loader.progress}%`);
};

PIXI.loader.
  add([Constants.SHIP_TEXTURE, Constants.BACKGROUND_TEXTURE, Constants.FLAME_SPRITESHEET, Constants.PROJECTILE_SPRITESHEET]).
  on("progress", loadProgressHandler).
  load(setup);

let tilingSprite;
let overlay;
let hud;
let radar;

function setup() {
  initializeHandlers();
  initializeConnection(HOST, PORT, PATH);


  const bgTexture = Assets.getTexture(Constants.BACKGROUND_TEXTURE);

  tilingSprite = new PIXI.extras.TilingSprite(bgTexture, renderer.width, renderer.height);
  stage.addChild(tilingSprite);

  overlay = new PIXI.Graphics();
  overlay.drawRect(0, 0, 10, 10);
  overlay.filterArea = new PIXI.Rectangle(0, 0, renderer.width, renderer.height);
  overlay.filters = [fogShader];
  stage.addChild(overlay);

  radar = new Radar();
  radar.show();

  hud = new Hud();
  hud.show();

  const dialog = new UsernameDialog();
  dialog.show()

  leaderboardDialog.show();

  main();
}

window.addEventListener("resize", () => {
  const windowSize = Utils.getCurrentWindowSize();
  renderer.resize(windowSize.width, windowSize.height);
  tilingSprite.width = windowSize.width;
  tilingSprite.height = windowSize.height;
  overlay.filterArea.width = windowSize.width;
  overlay.filterArea.height = windowSize.height;
});

const defaultViewport = { vx: 0, vy: 0, width: renderer.width, height: renderer.height };

// Draw everything
const render = function () {
  let myShip;

  for (var spaceship of globalState.spaceshipMap.values()) {
    spaceship.interpolateData();
  }

  if (globalState.spaceshipMap.size > 0) {
    myShip = globalState.spaceshipMap.get(globalState.clientId);
  }

  let viewport = myShip ? myShip.viewport() : defaultViewport;

  let backgroundPos = Utils.translateToViewport(0, 0, viewport);
  tilingSprite.tilePosition.set(backgroundPos.x, backgroundPos.y);

  globalState.spaceshipMap.forEach((spaceship) => spaceship.update(viewport));
  globalState.projectilesMap.forEach((projectile) => projectile.update(viewport));

  if(myShip) {
    radar.update(myShip, viewport);
  }

  hud.update();

  let x = myShip ? Math.floor(myShip.interpolatedPosition.x / 100) : '?';
  let y = myShip ? Math.floor(myShip.interpolatedPosition.y / 100) : '?';

  fogShader.worldCoordinates[0] = x;
  fogShader.worldCoordinates[1] = y;

  renderer.render(stage);
};

// The main game loop
const main = function () {
  render();
  // Request to do this again ASAP
  requestAnimationFrame(main);
};
