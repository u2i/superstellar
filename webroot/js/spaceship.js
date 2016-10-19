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

    if (__DEBUG__) {
      this.collisionSphere = new PIXI.Graphics();
      this.collisionSphere.beginFill(0xFF77FF);
      this.collisionSphere.alpha = 0.3;
      this.collisionSphere.drawCircle(this.sprite.width / 2, this.sprite.height / 2, 20);
    }

    this.hpTextStyle = {
          fontFamily: 'Helvetica',
          fontSize: '24px',
          fill: '#FFFFFF',
          align: 'left'
        };

    this.healthBar = new PIXI.Text('', this.hpTextStyle);
	this.healthBar.y = -40

    stage.addChild(this.container);
    this.container.addChild(this.sprite);
    this.container.addChild(this.thrustAnimation);

    if (__DEBUG__) {
      this.container.addChild(this.collisionSphere);
    }

	this.container.addChild(this.healthBar)

    this.container.pivot.set(this.sprite.width / 2, this.sprite.height / 2);
  }

  updateData ({ id, position, velocity, facing, inputThrust, hp, maxHp }) {
    this.id = id;
    this.position = position;
    this.velocity = velocity;
    this.facing = facing;
    this.inputThrust = inputThrust;
    this.hp = hp;
    this.maxHp = maxHp;
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

    this.healthBar.text = this.hp;
  }

  remove () {
    stage.removeChild(this.container);
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
