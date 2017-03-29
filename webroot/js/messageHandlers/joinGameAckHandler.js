import { initializeControls } from '../controls';
import { globalState, leaderboardDialog } from '../globals';

const joinGameAckHandler = (message) => {
  const { success, error } = message;

  if (success) {
    globalState.dialog.hide();
    leaderboardDialog.show();
    initializeControls();
  } else {
    globalState.dialog.showError(error);
  }
};

export default joinGameAckHandler;
