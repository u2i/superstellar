import { renderer } from './globals';

const shaderContent = require('raw!../shaders/annulus_fog.frag');

export default class AnnulusFilter {
  constructor () {
    PIXI.Filter.call(this, null, shaderContent);
    this.uniforms.worldCoordinates = new Float32Array([0.0, 0.0]);
    this.uniforms.worldSize = new Float32Array([1000, 1400]);
    this.uniforms.dimensions = new Float32Array([renderer.width, renderer.height, 0, 0]);
    this.uniforms.magicMatrix = new PIXI.Matrix;
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

  get worldSize () {
    return this.uniforms.worldSize;
  }

  set worldSize (value) {
    this.uniforms.worldSize = value;
  }

  get dimensions () {
    return this.uniforms.dimensions;
  }
}
