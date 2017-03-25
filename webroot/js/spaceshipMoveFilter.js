import SpaceshipSimulationFrame from './spaceshipSimulationFrame.js';

export default class SpaceshipMoveFilter {
  constructor(frameId) {
    this.frameId = frameId;
    this.simulationFrame = null;
  }

  update(updateFrameId, data) {
    this.simulationFrame = new SpaceshipSimulationFrame(updateFrameId, data);
    this.simulationFrame.predictTo(this.frameId);
  }

  predictTo(frameId) {
    this.frameId = frameId;
    this.simulationFrame.predictTo(frameId);
  }

  position() {
    return this.simulationFrame.position;
  }

  velocity() {
    return this.simulationFrame.velocity;
  }

  facing() {
    return this.simulationFrame.facing;
  }

  inputThrust() {
    return this.simulationFrame.inputThrust;
  }

  inputBoost() {
    return this.simulationFrame.inputBoost;
  }

  hp() {
    return this.simulationFrame.hp;
  }

  maxHp() {
    return this.simulationFrame.maxHp;
  }

  energy() {
    return this.simulationFrame.energy;
  }

  maxEnergy() {
    return this.simulationFrame.maxEnergy;
  }
}
