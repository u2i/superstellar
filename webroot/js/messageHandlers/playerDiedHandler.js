import GameOverDialog from '../dialogs/gameOverDialog';
import { globalState } from '../globals';

const playerDiedHandler = (message) => {
  const playerId = message.id;
  const killedBy = message.killedBy
  const killedByName = globalState.clientIdToName.get(killedBy)

  let spaceship = globalState.spaceshipMap.get(playerId);
  spaceship.remove();
  globalState.spaceshipMap.delete(playerId);

  if (globalState.killedBy === killedBy) {
    globalState.killedBy = null
  }

  if (playerId === globalState.clientId) {
    globalState.killedBy = killedBy

    const gameOverDialog = new GameOverDialog(killedByName);
    gameOverDialog.show();
  }

};

export default playerDiedHandler;
