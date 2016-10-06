import ProtoBuf from 'protobufjs';
import * as Constants from './constants';

let ws; 

const webSocketMessageReceived = (e) => {
  var fileReader = new FileReader();

  fileReader.onload = function() {
    handleProtoBufMessage(this.result);
  };

  fileReader.readAsArrayBuffer(e.data);
};

export const initializeConnection = (host, port, path) => {
  ws = new WebSocket("ws://" + host + ':' + port + path);
  ws.onmessage = webSocketMessageReceived;
};

export const handleProtoBufMessage = (protoBufMsg) => {
  const message = Message.decode(protoBufMsg);

  const handlers = messageHandlers.get(message.content);

  if (!handlers) { 
    console.log(`Handlers for ${message.content} are not registered`) 
  } else {
    const messageContent = message[message.content];

    handlers.forEach((handler) => {
      handler(messageContent);
    });
  }
};

export const sendUserInput = (thrust, direction) => {
};

export const sendMessage = (protobufMsg) => {
  let buffer = protobufMsg.encode();

  // TODO: we can probably handle this a bit better
  if (ws.readyState == WebSocket.OPEN) {
    ws.send(buffer.toArrayBuffer());
  }
}

const builder    = ProtoBuf.loadJsonFile(Constants.PROTOBUF_DEFINITION);
export const Message    = builder.build(Constants.MESSAGE_DEFINITION);
export const Space      = builder.build(Constants.SPACE_DEFINITION);
export const UserInput  = builder.build(Constants.USER_INPUT_DEFINITION);
export const PlayerLeft = builder.build(Constants.PLAYER_LEFT_DEFINITION);

const messageHandlers = new Map();

export const registerMessageHandler = (messageType, handler) => {
  let currentHandlers = messageHandlers.get(messageType) || [];

  currentHandlers.push(handler);

  messageHandlers.set(messageType, currentHandlers);
};

