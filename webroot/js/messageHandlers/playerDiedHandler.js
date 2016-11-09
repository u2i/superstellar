import GameOverDialog from '../dialogs/gameOverDialog';
import { globalState } from '../globals';

const playerDiedHandler = (message) => {
  const killedPlayer = message.id;
  const killedBy = message.killedBy;
  const myId = globalState.clientId;

  let spaceship = globalState.spaceshipMap.get(killedPlayer);
  spaceship.remove();
  globalState.spaceshipMap.delete(killedPlayer);

  if (killedPlayer === globalState.killedBy) {
    globalState.killedBy = null;
  }

  if (killedPlayer === myId) {
    globalState.killedBy = killedBy;

    const killedByName = globalState.clientIdToName.get(killedBy);
    const gameOverDialog = new GameOverDialog(killedByName);
    gameOverDialog.show();
  }

};

export default playerDiedHandler;
