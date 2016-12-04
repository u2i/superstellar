import Victor from 'victor';

const SpaceshipAcceleration = 50.0
const WorldRadius = 100000
const BoundaryAnnulusWidth = 20000
const FrictionCoefficient = 0.005
const SpaceshipNonlinearAngularAcceleration = 2
const SpaceshipLinearAngularAcceleration = 0.0001
const SpaceshipMaxAngularSpeed = 0.12
const SpaceshipAngularFriction = 0.2
const SpaceshipMaxSpeed = 600

const DIR_CENTER = null;
const DIR_RIGHT = 1;
const DIR_LEFT = 2;

Victor.prototype.scalarMultiply = function(scalar) {
  return this.multiply(new Victor(scalar, scalar));
}

export default class SimulationFrame {
  constructor(physicalFrameId, data) {
    this.physicalFrameId = physicalFrameId;

    this.update(data);
  }

  update({id, position, velocity, facing, angularSpeed, inputDirection, inputThrust}) {
    this.id = id;
    this.position = Victor.fromObject(position);
    this.velocity = Victor.fromObject(velocity);
    this.facing = facing;
    this.angularSpeed = angularSpeed;
    this.inputDirection = inputDirection;
    this.inputThrust = inputThrust;
  }

  predict() {
    this.applyInputThrust();
    this.applyAnnulus();
    this.limitMaxSpeed();

    this.position.add(this.velocity);

    this.applyTurn();

    this.facing += this.angularSpeed;
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

  applyAngularFriction() {
    this.angularSpeed *= (1.0 - SpaceshipAngularFriction);
  }

  turnLeft() {
    this.angularSpeed -= this.angularSpeedDelta();
    this.limitAngularSpeed();
  }

  turnRight() {
    this.angularSpeed += this.angularSpeedDelta();
    this.limitAngularSpeed();
  }

  angularSpeedDelta() {
    let nonlinearPart = SpaceshipNonlinearAngularAcceleration * Math.abs(this.angularSpeed);
    let linearPart = SpaceshipLinearAngularAcceleration;
    return nonlinearPart + linearPart;
  }

  limitAngularSpeed() {
    if (Math.abs(this.angularSpeed) > SpaceshipMaxAngularSpeed) {
      this.angularSpeed = SpaceshipMaxAngularSpeed * Math.sign(this.angularSpeed);
    }
  }


}
