import * as Utils from './utils';

const shaderContent = require('raw!../shaders/annulus_fog.frag');

export default class AnnulusFilter extends PIXI.Filter {
  constructor () {
    super(null, shaderContent);
    this._setWindowSize();

    this.uniforms.worldCoordinates = new Float32Array([0.0, 0.0]);
    this.uniforms.worldBoundarySize = new Float32Array([1000.0, 1400.0]);
    this.uniforms.magicMatrix = new PIXI.Matrix;

    window.addEventListener("resize", () => {
      this._setWindowSize();
    });
  }

  apply (filterManager, input, output) {
    filterManager.calculateNormalizedScreenSpaceMatrix(this.uniforms.magicMatrix);
    filterManager.applyFilter(this, input, output);
  }

  get worldCoordinates () {
    return this.uniforms.worldCoordinates;
  }
  
  set worldCoordinates (value) {
    this.uniforms.worldCoordinates = value;
  }

  get worldBoundarySize () {
    return this.uniforms.worldBoundarySize;
  }

  set worldBoundarySize (value) {
    this.uniforms.worldBoundarySize = value;
  }

  get windowSize () {
    return this.uniforms.windowSize;
  }

  _setWindowSize() {
    const { width, height } = Utils.getCurrentWindowSize();

    this.uniforms.windowSize = new Float32Array([width, height]);
  }
}
