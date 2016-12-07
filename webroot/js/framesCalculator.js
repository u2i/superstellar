const PhysicsFrameDuration = 20;

export default class FramesCalculator {
  constructor() {
    this.simulationStartTime = new Date();
  }

  receivedFrameId(receivedFrameId) {
    let computedElapsedSinceSimulationStart = (receivedFrameId - 1) * PhysicsFrameDuration;
    let computedSimulationStartTime = new Date() - computedElapsedSinceSimulationStart;

    if (computedSimulationStartTime < this.simulationStartTime) {
      this.simulationStartTime = computedSimulationStartTime;
      console.log("simulationStartTime: " + this.simulationStartTime);
    }
  }

  currentFrameId() {
    let elapsedSinceSimulationStart = new Date() - this.simulationStartTime
    return Math.floor(elapsedSinceSimulationStart / PhysicsFrameDuration);
  }
}