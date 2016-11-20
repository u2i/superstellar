import {globalState} from './globals';

export default class Hud {
  constructor() {
    this.domNode = document.getElementById('hud-debug');

    this.frameCounter = 0;
    this.fps = 0;
    this.ping = 0;
    this.lastTime = Date.now();
  }

  show() {
    this.domNode.style.display = 'block';
    if (__DEBUG__) {
      for (let definition of this.domNode.querySelectorAll('[class="debug"]')) {
        definition.className = '';
      }
    }
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

    this.getFpsNode().textContent = this.fps;
    this.getShipsNode().textContent = globalState.spaceshipMap.size;
    this.getPingNode().textContent = globalState.ping;

    if (__DEBUG__) {
      this.getSpeedNode().textContent = this.getSpeed();

      let [x, y] = this.getPosition(playerShip);
      this.getPositionNode('x').textContent = x;
      this.getPositionNode('y').textContent = y;
    }
  }

  getSpeed() {
    let currentShip = globalState.spaceshipMap.get(globalState.clientId);
    if (currentShip !== undefined) {
      let x2 = currentShip.velocity.x * currentShip.velocity.x;
      let y2 = currentShip.velocity.y * currentShip.velocity.y;
      let v = Math.sqrt(x2 + y2);
      return v.toFixed(0)
    }
    return '?';
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

  getPingNode() {
    return this.domNode.querySelector('[data-type="ping"]');
  }

  getPosition(playerShip) {
    let x = playerShip ? Math.floor(playerShip.position.x / 100) : '?';
    let y = playerShip ? Math.floor(playerShip.position.y / 100) : '?';

    return [x, y];
  }

  getPositionNode(dimension) {
    return this.domNode.querySelector(`[data-type="coordinate-${dimension}"]`);
  }
}
