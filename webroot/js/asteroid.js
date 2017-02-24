import { stage, constants } from './globals.js';
import { ASTEROID_01_TEXTURE } from './constants.js'
import * as Utils from './utils.js';
import Victor from 'victor';

Victor.prototype.scalarMultiply = function(scalar) {
  return this.multiply(new Victor(scalar, scalar));
}

export default class Asteroid {
  constructor (id) {
    this.id = id;

    this.sprite = new PIXI.Sprite(PIXI.Texture.fromImage(ASTEROID_01_TEXTURE));

    this.position = new Victor(200, 200);
    this.velocity = new Victor(900, 0);

    this.container = new PIXI.Container();
    this.container.position = this.position
    this.container.addChild(this.sprite);

    stage.addChild(this.container);
  }

  update (viewport) {
    this._updatePosition();

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

  _updatePosition() {
    this._applyAnnulus()
    this.position = new Victor(this.position.x + this.velocity.x, this.position.y + this.velocity.y);
  }

  _applyAnnulus() {
    if (this.position.length() > constants.worldRadius) {
      let outreachLength = this.position.length() - constants.worldRadius;
      let gravityAcceleration = -(outreachLength / constants.boundaryAnnulusWidth) * constants.spaceshipAcceleration;
      let deltaVelocity = this.position.clone().normalize().scalarMultiply(gravityAcceleration);
      this.velocity.add(deltaVelocity);
    }
  }
}

