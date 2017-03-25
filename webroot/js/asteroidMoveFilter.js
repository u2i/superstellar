import AsteroidSimulationFrame from './asteroidSimulationFrame.js';

export default class AsteroidMoveFilter {
  constructor(frameId) {
    this.frameId = frameId;
    this.simulationFrame = null;
  }

  update(updateFrameId, data) {
    this.simulationFrame = new AsteroidSimulationFrame(updateFrameId, data);
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
}
