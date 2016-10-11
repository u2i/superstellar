export const translateToViewport = (x, y, viewport) => {
	const newX = x - viewport.vx + viewport.width / 2;
	const newY = -y + viewport.vy + viewport.height / 2;
	return {x: newX, y: newY}
};
