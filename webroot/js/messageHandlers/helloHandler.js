import { globalState, constants } from '../globals';

const helloHandler = (message) => {
  globalState.clientId = message.myId;
  message.idToUsername.forEach((username, id) => {
    globalState.clientIdToName.set(id, username);
  });

  Object.assign(constants, message.constants);

  const anulusBorder = constants.worldRadius + 2 * constants.boundaryAnnulusWidth;
  globalState.worldSizeFilter.worldBoundarySize = new Float32Array([constants.worldRadius / 100, anulusBorder / 100]);
  globalState.constants = constants;
};

export default helloHandler;
