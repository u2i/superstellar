import { globalState } from '../globals';

const helloHandler = (message) => {
  globalState.clientId = message.myId;
  message.idToUsername.forEach((username, id) => {
    globalState.clientIdToName.set(id, username);
  });

};

export default helloHandler;
