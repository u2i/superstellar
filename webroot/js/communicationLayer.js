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

export const sendMessage = (protobufMsg) => {
  const message = new UserMessage();
  message.set(getMessageName(protobufMsg), protobufMsg);

  const buffer = message.encode();

  // TODO: we can probably handle this a bit better
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(buffer.toArrayBuffer());
  }
}

const getMessageName = (protobufMsg) => {
  const splitMessageName = protobufMsg.toString().split(".");

  const messageType = splitMessageName[splitMessageName.length - 1];

  return messageType.charAt(0).toLowerCase() + messageType.slice(1);
};

export const sendPing = () => {
  var pingId = nextPingId++;

  pingDates.set(pingId, new Date());

  const ping = new Ping();
  ping.set('Id', pingId)

  sendMessage(ping);
}

const builder = ProtoBuf.loadJsonFile(Constants.PROTOBUF_DEFINITION);

export const JoinGame    = builder.build(Constants.JOIN_GAME_DEFINITION);
export const Message     = builder.build(Constants.MESSAGE_DEFINITION);
export const Space       = builder.build(Constants.SPACE_DEFINITION);
export const UserMessage = builder.build(Constants.USER_MESSAGE_DEFINITION);
export const PlayerLeft  = builder.build(Constants.PLAYER_LEFT_DEFINITION);
export const UserEvent   = builder.build(Constants.USER_EVENT_DEFINITION);
export const UserAction  = builder.build(Constants.USER_ACTION_DEFINITION);
export const TargetAngle = builder.build(Constants.TARGET_ANGLE_DEFINITION);
export const Ping        = builder.build(Constants.PING_DEFINITION);
export const Pong        = builder.build(Constants.PONG_DEFINITION);


const messageHandlers = new Map();
var nextPingId = 0;
export const pingDates = new Map();

export const registerMessageHandler = (messageType, handler) => {
  let currentHandlers = messageHandlers.get(messageType) || [];

  currentHandlers.push(handler);

  messageHandlers.set(messageType, currentHandlers);
};

window.setInterval(function() {
  sendPing();
}, 1000)
