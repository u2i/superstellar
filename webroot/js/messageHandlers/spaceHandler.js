import { globalState } from '../globals';
import Spaceship from '../spaceship';
import Assets from '../assets';
import * as Constants from '../constants';

const spaceHandler = (space) => {
  globalState.physicsFrameID = space.physicsFrameID;
  const ships = space.spaceships;
  const shipTexture = Assets.getTexture(Constants.SHIP_TEXTURE);

  let shipThrustFrames = [];

  Constants.FLAME_SPRITESHEET_FRAME_NAMES.forEach((frameName) =>  {
    shipThrustFrames.push(Assets.getTextureFromFrame(frameName));
  });

  for (let i in ships) {
    const shipId = ships[i].id;

    if (!globalState.spaceshipMap.has(shipId)) {
      const newSpaceship = new Spaceship(shipTexture, shipThrustFrames, ships[i]);

      globalState.spaceshipMap.set(shipId, newSpaceship);
    } else {
      globalState.spaceshipMap.get(shipId).updateData(ships[i]);
    }
  }
};

export default spaceHandler;
