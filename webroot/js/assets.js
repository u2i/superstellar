const Assets = {
  getTexture: (texturePath) => {
    return PIXI.loader.resources[texturePath].texture;
  },

  getTextureFromFrame: (frameName) => {
    return PIXI.Texture.fromFrame(frameName);
  }
};

export default Assets;
