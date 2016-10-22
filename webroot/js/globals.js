import UsernameDialog from './usernameDialog';
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
  spaceshipMap: new Map(),
  physicsFrameID: 0,
  projectiles: []
};

export const usernameDialog = new UsernameDialog();
export const leaderboardDialog = new LeaderboardDialog();
