import * as Utils from "./utils";

const WIDTH = 300;
const HEIGHT = 150;

export default class GameOverDialog {
  constructor (killedBy) {
    this.domNode = document.createElement("div");

    this.domNode.className = 'game-dialog';

    this._updatePosition();
    this.domNode.innerHTML = `
    <div class="dialog-content">
      <p class="dialog-message">You died, captain!</p>
      <p>You were killed by: ` + killedBy + `</p>
    </div>
    `;
  }

  show () {
    document.body.appendChild(this.domNode);
    this.resizeListenerID = window.addEventListener("resize", () => { this._updatePosition() });
  }

  hide () {
    window.removeEventListener("resize", this.resizeListenerID);
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
}
