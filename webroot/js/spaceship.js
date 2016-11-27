import * as Utils from './utils.js';
import HealthBarFilter from './healthBarFilter';
import {globalState, renderer, stage} from './globals.js';

const healthBarRadius = 40;

export default class Spaceship {
  constructor(shipTexture, thrustAnimationFrames) {
    this.createHealthBarFilter();
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

    this.labelTextStyle = {
      fontFamily: 'Roboto',
      fontSize: '12px',
      fill: '#FFFFFF',
      align: 'center',
    };

    this.label = new PIXI.Text('', this.labelTextStyle);

    stage.addChild(this.container);
    this.container.addChild(this.sprite);
    this.container.addChild(this.thrustAnimation);
    this.addHealthBar();

    if (__DEBUG__) {
      this.container.addChild(this.collisionSphere);
    }

    stage.addChild(this.label);

    this.container.pivot.set(this.sprite.width / 2, this.sprite.height / 2);
  }

  updateData(timestamp, {id, position, velocity, facing, inputThrust, hp, maxHp}) {
    this.timestamp = timestamp;
    this.id = id;
    this.position = position;
    this.velocity = velocity;
    this.facing = facing;
    this.inputThrust = inputThrust;
    this.hp = hp;
    this.maxHp = maxHp;
    this.updateHealthBar();
  }

  interpolateData() {
    let now = new Date();
    let delta = now - this.timestamp;

    let interpolatedPositionX = this.position.x + this.velocity.x * delta / globalState.physicsFrameRate;
    let interpolatedPositionY = this.position.y + this.velocity.y * delta / globalState.physicsFrameRate;

    //console.log(now);
    //console.log(interpolatedPositionX);
    //console.log(this.velocity.x);
    //console.log(this.position.x);
    //console.log(this.velocity.x * delta / globalState.physicsFrameRate);
    //console.log(globalState.physicsFrameID);
    console.log('--------------')

    this.interpolatedPosition = {x: Math.round(interpolatedPositionX), y: Math.round(interpolatedPositionY)};
  }

  update(viewport) {
    if (this.inputThrust) {
      this.thrustAnimation.visible = true;
      this.thrustAnimation.play();
    } else {
      this.thrustAnimation.visible = false;
      this.thrustAnimation.stop();
    }

    const {x, y} = Utils.translateToViewport(
      this.interpolatedPosition.x / 100,
      this.interpolatedPosition.y / 100,
      viewport
    );

    this.container.position.set(x, y);

    if (this.isOutOfView(x, y, viewport)) {
      this.disableHealthBarFilter();
    } else {
      this.enableHealthBarFilter(x, y);
    }

    this.container.rotation = this.facing;

    if (globalState.clientId !== this.id) {
      this.label.text = globalState.clientIdToName.get(this.id);

      if (this.id === globalState.killedBy) {
        this.label.style.fill = '#FF0000'
      }

      this.label.position.set(x - (this.label.text.length * 6) / 2, y + this.sprite.height);
    }
  }

  isOutOfView(x, y, viewport) {
    return x - healthBarRadius < 0
      || y - healthBarRadius < 0
      || x + healthBarRadius > viewport.width
      || y + healthBarRadius > viewport.height;
  }

  addHealthBar() {
    this.healthBar = new PIXI.Graphics();
    this.healthBarRectangle = new PIXI.Rectangle(100, 100, healthBarRadius * 2, healthBarRadius * 2);
    this.healthBar.filterArea = this.healthBarRectangle;
    this.container.addChild(this.healthBar);
  }

  enableHealthBarFilter(x, y) {
    this.healthBarRectangle.x = x - healthBarRadius;
    this.healthBarRectangle.y = y - healthBarRadius;
    this.healthBar.filters = [this.healthBarFilter];
  }

  disableHealthBarFilter() {
    this.healthBar.filters = [];
  }

  createHealthBarFilter() {
    this.healthBarFilter = new HealthBarFilter();
  }

  updateHealthBar() {
    this.healthBarFilter.hps = [this.hp, this.maxHp];
  }

  remove() {
    stage.removeChild(this.container);
    stage.removeChild(this.label);
  }

  viewport() {
    return {
      vx: this.interpolatedPosition.x / 100,
      vy: this.interpolatedPosition.y / 100,
      width: renderer.width,
      height: renderer.height
    };
  }
}
