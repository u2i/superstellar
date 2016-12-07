import { stage, globalState } from './globals';
import * as Utils from './utils.js';

export default class Projectile {
  constructor (id, animationFrames, frameId, origin, ttl, velocity, facing) {
    this.id = id;
    this.frameId  = frameId;
    this.origin   = origin;
    this.ttl      = ttl;
    this.velocity = velocity;
    this.facing   = facing;

    this.animation = new PIXI.extras.MovieClip(animationFrames);

    this.position = new PIXI.Point();

    this._updatePosition();

    this.animation.position.set(this.position.x / 100, this.position.y / 100);
    this.animation.rotation = this.facing;
    this.animation.animationSpeed = 10;
    this.animation.pivot.set(0, this.animation.height / 2);
    this.animation.play();

    stage.addChild(this.animation);
  }

  update (viewport, currentFrameId) {
    if (currentFrameId > this.frameId + this.ttl) {
      this.remove();
      globalState.projectilesMap.delete(this.id);
    }

    this._updatePosition(currentFrameId);
    this.animation.play();

    const translatedPosition = Utils.translateToViewport(
      this.position.x / 100,
      this.position.y / 100,
      viewport
    )

    this.animation.position.set(translatedPosition.x, translatedPosition.y);
  }

  remove () {
    stage.removeChild(this.animation);
  }

  _updatePosition(currentFrameId) {
    const frameOffset = currentFrameId - this.frameId;

    this.position.set(this.origin.x + this.velocity.x * frameOffset, this.origin.y + this.velocity.y * frameOffset);
  }
}

