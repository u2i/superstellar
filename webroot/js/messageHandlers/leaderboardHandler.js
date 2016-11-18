import {leaderboardDialog, globalState} from '../globals';

class Rank {
  constructor(position, name, score, currentClient = false) {
    this.position = position;
    this.name = name;
    this.score = score;
    this.currentClient = currentClient;
  }
}
const leaderboardHandler = (leaderboard) => {
  const ranks = leaderboard.ranks;

  let result = ranks.map(
    (rank, index) => new Rank(index + 1, globalState.clientIdToName.get(rank.id), rank.score, index + 1 === leaderboard.userPosition)
  );

  const currentUserRank = new Rank(leaderboard.userPosition, globalState.clientIdToName.get(leaderboard.clientId), leaderboard.userScore, true);

  leaderboardDialog.update(result, currentUserRank)
};

export default leaderboardHandler;
