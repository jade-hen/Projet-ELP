// game.js
const { createInterface } = require("readline/promises");
const { stdin: input, stdout: output } = require("process");

const { createLogger } = require("./logger");
const { createDeckRuntime } = require("./deck");
const { makePlayer, playRound, isGameOver, getWinner } = require("./engine");

// Ajuste ici les quantités de cartes spéciales si tu veux coller exactement à ton deck physique
const DECK_CONFIG = {
  actions: { FREEZE: 3, FLIP_THREE: 3, SECOND_CHANCE: 3 },
  modifiers: { PLUS_2: 1, PLUS_4: 1, PLUS_6: 1, PLUS_8: 1, PLUS_10: 1, X2: 1 },
};

function makeId(prefix = "id") {
  return `${prefix}_${Math.random().toString(16).slice(2)}_${Date.now().toString(16)}`;
}

function isNo(s) {
  return ["n", "no", "non"].includes(String(s).trim().toLowerCase());
}

async function main() {
  const rl = createInterface({ input, output });

  const gameId = makeId("game");
  const logger = createLogger({});
  logger.log("GAME_START", { gameId });

  console.log("Flip7 (mode texte) - règles complètes");
  console.log(`Log: ${logger.filepath}\n`);

  let n;
  while (true) {
    const ans = await rl.question("Nombre de joueurs (2+): ");
    n = Number(ans);
    if (Number.isInteger(n) && n >= 2) break;
    console.log("Entrée invalide.");
  }

  const players = [];
  for (let i = 0; i < n; i++) {
    const name = (await rl.question(`Nom joueur ${i + 1}: `)).trim() || `J${i + 1}`;
    players.push(makePlayer(name));
  }
  logger.log("PLAYERS", { gameId, players: players.map((p) => p.name) });

  const deck = createDeckRuntime(DECK_CONFIG);

  let dealerIndex = 0;
  let roundNo = 1;

  while (true) {
    await playRound(players, deck, dealerIndex, rl, logger, gameId, roundNo);

    if (isGameOver(players)) break;

    dealerIndex = (dealerIndex + 1) % players.length;
    roundNo++;

    const cont = await rl.question("Continuer ? (o/n) ");
    if (isNo(cont)) break;
  }

  const winner = getWinner(players);
  logger.log("GAME_END", {
    gameId,
    winner: winner.name,
    totals: Object.fromEntries(players.map((p) => [p.name, p.total])),
  });

  console.log("\n=== FIN DE PARTIE ===");
  players.forEach((p) => console.log(`${p.name}: ${p.total}`));
  console.log(`Gagnant: ${winner.name}`);

  await rl.close();
}

main().catch((err) => {
  console.error("Erreur:", err);
  process.exitCode = 1;
});
