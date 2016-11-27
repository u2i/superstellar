import { globalState } from '../globals';
import Spaceship from '../spaceship';
import Assets from '../assets';
import * as Constants from '../constants';

const spaceHandler = (space) => {
  console.log(space.physicsFrameID - globalState.physicsFrameID)
  globalState.physicsFrameID = space.physicsFrameID;
  const ships = space.spaceships;
  const shipTexture = Assets.getTexture(Constants.SHIP_TEXTURE);

  let shipThrustFrames = [];

  Constants.FLAME_SPRITESHEET_FRAME_NAMES.forEach((frameName) =>  {
    shipThrustFrames.push(Assets.getTextureFromFrame(frameName));
  });

  for (let i in ships) {
    let shipId = ships[i].id;
    let timestamp = new Date();

    if (!globalState.spaceshipMap.has(shipId)) {
      const newSpaceship = new Spaceship(shipTexture, shipThrustFrames);
      globalState.spaceshipMap.set(shipId, newSpaceship);
    }

    globalState.spaceshipMap.get(shipId).updateData(timestamp, ships[i]);
  }
};

export default spaceHandler;
