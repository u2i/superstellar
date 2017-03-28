import GameOverDialog from '../dialogs/gameOverDialog';
import { globalState } from '../globals';

const playerDiedHandler = (message) => {
  const destroyedObjectId = message.id;
  const destroyedById = message.killedBy;
  const myId = globalState.clientId;

  if (globalState.spaceshipMap.has(destroyedObjectId)) {
    let spaceship = globalState.spaceshipMap.get(destroyedObjectId);
    spaceship.remove();
    globalState.spaceshipMap.delete(destroyedObjectId);

    if (destroyedObjectId === globalState.killedBy) {
      globalState.killedBy = null;
    }

    if (destroyedObjectId === myId) {
      globalState.killedBy = destroyedById;

      const killedByName = globalState.clientIdToName.get(destroyedById);
      const gameOverDialog = new GameOverDialog(killedByName);
      gameOverDialog.show();
    }
  } else if (globalState.asteroidsMap.has(destroyedObjectId)) {
    let asteroid = globalState.asteroidsMap.get(destroyedObjectId);
    asteroid.remove();
    globalState.asteroidsMap.delete(destroyedObjectId)
  }
};

export default playerDiedHandler;
