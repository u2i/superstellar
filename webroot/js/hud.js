import * as PIXI from "pixi.js";
import { globalState, stage, hudRightOffset } from './globals';

export default class Hud {
  constructor (canvasWidth) {
    this.hudTextStyle = {
      fontFamily: 'Helvetica',
      fontSize: '24px',
      fill: '#FFFFFF',
      align: 'left',
      textBaseline: 'top'
    };

    this.text = new PIXI.Text('', this.hudTextStyle);
    this.text.x = canvasWidth - hudRightOffset;
    this.text.y = 0;

    this.frameCounter = 0;
    this.fps = 0;
    this.lastTime = Date.now();

    stage.addChild(this.text);
  }

  update () {
    this.frameCounter++;

    if (this.frameCounter === 100) {
      this.frameCounter = 0;
      const now = Date.now();
      const delta = (now - this.lastTime) / 1000;
      this.fps = (100 / delta).toFixed(1);
      this.lastTime = now;
    }

    this.text.text = this._updateHudText();
  }

  _updateHudText () {
    const playerShip = globalState.spaceshipMap.get(globalState.clientId);

    let text = "Ships: " + globalState.spaceshipMap.size + "\n";

    let x = playerShip ? Math.floor(playerShip.position.x / 100) : '?';
    let y = playerShip ? Math.floor(playerShip.position.y / 100) : '?';

    text += "FPS: " + this.fps + "\n";
    text += "X: " + x + "\n";
    text += "Y: " + y + "\n";

    return text;
  }
}
