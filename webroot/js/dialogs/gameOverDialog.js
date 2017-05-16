import { sendMessage, JoinGame } from '../communicationLayer';
import { globalState } from '../globals';
import * as Utils from "../utils";

const WIDTH = 300;
const HEIGHT = 130;

export default class GameOverDialog {
  constructor (killedBy) {
    this.domNode = document.createElement("div");

    this.domNode.className = 'game-dialog';

    this._updatePosition();
    this.domNode.innerHTML = `
    <div class="dialog-content">
      <p class="dialog-message">You died, captain!</p>
      <p class="dialog-message">Your score was: ` + globalState.score + `!</p>
      <button id="submit" class="action-button" type="button" autofocus>Take revenge!</button>
    </div>
    `;
    setTimeout(function () {
      document.getElementById('submit').focus();
    }, 50);
  }

  show () {
    globalState.dialog = this

    document.body.appendChild(this.domNode);
    this.resizeListenerID = window.addEventListener("resize", () => { this._updatePosition() });
    this.submitListenerID  = document.
          getElementById("submit").
          addEventListener("click", (ev) => {
            ev.preventDefault();
            this._sendJoinGame();
          });
  }

  showError (errorMsg) {
    const dialog = document.getElementsByClassName("dialog-message")[0];
    dialog.textContent = errorMsg;
    this.domNode.classList.add("error");
    globalState.dialog = null

  }

  hide () {
    window.removeEventListener("resize", this.resizeListenerID);
    window.removeEventListener("submit", this.submitListenerID);
    document.body.removeChild(this.domNode);

    globalState.dialog = null
  }

  _updatePosition () {
    const { width, height } = Utils.getCurrentWindowSize();

    const x = (width - WIDTH) / 2 + 30;
    const y = (height - HEIGHT) / 2 + 30;

    this.domNode.style.top  = `${y}px`;
    this.domNode.style.left = `${x}px`;
    this.domNode.style.width = `${WIDTH}px`;
    this.domNode.style.height = `${HEIGHT}px`;
  }

  _sendJoinGame () {
    sendMessage(new JoinGame(globalState.nickname));
  }
}
