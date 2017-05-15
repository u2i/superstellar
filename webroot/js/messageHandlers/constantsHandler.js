import { globalState, constants } from '../globals';

const constantsHandler = (message) => {
  Object.assign(constants, message);

  const anulusBorder = constants.worldRadius + 2 * constants.boundaryAnnulusWidth;
  globalState.worldSizeFilter.worldBoundarySize = new Float32Array([constants.worldRadius / 100, anulusBorder / 100]);
  globalState.constants = constants;
};

export default constantsHandler;
