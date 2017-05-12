import { sendMessage, UserAction, TargetAngle } from './communicationLayer';
import { globalState, renderer } from './globals';

var mouseDown = false;
var lastTargetAngle = null;

const KEY_SPACE = 32;
const KEY_UP    = 38;
const KEY_LEFT  = 37;
const KEY_RIGHT = 39;
const KEY_W     = 87;
const KEY_A     = 65;
const KEY_D     = 68;
const KEY_SHIFT = 16;

const keysDown = new Map();

keysDown.set(KEY_SPACE, false);
keysDown.set(KEY_UP,    false);
keysDown.set(KEY_LEFT,  false);
keysDown.set(KEY_RIGHT, false);
keysDown.set(KEY_W, false);
keysDown.set(KEY_A, false);
keysDown.set(KEY_D, false);
keysDown.set(KEY_SHIFT, false);

const updateKeysState = (keyCode, isPressed) => {
  const lastState = keysDown.get(keyCode);
  if (lastState === undefined) {
    return;
  }

  if (lastState !== isPressed) {
    keysDown.set(keyCode, isPressed);
    sendInput(keyCode, isPressed);
  }
}

const updateMouseState = (isDown) => {
  if (mouseDown !== isDown) {
    mouseDown = isDown;

    let userAction = new UserAction(isDown ? "TURRET_FIRE_START" : "FIRE_STOP");
    sendMessage(userAction);
  }
}

const updateMousePosition = (event) => {
  globalState.crosshair.update(event.x, event.y);
}

const sendInput = (keyCode, isPressed) => {
  let userInput = "CENTER"

  switch(keyCode) {
  case KEY_UP:
  case KEY_W:
    userInput = isPressed ? "THRUST_ON" : "THRUST_OFF";
    break;
  case KEY_LEFT:
  case KEY_A:
    userInput = isPressed ? "LEFT" : "CENTER";
    break;
  case KEY_RIGHT:
  case KEY_D:
    userInput = isPressed ? "RIGHT" : "CENTER";
    break;
  case KEY_SPACE:
    userInput = isPressed ? "STRAIGHT_FIRE_START" : "FIRE_STOP";
    break;
  case KEY_SHIFT:
    userInput = isPressed ? "BOOST_ON" : "BOOST_OFF";
    break;
  }

  let userAction = new UserAction(userInput);

  sendMessage(userAction);
}

export const initializeControls = () => {
  addEventListener("keydown", function (e) {
    updateKeysState(e.keyCode, true);
  }, false);

  addEventListener("keyup", function (e) {
    updateKeysState(e.keyCode, false);
  }, false);

  addEventListener("mousedown", function() {
    updateMouseState(true);
  }, false);

  addEventListener("mouseup", function() {
    updateMouseState(false);
  }, false);

  addEventListener("mousemove", function(e) {
    updateMousePosition(e);
  }, false);
};

window.setInterval(function() {
  var mousePosition = renderer.plugins.interaction.mouse.global;

  let x = mousePosition.x - renderer.width / 2;
  let y = mousePosition.y - renderer.height / 2;
  let targetAngle = Math.atan2(y, x);

  if (lastTargetAngle != targetAngle) {
    lastTargetAngle = targetAngle;

    let targetAngleMsg = new TargetAngle(targetAngle);
    sendMessage(targetAngleMsg);
  }
}, 100)
