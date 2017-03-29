export default class ScoreBoardDialog {
  constructor() {
    this.domNode = document.getElementById('scoreboard');
    this.scoreBoard = document.getElementById('scoreboard-items').firstElementChild;
  }

  show() {
    this.domNode.style.display = 'block'
  }

  hide() {
    this.domNode.style.display = 'none'
  }

  update(items) {
    this.createRows(items);
    this.fillRows(items);
  }

  createRows(items) {
    const newCount = Math.min(items.length, 5);
    const actCount = this.scoreBoard.childElementCount;

    for(let i = actCount; i < newCount; i++) {
      const tr = ScoreBoardDialog.buildRow(i + 1);
      this.scoreBoard.appendChild(tr);
    }
  }

  static buildRow(position) {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td class='rank-position'>${position}</td>
      <td class='name'></td>
      <td class='score'></td>
    `;
    ScoreBoardDialog.hideElement(tr);
    return tr;
  }

  fillRows(items) {
    for (let i = 0; i < items.length; i++) {
      const tr = this.scoreBoard.children[i];
      const item = items[i];
      tr.children[1].textContent = item.name;
      tr.children[2].textContent = item.score;
      ScoreBoardDialog.showElement(tr);
    }
  }

  static showElement(element) {
    element.style.removeProperty('display')
  }

  static hideElement(element) {
    element.style.display = 'none'
  }
}
