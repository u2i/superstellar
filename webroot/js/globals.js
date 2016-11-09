import { getCurrentWindowSize } from './utils';
import LeaderboardDialog from "./leaderboardDialog";

const windowSize = getCurrentWindowSize();

export const renderer = new PIXI.WebGLRenderer(
  windowSize.width, windowSize.height, {autoResize: true}
);
export const stage = new PIXI.Container();

export const globalState = {
  clientId: null,
  clientIdToName: new Map(),
  nickname: null,
  spaceshipMap: new Map(),
  physicsFrameID: 0,
  projectiles: [],
  dialog: null
};

export const leaderboardDialog = new LeaderboardDialog();
