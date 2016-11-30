import { globalState, stage } from './globals.js';
import * as Utils from './utils';
import { RADER_REFRESH_TEXTURE } from './constants.js';

export default class Radar {
  constructor() {
    this.container = new PIXI.Container();
    this.radius = 70;
    this.borderOffset = 20;
    this.radarColor = 0x00FF00;
    this.enemyColor = 0xFF0000;
    this.maxRadarRange = 3000;
    this.radarScale = this.radius / this.maxRadarRange;
    this.shipsGraphics = new PIXI.Graphics();
    this.refreshSprite = new PIXI.Sprite(this.loadRefreshTexture());
  }

  show() {
    this.container.x = -200;
    this.container.y = -200;

    let backgroundSprite = new PIXI.Sprite(this.generateBackgroundTexture());
    this.container.addChild(backgroundSprite);

    this.container.addChild(this.refreshSprite);
    this.refreshSprite.setTransform(this.radius, this.radius);

    this.container.addChild(this.shipsGraphics);
    this.shipsGraphics.setTransform(this.radius, this.radius);
    stage.addChild(this.container);
  }

  update(myShip, viewport) {
    let radarPos = Utils.translateToViewport(myShip.position.x/100, myShip.position.y/100, viewport);
    this.container.x = radarPos.x + (viewport.width / 2 - (2*this.radius)) - this.borderOffset;
    this.container.y = radarPos.y + (viewport.height / 2 - (2*this.radius)) - this.borderOffset;

    this.refreshSprite.rotation += 0.025;
    this.drawShips(this.otherShips(myShip, viewport));
  }

  drawShips(otherShips) {
    this.shipsGraphics.clear();
    for(let ship of otherShips) {
      if(this._withinCircle(ship.x, ship.y, this.maxRadarRange)) {
        this.drawShip(ship.x * this.radarScale, ship.y * this.radarScale);
      }
    }
  }

  _withinCircle(x, y, circleRadius) {
    return Math.pow(x,2) + Math.pow(y,2) < Math.pow(circleRadius, 2);
  }

  drawShip(x,y) {
    this.shipsGraphics.beginFill(this.enemyColor, 1);
    this.shipsGraphics.drawCircle(x, y, 3);
  }

  otherShips(myShip, viewport) {
    let result = [];
    for(let ship of globalState.spaceshipMap.values()) {
      if(ship.id !== myShip.id) {
        let viewportRespectCoords = Utils.translateToViewport(ship.position.x / 100, ship.position.y / 100, viewport);
        viewportRespectCoords.x -= viewport.width / 2;
        viewportRespectCoords.y -= viewport.height / 2;
        result.push(viewportRespectCoords);
      }
    }
    return result;
  }

  loadRefreshTexture() {
    return PIXI.Texture.fromImage(RADER_REFRESH_TEXTURE);
  }

  generateBackgroundTexture() {
    let radarGrahpics = new PIXI.Graphics();
    radarGrahpics.lineStyle(1, this.radarColor, 1);

    // inner circles
    radarGrahpics.drawCircle(0, 0, this.radius * 0.25);
    radarGrahpics.drawCircle(0, 0, this.radius * 0.5);
    radarGrahpics.drawCircle(0, 0, this.radius * 0.75);

    // inner lines
    radarGrahpics.lineStyle(1, this.radarColor, 1);
    radarGrahpics.moveTo(-this.radius, 0);
    radarGrahpics.lineTo(this.radius, 0);
    radarGrahpics.moveTo(0, -this.radius);
    radarGrahpics.lineTo(0, this.radius);
    radarGrahpics.moveTo(this.radius * Math.cos(Math.PI * 0.75), this.radius * Math.sin(Math.PI * 0.75));
    radarGrahpics.lineTo(this.radius * Math.cos(-Math.PI * 0.25), this.radius * Math.sin(-Math.PI * 0.25));
    radarGrahpics.moveTo(this.radius * Math.cos(Math.PI * 0.25), this.radius * Math.sin(Math.PI * 0.25));
    radarGrahpics.lineTo(this.radius * Math.cos(Math.PI * 1.25), this.radius * Math.sin(Math.PI * 1.25));
    radarGrahpics.moveTo(0, 0);
    // background
    radarGrahpics.lineStyle(1.5, 0x00FF00, 1);
    radarGrahpics.beginFill(0x00FF00, 0.1);
    radarGrahpics.drawCircle(0, 0, this.radius);
    radarGrahpics.endFill();

    return radarGrahpics.generateTexture();
  }
}
