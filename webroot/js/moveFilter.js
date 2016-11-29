export default class MoveFilter {
  constructor(x, y, vx, vy) {
    this.dt = 10

    // state model
    this.X = math.matrix([[x], [y], [vx], [vy]])

    // covariance matrix
    this.P = math.matrix[0.1, 0.1, 0.1, 0.1]

    // state transition matrix
    this.F = math.matrix([[1, 0, this.dt, 0], [0, 1, 0, this.dt], [0, 0, 1, 0], [0, 0, 0, 1]]);
  }

  predict {
    this.X = math.multiply(this.F, this.X)
    this.P = math.multiply(this.F, this.P, math.transpose(this.F))
  }
}