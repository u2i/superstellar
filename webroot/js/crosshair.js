import { globalState, stage, renderer } from './globals.js';
import * as Utils from './utils';
import { CROSSHAIR_TEXTURE } from './constants.js';
import Assets from './assets';


export default class Crosshair {
  constructor() {
    this.container = new PIXI.Container();
  }

  show() {
    this.container.x = 0;
    this.container.y = 0;

    let crosshairSprite = new PIXI.Sprite(Assets.getTexture(CROSSHAIR_TEXTURE));
    this.container.addChild(crosshairSprite);
    stage.addChild(this.container);
  }

  update(x, y) {
    let relX = x - renderer.width / 2;
    let relY = y - renderer.height / 2;
    let targetAngle = Math.atan2(relY, relX);

//    this.container.x = renderer.width / 2 + Math.cos(targetAngle) * 150 - 16;
//    this.container.y = renderer.height / 2 + Math.sin(targetAngle) * 150 - 16;

    this.container.x = x - 16;
    this.container.y = y - 16;
  }
}
