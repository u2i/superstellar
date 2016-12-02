import math  from 'mathjs'

export default class MoveFilter {
  constructor(position, velocity) {
    this.dt = 1;

    // state model
    this.X = math.matrix([[position.x], [position.y], [velocity.x], [velocity.y]]);

    // update projection matric
    this.H = math.eye(4);

    this.Ht = math.transpose(this.H);

    // covariance matrix
    this.P = math.multiply(0.001, math.eye(4));

    this.R = math.multiply(0.001, math.eye(4));

    // state transition matrix
    this.F = math.matrix([[1, 0, this.dt, 0], [0, 1, 0, this.dt], [0, 0, 1, 0], [0, 0, 0, 1]]);
  }

  predict() {
    this.X = math.multiply(this.F, this.X);
    this.P = math.multiply(math.multiply(this.F, this.P), math.transpose(this.F));
  }

  update(position, velocity) {
    console.log(position, velocity);
    let Z = math.matrix([[position.x], [position.y], [velocity.x], [velocity.y]]);

    let ph = math.multiply(this.P, this.Ht);
    let hph = math.multiply(math.multiply(this.H, this.P), this.Ht);

    let K = math.multiply(ph, math.inv(math.add(hph, this.R)));

    let zhx = math.subtract(Z, this.X);
    console.log(math.subset(this.X, math.index(0, 0)), math.subset(this.X, math.index(1, 0)))
    console.log(math.subset(zhx, math.index(0, 0)), math.subset(zhx, math.index(1, 0)))


    this.X = math.add(this.X, math.multiply(K, zhx));
    this.P = math.subtract(this.P, math.multiply(math.multiply(K, this.H), this.P));
    console.log('P');

  }

  position() {
    var x = math.subset(this.X, math.index(0, 0));
    var y = math.subset(this.X, math.index(1, 0));
    var position = {x: x, y: y};

    return position;
  }
}