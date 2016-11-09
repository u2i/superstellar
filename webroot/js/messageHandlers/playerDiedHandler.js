import GameOverDialog from '../dialogs/gameOverDialog';
import { globalState } from '../globals';

const playerDiedHandler = (message) => {
  const playerId = message.id;
  const killedBy = message.killedBy
  const killedByName = globalState.clientIdToName.get(killedBy)

  let spaceship = globalState.spaceshipMap.get(playerId);

  spaceship.remove();

  globalState.spaceshipMap.delete(playerId);

  if (playerId === globalState.clientId) {
    const gameOverDialog = new GameOverDialog(killedByName);
    gameOverDialog.show();
  }

};

export default playerDiedHandler;
