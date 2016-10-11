import { renderer, stage, globalState } from './globals';
import * as Utils from './utils.js';
import * as Constants from './constants';

export default class Projectile {
  constructor (animationFrames, frameId, origin, facing, range) {
    this.frameId = frameId;
    this.origin  = origin;
    this.facing  = facing;
    this.range   = range;

    this.velocity = new PIXI.Point(Math.cos(facing) * 1000, Math.sin(facing) * (-1000));

    console.log("range", this.range);

    this.animation = new PIXI.extras.MovieClip(animationFrames);

    const frameOffset = frameId - globalState.physicsFrameID;

    this.position = new PIXI.Point();

    this._updatePosition();

    console.log("position", this.position);

    this.animation.position.set(this.position.x / 100, this.position.y / 100);
    this.animation.rotation = this.facing;
    this.animation.animationSpeed = 10;
    this.animation.play();

    stage.addChild(this.animation);
  }

  update (viewport) {
    this._updatePosition();
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

  _updatePosition() {
    const frameOffset = globalState.physicsFrameID - this.frameId;

    this.position.set(this.origin.x + this.velocity.x * frameOffset, this.origin.y + this.velocity.y * frameOffset);
  }
}

