import { globalState } from '../globals';

const helloHandler = (message) => {
  const { myId, idToUsername, worldRadius, boundaryAnnulusWidth, firstPhysicsFrameTimestamp, physicsFrameRate } = message;
  globalState.clientId = myId;
  idToUsername.forEach((username, id) => {
    globalState.clientIdToName.set(id, username);
  });
  const anulusBorder = worldRadius + 2*boundaryAnnulusWidth;
  globalState.worldSizeFilter.worldBoundarySize = new Float32Array([worldRadius, anulusBorder]);
  global.firstPhysicsFrameTimestamp = firstPhysicsFrameTimestamp;
  global.physicsFrameRate = physicsFrameRate;
};

export default helloHandler;
