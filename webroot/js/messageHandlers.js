import * as Constants from './constants';
import { registerMessageHandler } from './communicationLayer';

import playerLeftHandler from './messageHandlers/playerLeftHandler';
import helloHandler from './messageHandlers/helloHandler';
import joinGameAckHandler from './messageHandlers/joinGameAckHandler';
import spaceHandler from './messageHandlers/spaceHandler';
import projectileFiredHandler from './messageHandlers/projectileFiredHandler';
import playerJoinedHandler from './messageHandlers/playerJoinedHandler';
import leaderboardHandler from './messageHandlers/leaderboardHandler';
import playerDiedHandler from './messageHandlers/playerDiedHandler';
import pongHandler from './messageHandlers/pongHandler';

export const initializeHandlers = () => {
  registerMessageHandler(Constants.HELLO_MESSAGE,            helloHandler);
  registerMessageHandler(Constants.JOIN_GAME_ACK_MESSAGE,    joinGameAckHandler);
  registerMessageHandler(Constants.PLAYER_LEFT_MESSAGE,      playerLeftHandler);
  registerMessageHandler(Constants.PLAYER_JOINED_MESSAGE,    playerJoinedHandler);
  registerMessageHandler(Constants.PROJECTILE_FIRED_MESSAGE, projectileFiredHandler);
  registerMessageHandler(Constants.SPACE_MESSAGE,            spaceHandler);
  registerMessageHandler(Constants.LEADERBOARD_MESSAGE,      leaderboardHandler);
  registerMessageHandler(Constants.PLAYER_DIED_MESSAGE,      playerDiedHandler);
  registerMessageHandler(Constants.PONG_MESSAGE,             pongHandler);
};
