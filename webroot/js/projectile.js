import { renderer, stage, globalState } from './globals';
import * as Utils from './utils.js';
import * as Constants from './constants';

export default class Projectile {
  constructor (animationFrames, frameId, origin, facing, ttl, speed) {
    this.frameId = frameId;
    this.origin  = origin;
    this.facing  = facing;
    this.ttl     = ttl;
    this.speed   = speed;


    this.velocity = new PIXI.Point(Math.cos(facing) * speed, Math.sin(facing) * (-speed));

    this.animation = new PIXI.extras.MovieClip(animationFrames);

    const frameOffset = frameId - globalState.physicsFrameID;

    this.position = new PIXI.Point();

    this._updatePosition();

    this.animation.position.set(this.position.x / 100, this.position.y / 100);
    this.animation.rotation = this.facing;
    this.animation.animationSpeed = 10;
    this.animation.play();

    stage.addChild(this.animation);
  }

  update (viewport) {
    if(globalState.physicsFrameID > this.frameId + this.ttl) {
        this.remove();
        globalState.projectiles.splice(globalState.projectiles.indexOf(this), 1);
    }

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

