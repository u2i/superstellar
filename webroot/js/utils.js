export const translateToViewport = (x, y, viewport) => {
  const newX = x - viewport.vx + viewport.width / 2;
  const newY = -y + viewport.vy + viewport.height / 2;
  return {x: newX, y: newY}
};

export const getCurrentWindowSize = () => {
  const width = window.innerWidth
  || document.documentElement.clientWidth
  || document.body.clientWidth;

  const height = window.innerHeight
  || document.documentElement.clientHeight
  || document.body.clientHeight;

  return {width: width, height: height};
};
