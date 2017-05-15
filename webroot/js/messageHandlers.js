import * as Constants from './constants';
import { registerMessageHandler } from './communicationLayer';

import playerLeftHandler from './messageHandlers/playerLeftHandler';
import helloHandler from './messageHandlers/helloHandler';
import constantsHandler from './messageHandlers/constantsHandler';
import joinGameAckHandler from './messageHandlers/joinGameAckHandler';
import spaceHandler from './messageHandlers/spaceHandler';
import projectileFiredHandler from './messageHandlers/projectileFiredHandler';
import projectileHitHandler from './messageHandlers/projectileHitHandler';
import playerJoinedHandler from './messageHandlers/playerJoinedHandler';
import leaderboardHandler from './messageHandlers/leaderboardHandler';
import playerDiedHandler from './messageHandlers/playerDiedHandler';
import pongHandler from './messageHandlers/pongHandler';
import scoreBoardHandler from './messageHandlers/scoreBoardHandler';

export const initializeHandlers = () => {
  registerMessageHandler(Constants.HELLO_MESSAGE,            helloHandler);
  registerMessageHandler(Constants.CONSTANTS_MESSAGE,        constantsHandler);
  registerMessageHandler(Constants.JOIN_GAME_ACK_MESSAGE,    joinGameAckHandler);
  registerMessageHandler(Constants.PLAYER_LEFT_MESSAGE,      playerLeftHandler);
  registerMessageHandler(Constants.PLAYER_JOINED_MESSAGE,    playerJoinedHandler);
  registerMessageHandler(Constants.PROJECTILE_FIRED_MESSAGE, projectileFiredHandler);
  registerMessageHandler(Constants.PROJECTILE_HIT_MESSAGE,   projectileHitHandler);
  registerMessageHandler(Constants.SPACE_MESSAGE,            spaceHandler);
  registerMessageHandler(Constants.LEADERBOARD_MESSAGE,      leaderboardHandler);
  registerMessageHandler(Constants.PLAYER_DIED_MESSAGE,      playerDiedHandler);
  registerMessageHandler(Constants.PONG_MESSAGE,             pongHandler);
  registerMessageHandler(Constants.SCORE_BOARD_MESSAGE,      scoreBoardHandler);
};
