import math from 'mathjs'

export default class MoveFilter {
  constructor(position, velocity) {
    this.dt = 10;

    // state model
    this.X = math.matrix([[position.x], [position.y], [velocity.x], [velocity.y]]);

    // covariance matrix
    this.P = math.matrix[0.1, 0.1, 0.1, 0.1];

    // state transition matrix
    this.F = math.matrix([[1, 0, this.dt, 0], [0, 1, 0, this.dt], [0, 0, 1, 0], [0, 0, 0, 1]]);
  }

  predict() {
    this.X = math.multiply(this.F, this.X);
    this.P = math.multiply(this.F, this.P, math.transpose(this.F));
  }

  position() {
    var pos = math.subset(this.X, math.index(0, [0, 1]));
    return {x: pos[0], y: pos[1]};
  }
}