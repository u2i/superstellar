import { stage, constants } from './globals.js';
import { ASTEROID_01_TEXTURE } from './constants.js'
import * as Utils from './utils.js';
import AsteroidMoveFilter from './asteroidMoveFilter.js';

export default class Asteroid {
  constructor (frameId) {
    this.moveFilter = new AsteroidMoveFilter(frameId);

    this.sprite = new PIXI.Sprite(PIXI.Texture.fromImage(ASTEROID_01_TEXTURE));
    this.sprite.width = 50;
    this.sprite.height = 50;

    this.container = new PIXI.Container();
    this.container.addChild(this.sprite);

    stage.addChild(this.container);

    this.container.pivot.set(this.sprite.width / 2, this.sprite.height / 2);
  }

  updateData(updateFrameId, data) {
    this.moveFilter.update(updateFrameId, data);
    this.position = this.moveFilter.position();
    this.velocity = this.moveFilter.velocity();
    this.facing = this.moveFilter.facing();

    this.id = data.id;
  }

  predictTo(frameId) {
    this.moveFilter.predictTo(frameId);
    this.position = this.moveFilter.position();
    this.velocity = this.moveFilter.velocity();
    this.facing = this.moveFilter.facing();
  }

  update (viewport) {
    const translatedPosition = Utils.translateToViewport(
      this.position.x / 100,
      this.position.y / 100,
      viewport
    )
    this.container.position = translatedPosition;
  }

  remove () {
    stage.removeChild(this.container);
  }

  obsolete() {
    return this.position.length() > constants.worldRadius * 2
  }
}

