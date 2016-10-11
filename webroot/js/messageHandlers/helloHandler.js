import { globalState } from '../globals';

const helloHandler = (message) => {
  globalState.clientId = message.myId;
};

export default helloHandler;
