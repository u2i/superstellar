import * as PIXI from "pixi.js";
import {globalState, stage} from './globals';

const rightOffset = 150;

export default class Hud {
  static get rightOffset() {
    return rightOffset;
  }

  constructor() {
    this.domNode = document.getElementById('hud-debug');

    this.frameCounter = 0;
    this.fps = 0;
    this.lastTime = Date.now();
  }

  show() {
    this.domNode.style.display = 'block'
  }

  hide() {
    this.domNode.style.display = 'none'
  }

  update() {
    this.frameCounter++;

    if (this.frameCounter === 100) {
      this.frameCounter = 0;
      const now = Date.now();
      const delta = (now - this.lastTime) / 1000;
      this.fps = (100 / delta).toFixed(1);
      this.lastTime = now;
    }

    this._updateHtml()
  }

  _updateHtml() {
    const playerShip = globalState.spaceshipMap.get(globalState.clientId);

    this.getFpsNode().innerHTML = this.fps;
    this.getShipsNode().innerHTML = globalState.spaceshipMap.size;
    this.getSpeedNode().innerHTML = this.getSpeed();

    if (__DEBUG__) {
      let x = playerShip ? Math.floor(playerShip.position.x / 100) : '?';
      let y = playerShip ? Math.floor(playerShip.position.y / 100) : '?';

      text += "X: " + x + "\n";
      text += "Y: " + y + "\n";
    }
  }

  getSpeed() {
    let currentShip = globalState.spaceshipMap.get(globalState.clientId);
    if (currentShip !== undefined) {
      let x2 = currentShip.velocity.x * currentShip.velocity.x;
      var y2 = currentShip.velocity.y * currentShip.velocity.y;
      let v = Math.sqrt(x2 + y2);
      return v.toFixed(0)
    }
    return 'n/a';
  }

  getFpsNode() {
    return this.domNode.querySelector('[data-type="fps"]');
  }

  getShipsNode() {
    return this.domNode.querySelector('[data-type="ships"]');
  }

  getSpeedNode() {
    return this.domNode.querySelector('[data-type="speed"]');
  }
}
