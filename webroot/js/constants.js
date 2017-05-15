// ********** TEXTURES **********
export const BACKGROUND_TEXTURE = "images/background1.png";
export const SHIP_TEXTURE = "images/ship.png";
export const RADER_REFRESH_TEXTURE = "images/radar_refresh.png";
export const CROSSHAIR_TEXTURE = "images/crosshair_32.png";
export const ASTEROID_01_TEXTURE = "images/asteroid_01.png";

// ********** SPRITESHEETS **********
export const FLAME_SPRITESHEET = "spritesheets/flame_yellow.json";
export const FLAME_SPRITESHEET_FRAME_NAMES = [0, 1, 2, 3].map((i) => `thrust_yellow_${i}.png`);

export const PROJECTILE_SPRITESHEET = "spritesheets/projectile.json";
export const PROJECTILE_SPRITESHEET_FRAME_NAMES = [1, 2, 3].map((i) => `bullet${i}.png`);

export const BOOST_SPRITESHEET = "spritesheets/boost.json";
export const BOOST_SPRITESHEET_FRAME_NAMES = [1, 2].map((i) => `boost_0${i}.png`);

// ********** PROTOBUFS **********
export const PROTOBUF_DEFINITION      = "js/superstellar_proto.json";
export const SPACE_DEFINITION         = "superstellar.Space";
export const JOIN_GAME_DEFINITION     = "superstellar.JoinGame";
export const USER_ACTION_DEFINITION   = "superstellar.UserAction";
export const USER_MESSAGE_DEFINITION  = "superstellar.UserMessage";
export const TARGET_ANGLE_DEFINITION  = "superstellar.TargetAngle";
export const PLAYER_LEFT_DEFINITION   = "superstellar.PlayerLeft";
export const SHOT_DEFINITION          = "superstellar.Shot";
export const MESSAGE_DEFINITION       = "superstellar.Message";
export const PING_DEFINITION          = "superstellar.Ping";
export const PONG_DEFINITION          = "superstellar.Pong";

export const HELLO_MESSAGE             = "hello";
export const CONSTANTS_MESSAGE         = "constants";
export const JOIN_GAME_ACK_MESSAGE     = "joinGameAck";
export const SPACE_MESSAGE             = "space";
export const LEADERBOARD_MESSAGE       = "leaderboard";
export const PLAYER_LEFT_MESSAGE       = "playerLeft";
export const PLAYER_JOINED_MESSAGE     = "playerJoined";
export const PROJECTILE_FIRED_MESSAGE  = "projectileFired";
export const PROJECTILE_HIT_MESSAGE    = "projectileHit";
export const PLAYER_DIED_MESSAGE       = "playerDied";
export const PING_MESSAGE              = "ping";
export const PONG_MESSAGE              = "pong";
export const SCORE_BOARD_MESSAGE       = "scoreBoard";



