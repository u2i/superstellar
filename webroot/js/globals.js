import { getCurrentWindowSize } from './utils';

export const renderer = getCurrentWindowSize(
  (width, height) => new PIXI.WebGLRenderer(width, height, {autoResize: true})
);
export const stage = new PIXI.Container();

export const globalState = {
  clientId: null,
  spaceshipMap: new Map(),
  physicsFrameID: 0,
  projectiles: []
};
