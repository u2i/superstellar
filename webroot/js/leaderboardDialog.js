export default class LeaderboardDialog {
  constructor () {
    this.domNode = document.getElementById('leaderboard');
  }

  show () {
    this.domNode.style.display = 'block'
  }

  hide() {
    this.domNode.style.display = 'none'
  }

  update(ranks) {
    let tbody = this.domNode.firstElementChild.firstElementChild;
    tbody.innerHTML = '';

    for(let rank of ranks) {
      let tr = document.createElement("tr");

      tr.appendChild(this.buildCell(rank.rank, 'rank'));
      tr.appendChild(this.buildCell(rank.name, 'name'));
      tr.appendChild(this.buildCell(rank.score, 'score'));

      tbody.appendChild(tr);
    }
  }

  buildCell(html, className) {
    let td = document.createElement("td");
    td.innerHTML = html;
    td.className = className;
    return td;
  }
}
