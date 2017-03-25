import { constants } from './globals';
import Victor from 'victor';

const DIR_CENTER = null;
const DIR_RIGHT = 1;
const DIR_LEFT = 2;

Victor.prototype.scalarMultiply = function(scalar) {
  return this.multiply(new Victor(scalar, scalar));
}

export default class SpaceshipSimulationFrame {
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
    this.inputDirection = data.inputDirection;
    this.inputThrust = data.inputThrust;
    this.inputBoost = data.inputBoost;
    this.hp = data.hp;
    this.maxHp = data.maxHp;
    this.energy = data.energy;
    this.maxEnergy = data.maxEnergy;
    this.autoRepairDelay = data.autoRepairDelay;
  }

  predict() {
    this.applyInputThrust();
    this.applyAnnulus();
    this.limitMaxSpeed();

    this.applyTurn();

    this.position.add(this.velocity);
    this.updateAngularVelocity();

    this.applyAutoRepair();
    this.applyAutoEnergyRecharge();

    this.frameId++;
  }

  predictTo(targetFrameId) {
    while (this.frameId < targetFrameId) {
      this.predict()
    }
  }

  applyInputThrust() {
    if (this.inputThrust || this.inputBoost) {
      let facingVector = new Victor(Math.cos(this.facing), -Math.sin(this.facing));
      let deltaVelocity = facingVector.scalarMultiply(constants.spaceshipAcceleration);
      this.velocity.add(deltaVelocity);
    } else {
      if (this.velocity.length() > 0) {
        this.velocity.scalarMultiply(1.0 - constants.frictionCoefficient);
      }

      if (this.velocity.length() < 1.0) {
        this.velocity.scalarMultiply(0.0);
      }
    }
  }

  applyAnnulus() {
    if (this.position.length() > constants.worldRadius) {
      let outreachLength = this.position.length() - constants.worldRadius;
      let gravityAcceleration = -(outreachLength / constants.boundaryAnnulusWidth) * constants.spaceshipAcceleration;
      let deltaVelocity = this.position.clone().normalize().scalarMultiply(gravityAcceleration);
      this.velocity.add(deltaVelocity);
    }
  }

  limitMaxSpeed() {
    let boostActive = false;

    if (this.inputBoost) {
      if (this.energy >= constants.boostPerFrameEnergyCost) {
        this.energy -= constants.boostPerFrameEnergyCost;
        boostActive = true;
      }
    }

    this.inputBoost = boostActive;

    let maxVelocity = constants.spaceshipMaxSpeed;
    if (boostActive) {
      maxVelocity *= constants.spaceshipBoostFactor;
    }

    if (this.velocity.length() > maxVelocity) {
      this.velocity.normalize().scalarMultiply(maxVelocity);
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
    this.angularVelocity *= (1.0 - constants.spaceshipAngularFriction);
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
    let nonlinearPart = constants.spaceshipNonlinearAngularAcceleration * Math.abs(this.angularVelocity);
    let linearPart = constants.spaceshipLinearAngularAcceleration;
    return nonlinearPart + linearPart;
  }

  limitAngularVelocityDelta() {
    let potentialAngularVelocity = this.angularVelocity + this.angularVelocityDelta;
    let diff = Math.abs(potentialAngularVelocity) - constants.spaceshipMaxAngularVelocity;

    if (diff > 0) {
      this.angularVelocityDelta -= Math.abs(diff) * Math.sign(this.angularVelocity);
    }
  }

  applyAutoRepair() {
    if (this.autoRepairDelay === 0) {
      if (this.hp < this.maxHp) {
        this.hp = Math.min(this.hp + constants.autoRepairAmount, this.maxHp)
        this.autoRepairDelay = constants.autoRepairInterval;
      }
    } else {
      this.autoRepairDelay--;
    }
  }

  applyAutoEnergyRecharge() {
    this.energy = Math.min(this.energy + constants.autoEnergyRechargeAmount, this.maxEnergy);
  }
}
