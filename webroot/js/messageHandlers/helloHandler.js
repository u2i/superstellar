import { globalState, usernameDialog } from '../globals';
import { initializeControls } from '../controls';

const helloHandler = (message) => {
  usernameDialog.hide();
  globalState.clientId = message.myId;
  initializeControls();
};

export default helloHandler;
