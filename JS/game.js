const { createInterface } = require("readline/promises");
const { stdin: input, stdout: output } = require("process");

const { createLogger } = require("./logger");
const { createDeckRuntime } = require("./deck");
const { makePlayer, playRound, isGameOver, getWinner } = require("./engine");

const DECK_CONFIG = {
  actions: { FREEZE: 3, FLIP_THREE: 3, SECOND_CHANCE: 3 },
  modifiers: { PLUS_2: 1, PLUS_4: 1, PLUS_6: 1, PLUS_8: 1, PLUS_10: 1, X2: 1 },
};

function isNo(s) {
  return ["n", "no", "non"].includes(String(s).trim().toLowerCase());
}

async function main() {
  const rl = createInterface({ input, output }); // rl = readline

  //const gameId = makeId("game");
  const logger = createLogger({});
  logger.log("GAME_START"); //, { gameId }

  console.log("Flip7");
  console.log(`Log: ${logger.filepath}\n`);

  //Nombre de joueurs
  let n;
  while (true) {
    const ans = await rl.question("Nombre de joueurs (2+): ");
    n = Number(ans);
    if (Number.isInteger(n) && n >= 2) break;
    console.log("Entrée invalide.");
  }

  //Noms des joueurs
  const players = [];
  for (let i = 0; i < n; i++) {
    const name = (await rl.question(`Nom joueur ${i + 1}: `)).trim() || `J${i + 1}`;
    players.push(makePlayer(name));
  }
  logger.log("PLAYERS", { players: players.map((p) => p.name) }); 

  nbPaquets = Math.floor(n/(18))+(n%18!=0)
  if (n>18) {
    console.log("\nIl y a plus de 18 joueurs => Nécessité de rajouter", nbPaquets-1, nbPaquets-1>1 ? "paquets":"paquet");
  }
  logger.log("CREATE_DECK", { decksNumber: nbPaquets })

  const deck = createDeckRuntime(DECK_CONFIG, n); // rajouter le nombre de joueur pour savoir avec combien de paquets de cartes il faut jouer (+ ajouter un print pour montrer que ça marche bien)

  let dealerIndex = 0;
  let roundNo = 1;
  while (true) {
    await playRound(players, deck, dealerIndex, rl, logger, roundNo); //gameId, 

    if (isGameOver(players)) break;

    dealerIndex = (dealerIndex + 1) % players.length;
    roundNo++;

    const cont = await rl.question("Continuer ? (o/n) ");//finir plus tôt
    if (isNo(cont)) break;
  }

  const winner = getWinner(players);
  logger.log("GAME_END", {winner: winner.name, //gameId, 
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
