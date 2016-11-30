const shaderContent = require('raw!../shaders/circle_bar.frag');

export default class CircleBarFilter extends PIXI.Filter {
  constructor(basicColor) {
    super(null, shaderContent);

    this.uniforms.magicMatrix = new PIXI.Matrix;
    this.uniforms.basicColor = new Float32Array(basicColor);
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
