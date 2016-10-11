import * as Constants from './constants';
import { registerMessageHandler } from './communicationLayer';

import playerLeftHandler from './messageHandlers/playerLeftHandler';
import helloHandler from './messageHandlers/helloHandler';
import spaceHandler from './messageHandlers/spaceHandler';
import projectileFiredHandler from './messageHandlers/projectileFiredHandler';


export const initializeHandlers = () => {
  registerMessageHandler(Constants.HELLO_MESSAGE,            helloHandler);
  registerMessageHandler(Constants.SPACE_MESSAGE,            spaceHandler);
  registerMessageHandler(Constants.PLAYER_LEFT_MESSAGE,      playerLeftHandler);
  registerMessageHandler(Constants.PROJECTILE_FIRED_MESSAGE, projectileFiredHandler);
};

