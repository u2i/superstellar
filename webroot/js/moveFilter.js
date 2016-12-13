import SimulationFrame from './simulationFrame.js';

export default class MoveFilter {
  constructor(frameId) {
    this.frameId = frameId;
    this.simulationFrame = null;
  }

  update(updateFrameId, data) {
    this.simulationFrame = new SimulationFrame(updateFrameId, data);
    this.simulationFrame.predictTo(this.frameId);
  }

  predictTo(frameId) {
    this.frameId = frameId;
    this.simulationFrame.predictTo(frameId);
  }

  position() {
    return this.simulationFrame.position;
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
