import * as Utils from './utils.js';
import { renderer, stage } from './globals.js';

export default class Spaceship {
  constructor (shipTexture, thrustAnimationFrames, data) {
    this.updateData(data);
    this.container = new PIXI.Container();
    this.sprite = new PIXI.Sprite(shipTexture);
    this.thrustAnimation = new PIXI.extras.MovieClip(thrustAnimationFrames);

    this.thrustAnimation.position.set(-27, 7);
    this.thrustAnimation.animationSpeed = 0.5;
    
    stage.addChild(this.container);
    this.container.addChild(this.sprite);
    this.container.addChild(this.thrustAnimation);
    this.container.pivot.set(this.sprite.width / 2, this.sprite.height / 2);
  }

  updateData ({ id, position, velocity, facing, inputThrust }) {
    this.id = id;
    this.position = position;
    this.velocity = velocity;
    this.facing = facing;
    this.inputThrust = inputThrust;
  }

  update (viewport) {
    if (this.inputThrust) {
      this.thrustAnimation.visible = true;
      this.thrustAnimation.play();
    } else {
      this.thrustAnimation.visible = false;
      this.thrustAnimation.stop();
    }

    const translatedPosition = Utils.translateToViewport(
      this.position.x / 100, 
      this.position.y / 100, 
      viewport
    )

    this.container.position.set(translatedPosition.x, translatedPosition.y);
    this.container.rotation = this.facing;
  }

  viewport () {
    return {
      vx: this.position.x / 100,
      vy: this.position.y / 100,
      width: renderer.width,
      height: renderer.height
    };
  }
};
