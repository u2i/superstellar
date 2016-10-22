import {leaderboardDialog, globalState} from '../globals';

class Rank {
  constructor(rank, name, score) {
    this.rank = rank;
    this.name = name;
    this.score = score;
  }
}
const leaderboardHandler = (leaderboard) => {
  const ranks = leaderboard.ranks;

  let result = ranks.map(
    (rank, index) => new Rank(index + 1, globalState.clientIdToName.get(rank.id), rank.score)
  );

  leaderboardDialog.update(result)
};

export default leaderboardHandler;
