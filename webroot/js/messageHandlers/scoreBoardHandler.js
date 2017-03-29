import {scoreBoardDialog} from '../globals';

class ScoreBoardItem {
  constructor(position, name, score) {
    this.position = position;
    this.name = name;
    this.score = score;
  }
}
const scoreBoardHandler = (scoreBoard) => {
  const items = scoreBoard.items;

  let result = items.map((rank, index) => new ScoreBoardItem(index + 1, rank.name, rank.score));
  console.log(result)
  scoreBoardDialog.update(items)
};

export default scoreBoardHandler;
