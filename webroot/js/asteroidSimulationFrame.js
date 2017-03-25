import Victor from 'victor';

Victor.prototype.scalarMultiply = function(scalar) {
  return this.multiply(new Victor(scalar, scalar));
}

export default class AsteroidSimulationFrame {
  constructor(frameId, data) {
    this.frameId = frameId;

    this.update(data);
  }

  update(data) {
    this.id = data.id;
    this.position = Victor.fromObject(data.position);
    this.velocity = Victor.fromObject(data.velocity);
    this.facing = data.facing;
    this.angularVelocity = data.angularVelocity;
    this.angularVelocityDelta = 0.0;
  }

  predict() {
    this.position.add(this.velocity);
    this.updateAngularVelocity();

    this.frameId++;
  }

  predictTo(targetFrameId) {
    while (this.frameId < targetFrameId) {
      this.predict()
    }
  }

  updateAngularVelocity() {
    this.angularVelocity += this.angularVelocityDelta
    this.angularVelocityDelta = 0.0;
    this.facing -= this.angularVelocity;
  }
}
