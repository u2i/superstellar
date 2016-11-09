import { initializeControls } from '../controls';
import { globalState } from '../globals';

const joinGameAckHandler = (message) => {
  const { success, error } = message;

  if (success) {
    globalState.dialog.hide();
    initializeControls();
  } else {
    globalState.dialog.showError(error);
  }
};

export default joinGameAckHandler;
