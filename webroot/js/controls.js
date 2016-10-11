import { sendMessage, UserMessage } from './communicationLayer';

const KEY_SPACE = 32;
const KEY_UP    = 38;
const KEY_LEFT  = 37;
const KEY_RIGHT = 39;

const keysDown = new Map();

keysDown.set(KEY_SPACE, false);
keysDown.set(KEY_UP,    false);
keysDown.set(KEY_LEFT,  false);
keysDown.set(KEY_RIGHT, false);

const updateKeysState = (keyCode, isPressed) => {
  const lastState = keysDown.get(keyCode);
  if (lastState == undefined) {
    return;
  }

  if (lastState != isPressed) {
    keysDown.set(keyCode, isPressed);
    sendInput(keyCode, isPressed);
  }
}

const sendInput = (keyCode, isPressed) => {
  let userInput = "CENTER"

  switch(keyCode) {
    case KEY_UP:
      userInput = isPressed ? "THRUST_ON" : "THRUST_OFF";
      break;
    case KEY_LEFT:
      userInput = isPressed ? "LEFT" : "CENTER"
      break;
    case KEY_RIGHT:
      userInput = isPressed ? "RIGHT" : "CENTER"
      break;
    case KEY_SPACE:
      userInput = isPressed ? "FIRE_START" : "FIRE_STOP"
      break;
  }

  let userMessage = new UserMessage(userInput);

  sendMessage(userMessage);
}

export const initializeControls = () => {
  addEventListener("keydown", function (e) {
    updateKeysState(e.keyCode, true);
  }, false);

  addEventListener("keyup", function (e) {
    updateKeysState(e.keyCode, false);
  }, false);
};
