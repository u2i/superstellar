import { pingDates } from '../communicationLayer.js';
import { globalState } from '../globals';

const pongHandler = (message) => {
  var pongDate = new Date();

  if (pingDates.has(message.Id)) {
    let pingDate = pingDates.get(message.Id);
    let delta = pongDate - pingDate;
    globalState.ping = delta

    pingDates.delete(message.Id);
  }
};

export default pongHandler;
