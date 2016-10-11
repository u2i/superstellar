import * as Constants from '../constants';
import Assets from '../assets';
import { globalState } from '../globals';
import Projectile from '../projectile';

const projectileFiredHandler = (message) => {
  let { frameId, origin, ttl, velocity } = message;

  let animationFrames = [];

  Constants.PROJECTILE_SPRITESHEET_FRAME_NAMES.forEach((frameName) => {
    animationFrames.push(Assets.getTextureFromFrame(frameName));
  });

  globalState.projectiles.push(new Projectile(animationFrames, frameId, origin, ttl, velocity));
};

export default projectileFiredHandler;
