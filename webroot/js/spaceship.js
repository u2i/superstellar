import * as Utils from './utils.js';
import { globalState, renderer, stage } from './globals.js';

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
      fontSize: '15px',
      fill: '#FFFFFF',
      align: 'center'
    };

    this.label = new PIXI.Text('', this.hpTextStyle);

    this.healthBar = new PIXI.Text('', this.hpTextStyle);
    this.healthBar.y = -30;

    stage.addChild(this.container);
    this.container.addChild(this.sprite);
    this.container.addChild(this.thrustAnimation);

    if (__DEBUG__) {
      this.container.addChild(this.collisionSphere);
    }

    stage.addChild(this.label);
    this.container.addChild(this.healthBar);

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

    const { x, y } = Utils.translateToViewport(
      this.position.x / 100,
      this.position.y / 100,
      viewport
    )

    this.container.position.set(x, y);
    this.container.rotation = this.facing;

    this.label.text = globalState.clientIdToName.get(this.id);
    this.label.position.set(x - (this.label.text.length * 8) / 2, y - this.sprite.height);

    this.healthBar.text = this.hp;
  }

  remove () {
    stage.removeChild(this.container);
    stage.removeChild(this.label);
  }

  viewport () {
    return {
      vx: this.position.x / 100,
      vy: this.position.y / 100,
      width: renderer.width,
      height: renderer.height
    };
  }
}
