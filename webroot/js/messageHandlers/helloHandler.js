import { globalState, usernameDialog } from '../globals';
import { initializeControls } from '../controls';

const helloHandler = (message) => {
  const { myId, idToUsername } = message;
  usernameDialog.hide();
  globalState.clientId = myId;
  idToUsername.forEach((username, id) => {
    globalState.clientIdToName.set(id, username);
  });
  initializeControls();
};

export default helloHandler;
