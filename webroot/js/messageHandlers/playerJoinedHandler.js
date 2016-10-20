import { globalState } from '../globals';

const playerJoinedHandler = (message) => {
  const { id, username } = message;

  globalState.clientIdToName.set(id, username);
};

export default playerJoinedHandler;
