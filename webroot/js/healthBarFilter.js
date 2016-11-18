const shaderContent = require('raw!../shaders/circle_healthbar.frag');

export default class HealthBarFilter extends PIXI.Filter {
  constructor() {
    super(null, shaderContent);

    this.uniforms.magicMatrix = new PIXI.Matrix;
  }

  apply(filterManager, input, output) {
    filterManager.calculateNormalizedScreenSpaceMatrix(this.uniforms.magicMatrix);
    filterManager.applyFilter(this, input, output);
  }

  get hps() {
    return this.uniforms.hps;
  }

  set hps(value) {
    this.uniforms.hps = new Float32Array(value);
  }
}
