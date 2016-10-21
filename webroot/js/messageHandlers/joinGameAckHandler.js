import { initializeControls } from '../controls';
import { usernameDialog } from '../globals';

const joinGameAckHandler = (message) => {
  const { success, error } = message;

  if (success) {
    usernameDialog.hide();
    initializeControls();
  } else {
    usernameDialog.showError(error);
  }
};

export default joinGameAckHandler;
