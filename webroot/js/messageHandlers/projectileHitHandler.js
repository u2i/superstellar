import { globalState } from '../globals';

const projectileHitHandler = (message) => {
  let { id } = message;

  let projectile = globalState.projectilesMap.get(id)
  projectile.remove();
  globalState.projectilesMap.delete(projectile.id);
};

export default projectileHitHandler;
