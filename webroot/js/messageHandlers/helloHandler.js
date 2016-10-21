import { globalState } from '../globals';

const helloHandler = (message) => {
  const { myId, idToUsername } = message;
  globalState.clientId = myId;
  idToUsername.forEach((username, id) => {
    globalState.clientIdToName.set(id, username);
  });
};

export default helloHandler;
