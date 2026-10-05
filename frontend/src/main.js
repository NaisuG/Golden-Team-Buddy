import './style.css';
import './board.css';

import { GenerateVariants, ListChampions } from '../wailsjs/go/main/App';

const MAX_BOARD = 10;
const MAX_BENCH = 9;
const BOARD_ROWS = [7, 7, 7, 7];

let allChampions = [];
let boardChampions = [];
let benchChampions = [];
let searchQuery = '';

document.querySelector('#app').innerHTML = `
  <div class="gtb-layout">
    <div id="board-grid" class="board-grid"></div>

    <div class="bench-label">Banca (<span id="bench-count">0</span>/${MAX_BENCH})</div>
    <div id="bench-row" class="bench-row"></div>

    <input id="champion-search" class="champion-search" type="text"
           placeholder="Buscar campeón por nombre" autocomplete="off" />
    <div id="champion-picker" class="champion-picker"></div>

    <button id="generate-btn" class="generate-btn">Generar variantes</button>
    <div id="variants-output" class="variants-output"></div>
  </div>
`;

const boardGrid = document.getElementById('board-grid');
const benchRow = document.getElementById('bench-row');
const championPicker = document.getElementById('champion-picker');
const benchCountEl = document.getElementById('bench-count');
const variantsOutput = document.getElementById('variants-output');
const searchInput = document.getElementById('champion-search');
const generateBtn = document.getElementById('generate-btn');

generateBtn.onclick = runGenerateVariants;

searchInput.addEventListener('input', (e) => {
  searchQuery = e.target.value;
  renderPicker();
});

boardGrid.addEventListener('dragover', (e) => e.preventDefault());
boardGrid.addEventListener('drop', (e) => {
  e.preventDefault();
  const key = e.dataTransfer.getData('text/plain');
  if (key) placeChampion(key, 'board');
});

benchRow.addEventListener('dragover', (e) => e.preventDefault());
benchRow.addEventListener('drop', (e) => {
  e.preventDefault();
  const key = e.dataTransfer.getData('text/plain');
  if (key) placeChampion(key, 'bench');
});

function nameForKey(key) {
  const champ = allChampions.find(c => c.key === key);
  return champ ? champ.name : key;
}

function isPlaced(key) {
  return boardChampions.includes(key) || benchChampions.includes(key);
}

function placeChampion(key, target) {
  const targetList = target === 'board' ? boardChampions : benchChampions;
  const sourceList = target === 'board' ? benchChampions : boardChampions;
  const max = target === 'board' ? MAX_BOARD : MAX_BENCH;

  if (targetList.includes(key) || targetList.length >= max) return;

  const sourceIndex = sourceList.indexOf(key);
  if (sourceIndex !== -1) sourceList.splice(sourceIndex, 1);

  targetList.push(key);
  render();
}

function removeFromBoard(index) {
  boardChampions.splice(index, 1);
  render();
}

function removeFromBench(index) {
  benchChampions.splice(index, 1);
  render();
}

function filteredChampions() {
  const q = searchQuery.trim().toLowerCase();
  if (!q) return allChampions;
  return allChampions.filter(c => c.name.toLowerCase().includes(q));
}

function makeDraggable(el) {
  el.addEventListener('dragstart', (e) => {
    e.dataTransfer.setData('text/plain', el.dataset.key);
  });
}

function renderPicker() {
  const list = filteredChampions();
  championPicker.innerHTML = list.map(champ => {
    const placed = isPlaced(champ.key);
    return `<button class="champion-chip" data-key="${champ.key}" draggable="${!placed}" ${placed ? 'disabled' : ''}>${champ.name}</button>`;
  }).join('');

  championPicker.querySelectorAll('.champion-chip:not(:disabled)').forEach(btn => {
    btn.onclick = () => placeChampion(btn.dataset.key, 'board');
    makeDraggable(btn);
  });
}

function hexHTML(key, index) {
  return key
    ? `<div class="hex filled" draggable="true" data-key="${key}" data-index="${index}">${nameForKey(key)}</div>`
    : `<div class="hex"></div>`;
}

function renderBoard() {
  let html = '';
  let slot = 0;
  BOARD_ROWS.forEach((count, rowIndex) => {
    html += `<div class="hex-row ${rowIndex % 2 === 1 ? 'offset' : ''}">`;
    for (let i = 0; i < count; i++) {
      html += hexHTML(boardChampions[slot], slot);
      slot++;
    }
    html += `</div>`;
  });
  boardGrid.innerHTML = html;
  boardGrid.querySelectorAll('.hex.filled').forEach(el => {
    el.onclick = () => removeFromBoard(parseInt(el.dataset.index, 10));
    makeDraggable(el);
  });
}

function renderBench() {
  let html = '';
  for (let i = 0; i < MAX_BENCH; i++) {
    html += hexHTML(benchChampions[i], i);
  }
  benchRow.innerHTML = html;
  benchRow.querySelectorAll('.hex.filled').forEach(el => {
    el.onclick = () => removeFromBench(parseInt(el.dataset.index, 10));
    makeDraggable(el);
  });
  benchCountEl.innerText = benchChampions.length;
}

function render() {
  renderBoard();
  renderBench();
  renderPicker();
}

function showMessage(text, isError = false) {
  variantsOutput.innerHTML = `<p class="variants-message ${isError ? 'error' : ''}"></p>`;
  variantsOutput.firstElementChild.textContent = text;
}

function variantHTML(variant, label, isChild) {
  const champions = variant.Champions.map(nameForKey).join(', ');
  const traits = (variant.Traits || [])
    .map(t => `<span class="trait-chip ${t.Style}">${t.Count} ${t.Name}</span>`)
    .join('');
  return `
    <div class="variant ${isChild ? 'child' : ''}">
      <div class="variant-champions"><strong>${label})</strong> ${champions}</div>
      <div class="variant-traits">${traits}</div>
    </div>`;
}

function renderVariants(variants) {
  if (!variants || variants.length === 0) {
    showMessage('Sin variantes: revisa que haya al menos 1 campeón en el tablero.');
    return;
  }
  variantsOutput.innerHTML = variants.map((v, i) => {
    const children = (v.Children || []).map((child, j) => {
      const letter = String.fromCharCode('a'.charCodeAt(0) + j);
      return variantHTML(child, `${i + 1}${letter}`, true);
    }).join('');
    return variantHTML(v, `${i + 1}`, false) + children;
  }).join('');
}

async function runGenerateVariants() {
  generateBtn.disabled = true;
  generateBtn.innerText = 'Generando...';
  try {
    const variants = await GenerateVariants(boardChampions, benchChampions);
    renderVariants(variants);
  } catch (err) {
    showMessage('Error: ' + err, true);
    console.error(err);
  } finally {
    generateBtn.disabled = false;
    generateBtn.innerText = 'Generar variantes';
  }
}

async function init() {
  try {
    allChampions = await ListChampions();
  } catch (err) {
    showMessage('No se pudo cargar la lista de campeones: ' + err, true);
    console.error(err);
  }
  render();
}

init();
