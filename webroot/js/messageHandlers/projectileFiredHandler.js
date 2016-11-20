import * as Constants from '../constants';
import Assets from '../assets';
import { globalState } from '../globals';
import Projectile from '../projectile';

const projectileFiredHandler = (message) => {
  let { id, frameId, origin, ttl, velocity, facing } = message;

  let animationFrames = [];

  Constants.PROJECTILE_SPRITESHEET_FRAME_NAMES.forEach((frameName) => {
    animationFrames.push(Assets.getTextureFromFrame(frameName));
  });

  globalState.projectilesMap.set(id, new Projectile(id, animationFrames, frameId, origin, ttl, velocity, facing));
};

export default projectileFiredHandler;
