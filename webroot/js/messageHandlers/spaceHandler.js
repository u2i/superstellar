import { globalState } from '../globals';
import Spaceship from '../spaceship';
import Assets from '../assets';
import * as Constants from '../constants';

const spaceHandler = (space) => {
  const ships = space.spaceships;
  const shipTexture = Assets.getTexture(Constants.SHIP_TEXTURE);

  globalState.framesCalculator.receivedFrameId(space.physicsFrameID);

  let shipThrustFrames = [];

  Constants.FLAME_SPRITESHEET_FRAME_NAMES.forEach((frameName) =>  {
    shipThrustFrames.push(Assets.getTextureFromFrame(frameName));
  });

  for (let i in ships) {
    const shipId = ships[i].id;

    if (!globalState.spaceshipMap.has(shipId)) {
      const newSpaceship = new Spaceship(shipTexture, shipThrustFrames, space.physicsFrameID, ships[i]);

      globalState.spaceshipMap.set(shipId, newSpaceship);
    }
    globalState.spaceshipMap.get(shipId).updateData(space.physicsFrameID, ships[i]);
  }
};

export default spaceHandler;
