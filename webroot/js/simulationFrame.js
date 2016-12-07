import Victor from 'victor';

const SpaceshipAcceleration = 50.0
const WorldRadius = 100000
const BoundaryAnnulusWidth = 20000
const FrictionCoefficient = 0.005
const SpaceshipNonlinearAngularAcceleration = 2
const SpaceshipLinearAngularAcceleration = 0.0001
const SpaceshipMaxAngularVelocity = 0.12
const SpaceshipAngularFriction = 0.2
const SpaceshipMaxSpeed = 600

const DIR_CENTER = null;
const DIR_RIGHT = 1;
const DIR_LEFT = 2;

Victor.prototype.scalarMultiply = function(scalar) {
  return this.multiply(new Victor(scalar, scalar));
}

export default class SimulationFrame {
  constructor(frameId, data) {
    this.frameId = frameId;

    this.update(data);
  }

  update({id, position, velocity, facing, angularVelocity, inputDirection, inputThrust}) {
    this.id = id;
    this.position = Victor.fromObject(position);
    this.velocity = Victor.fromObject(velocity);
    this.facing = facing;
    this.angularVelocity = angularVelocity;
    this.angularVelocityDelta = 0.0;
    this.inputDirection = inputDirection;
    this.inputThrust = inputThrust;
  }

  predict() {
    this.applyInputThrust();
    this.applyAnnulus();
    this.limitMaxSpeed();

    this.position.add(this.velocity);

    this.applyTurn();
    this.updateAngularVelocity();

    this.frameId++;
  }

  predictTo(targetFrameId) {
    while (this.frameId < targetFrameId) {
      this.predict()
    }
  }

  applyInputThrust() {
    if (this.inputThrust) {
      let facingVector = new Victor(Math.cos(this.facing), -Math.sin(this.facing));
      let deltaVelocity = facingVector.scalarMultiply(SpaceshipAcceleration);
      this.velocity.add(deltaVelocity);
    } else {
      if (this.velocity.length() > 0) {
        this.velocity.scalarMultiply(1.0 - FrictionCoefficient);
      }

      if (this.velocity.length() < 1.0) {
        this.velocity.scalarMultiply(0.0);
      }
    }
  }

  applyAnnulus() {
    if (this.position.length() > WorldRadius) {
      let outreachLength = this.position.length() - WorldRadius;
      let gravityAcceleration = -(outreachLength / BoundaryAnnulusWidth) * SpaceshipAcceleration;
      let deltaVelocity = this.position.clone().normalize().scalarMultiply(gravityAcceleration);
      this.velocity.add(deltaVelocity);
    }
  }

  limitMaxSpeed() {
    if (this.velocity.length() > SpaceshipMaxSpeed) {
      this.velocity.normalize().scalarMultiply(SpaceshipMaxSpeed);
    }
  }

  applyTurn() {
    switch(this.inputDirection) {
    case DIR_CENTER:
      this.applyAngularFriction();
      break;
    case DIR_RIGHT:
      this.turnRight();
      break;
    case DIR_LEFT:
      this.turnLeft();
      break;
    }
  }

  updateAngularVelocity() {
    this.angularVelocity += this.angularVelocityDelta
    this.angularVelocityDelta = 0.0;
    this.facing -= this.angularVelocity;
  }

  applyAngularFriction() {
    this.angularVelocity *= (1.0 - SpaceshipAngularFriction);
  }

  turnLeft() {
    this.angularVelocityDelta = this.angularVelocityDeltaValue();
    this.limitAngularVelocityDelta();
  }

  turnRight() {
    this.angularVelocityDelta = -this.angularVelocityDeltaValue();
    this.limitAngularVelocityDelta();
  }

  angularVelocityDeltaValue() {
    let nonlinearPart = SpaceshipNonlinearAngularAcceleration * Math.abs(this.angularVelocity);
    let linearPart = SpaceshipLinearAngularAcceleration;
    return nonlinearPart + linearPart;
  }

  limitAngularVelocityDelta() {
    let potentialAngularVelocity = this.angularVelocity + this.angularVelocityDelta;
    let diff = Math.abs(potentialAngularVelocity) - SpaceshipMaxAngularVelocity;

    if (diff > 0) {
      this.angularVelocityDelta -= Math.abs(diff) * Math.sign(this.angularVelocity);
    }
  }


}
