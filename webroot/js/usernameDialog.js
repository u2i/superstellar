import { sendMessage, JoinGame } from './communicationLayer';
import Cookie from 'js-cookie';
import * as Utils from "./utils";

const WIDTH = 300;
const HEIGHT = 150;

export default class UsernameDialog {
  constructor () {
    this.domNode = document.createElement("div");

    this.domNode.className = 'game-dialog';

    let previousNickname = Cookie.get('nickname') || '';

    this._updatePosition();
    this.domNode.innerHTML = `
    <div class="dialog-content">
      <p class="dialog-message">Welcome Captain... errhm... what was your name again?</p>
      <form id="submit-username-form">
	<input autofocus class="underline-input" id="insert-name-input" type="text" minlength="3" maxlength="25" value="${previousNickname}" required/>
	<input class="action-button" type="submit" value="Blast'em Off!" />
      </form>
    </div>
    `;
  }

  show () {
    document.body.appendChild(this.domNode);
    this.resizeListenerID = window.addEventListener("resize", () => { this._updatePosition() });
    this.submitListenerID  = document.
      getElementById("submit-username-form").
      addEventListener("submit", (ev) => {
        ev.preventDefault();
        this._sendJoinGame();
      });
  }

  showError (errorMsg) {
    const dialog = document.getElementsByClassName("dialog-message")[0];
    dialog.innerText = errorMsg;

    this.domNode.classList.add("error");
  }

  hide () {
    window.removeEventListener("resize", this.resizeListenerID);
    window.removeEventListener("submit", this.submitListenerID);
    document.body.removeChild(this.domNode);
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
    const nickname = document.getElementById("insert-name-input").value;
    Cookie.set('nickname', nickname, {expires: 7});
    sendMessage(new JoinGame(nickname));
  }
}
