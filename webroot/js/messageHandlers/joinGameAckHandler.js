import { initializeControls } from '../controls';
import { globalState, leaderboardDialog, scoreBoardDialog } from '../globals';

const joinGameAckHandler = (message) => {
  const { success, error } = message;

  if (success) {
    globalState.dialog.hide();
    leaderboardDialog.show();
    scoreBoardDialog.show();
    initializeControls();
  } else {
    globalState.dialog.showError(error);
  }
};

export default joinGameAckHandler;
