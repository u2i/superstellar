import { globalState } from '../globals';

const playerLeftHandler = (message) => {
  const playerId = message.id;

  let spaceship = globalState.spaceshipMap.get(playerId);
  if(typeof spaceship !== "undefined") {
    spaceship.remove();
    globalState.spaceshipMap.delete(playerId);
  }
};

export default playerLeftHandler;
