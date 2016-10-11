export const renderer = new PIXI.WebGLRenderer(800, 600);
export const stage = new PIXI.Container();

export const globalState = {
  clientId: null,
  spaceshipMap: new Map(),
  physicsFrameID: 0,
  projectiles: []
};
