import { stage } from './globals.js';
import { ASTEROID_01_TEXTURE } from './constants.js';
import * as Utils from './utils.js';

export default class Asteroid {
  constructor (id) {
    this.id = id;

    this.sprite = new PIXI.Sprite(PIXI.Texture.fromImage(ASTEROID_01_TEXTURE));

    this.position = new PIXI.Point();
    this.position.x = 200;
    this.position.y = 200;

    this.container = new PIXI.Container();
    this.container.position = this.position
    this.container.addChild(this.sprite);

    stage.addChild(this.container);
  }

  update (viewport) {
    const translatedPosition = Utils.translateToViewport(
      this.position.x + 1 / 100,
      this.position.y / 100,
      viewport
    )
    this.container.position = translatedPosition;
  }

  remove () {
    stage.removeChild(this.container);
  }
}

