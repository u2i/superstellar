import * as Utils from './utils.js';
import CircleBarFilter from './circleBarFilter';
import {globalState, renderer, stage} from './globals.js';
import SimulationFrame from './simulationFrame.js';

const healthBarRadius = 40;
const energyBarRadius = 50;

export default class Spaceship {
  constructor(shipTexture, thrustAnimationFrames, data) {
    this.createHealthBarFilter();
    this.createEnergyBarFilter();

    this.simulationFrame = new SimulationFrame(0, data);

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
    this.addEnergyBar();

    if (__DEBUG__) {
      this.container.addChild(this.collisionSphere);
    }

    stage.addChild(this.label);

    this.container.pivot.set(this.sprite.width / 2, this.sprite.height / 2);
  }

  updateData(data) {
    this.simulationFrame.update(data)
    this.position = this.simulationFrame.position;
    this.facing = this.simulationFrame.facing;

    this.id = data.id;
    this.hp = data.hp;
    this.maxHp = data.maxHp;
    this.energy = data.energy;
    this.maxEnergy = data.maxEnergy;
    this.updateHealthBar();
    this.updateEnergyBar();
  }

  predict() {
    this.simulationFrame.predict();
    this.position = this.simulationFrame.position;
    this.facing = this.simulationFrame.facing;
  }

  update(viewport) {
    if (this.simulationFrame.inputThrust) {
      this.thrustAnimation.visible = true;
      this.thrustAnimation.play();
    } else {
      this.thrustAnimation.visible = false;
      this.thrustAnimation.stop();
    }

    const {x, y} = Utils.translateToViewport(
      this.position.x / 100,
      this.position.y / 100,
      viewport
    );

    this.container.position.set(x, y);

    if (this.isOutOfView(x, y, healthBarRadius, viewport)) {
      this.disableHealthBarFilter();
    } else {
      this.enableHealthBarFilter(x, y);
    }

    if (this.isOutOfView(x, y, energyBarRadius, viewport)) {
      this.disableEnergyBarFilter();
    } else {
      this.enableEnergyBarFilter(x, y);
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

  isOutOfView(x, y, barRadius, viewport) {
    return x - barRadius < 0
      || y - barRadius < 0
      || x + barRadius > viewport.width
      || y + barRadius > viewport.height;
  }

  addHealthBar() {
    this.healthBar = new PIXI.Graphics();
    this.healthBarRectangle = new PIXI.Rectangle(100, 100, healthBarRadius * 2, healthBarRadius * 2);
    this.healthBar.filterArea = this.healthBarRectangle;
    this.container.addChild(this.healthBar);
  }

  addEnergyBar() {
    this.energyBar = new PIXI.Graphics();
    this.energyBarRectangle = new PIXI.Rectangle(100, 100, energyBarRadius * 2, energyBarRadius * 2);
    this.energyBar.filterArea = this.energyBarRectangle;
    this.container.addChild(this.energyBar);
  }

  enableHealthBarFilter(x, y) {
    this.healthBarRectangle.x = x - healthBarRadius;
    this.healthBarRectangle.y = y - healthBarRadius;
    this.healthBar.filters = [this.healthBarFilter];
  }

  enableEnergyBarFilter(x, y) {
    this.energyBarRectangle.x = x - energyBarRadius;
    this.energyBarRectangle.y = y - energyBarRadius;
    this.energyBar.filters = [this.energyBarFilter];
  }

  disableHealthBarFilter() {
    this.healthBar.filters = [];
  }

  disableEnergyBarFilter() {
    this.energyBar.filters = [];
  }

  createHealthBarFilter() {
    this.healthBarFilter = new CircleBarFilter([0.6, 1.0, 0.6]);
  }

  createEnergyBarFilter() {
    this.energyBarFilter = new CircleBarFilter([0.6, 0.6, 1.0]);
  }

  updateHealthBar() {
    this.healthBarFilter.hps = [this.hp, this.maxHp];
  }

  updateEnergyBar() {
    this.energyBarFilter.hps = [this.energy, this.maxEnergy]
  }

  remove() {
    stage.removeChild(this.container);
    stage.removeChild(this.label);
  }

  viewport() {
    return {
      vx: this.position.x / 100,
      vy: this.position.y / 100,
      width: renderer.width,
      height: renderer.height
    };
  }
}
