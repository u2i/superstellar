export default class LeaderboardDialog {
  constructor() {
    this.domNode = document.getElementById('leaderboard');
    this.leaderboard = document.getElementById('leaderboard-head').firstElementChild;
    this.tailLeaderboard = document.getElementById('leaderboard-tail');
    this.tailEllipsisRow = this.tailLeaderboard.firstElementChild.firstElementChild;
  }

  show() {
    this.domNode.style.display = 'block'
  }

  hide() {
    this.domNode.style.display = 'none'
  }

  update(ranks, currentUserRank) {
    this.fillHeadLeaderboard(ranks);
    this.fillTailLeaderboard(currentUserRank, ranks.length);
  }

  fillHeadLeaderboard(ranks) {
    const newCount = ranks.length;
    const oldCount = this.leaderboard.childElementCount;
    const maxCount = Math.max(oldCount, newCount);
    this.buildMissingRows(oldCount, newCount);
    this.fillRows(maxCount, newCount, ranks);
    return newCount;
  }

  buildMissingRows(oldCount, newCount) {
    for (let i = oldCount; i < newCount; i++) {
      const tr = LeaderboardDialog.buildEmptyRow(i + 1);
      this.leaderboard.appendChild(tr);
    }
  }

  static buildEmptyRow(position) {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td class='rank-position'>${position}</td>
      <td class='name'></td>
      <td class='score'></td>
    `;
    LeaderboardDialog.hideElement(tr);
    return tr;
  }

  fillRows(maxCount, newCount, ranks) {
    for (let i = 0; i < maxCount; i++) {
      const tr = this.leaderboard.children[i];
      if (i < newCount) {
        const rank = ranks[i];
        this.highlightCurrentUser(tr, rank);
        tr.children[1].textContent = rank.name;
        tr.children[2].textContent = rank.score;
        LeaderboardDialog.showElement(tr);
      } else {
        tr.children[1].textContent = '';
        tr.children[2].textContent = '';
        LeaderboardDialog.hideElement(tr);
      }
    }
  }

  highlightCurrentUser(tr, rank) {
    tr.className = rank.currentClient ? 'current' : '';
  }

  fillTailLeaderboard(currentUserRank, newCount) {
    const tr = document.getElementById('current-user-rank');
    tr.children[0].textContent = currentUserRank.position;
    tr.children[1].textContent = currentUserRank.name;
    tr.children[2].textContent = currentUserRank.score;

    this.showTailLeaderboard(currentUserRank, newCount);
  }

  showTailLeaderboard(currentUserRank, newCount) {
    if (currentUserRank.position <= newCount) {
      LeaderboardDialog.hideElement(this.tailLeaderboard)
    } else {
      LeaderboardDialog.showElement(this.tailLeaderboard);
      this.showTailLeaderboardEllipsis(currentUserRank, newCount);
    }
  }

  showTailLeaderboardEllipsis(currentUserRank, newCount) {
    if (currentUserRank.position === newCount + 1) {
      LeaderboardDialog.hideElement(this.tailEllipsisRow)
    } else {
      LeaderboardDialog.showElement(this.tailEllipsisRow)
    }
  }

  static showElement(element) {
    element.style.removeProperty('display')
  }

  static hideElement(element) {
    element.style.display = 'none'
  }
}
