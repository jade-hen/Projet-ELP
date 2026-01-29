// engine.js
const { drawCard, discardCard, cardToString } = require("./deck");

const TARGET_SCORE = 200;
const FLIP7_BONUS = 15;

function makeId(prefix = "id") {//id pour quoi ? besoin que ce soit aussi compliqué ?
  return `${prefix}_${Math.random().toString(16).slice(2)}_${Date.now().toString(16)}`;
}

function makePlayer(name) { // pourquoi besoin d'un id ? pour identifier de façon unique un joueur (pas juste avec le nom)
  return { id: makeId("p"), name, total: 0, round: null };
}

function initRoundState(player) {//état d'un joueur
  player.round = {
    active: true,
    stood: false,
    numbers: [],
    modifiersPlus: 0,
    hasX2: false,
    hasSecondChance: false,
    flip7: false,
    bustedByDuplicate: false, //éliminé par une carte en double
    bustedByFreeze: false, //éliminé par une carte freeze
  };
}

function countUniqueNumbers(player) {// pas obligé de faire une fonction juste pour ça ?
  return new Set(player.round.numbers).size;
}

function computeRoundScore(player) { // calculer le score d'un tour pour un joueur
  const r = player.round;
  if (!r) return 0;

  // PRIORITÉ: si Freeze ou doublon => 0 quoi qu'il arrive
  if (r.bustedByFreeze || r.bustedByDuplicate) return 0;

  const sumNumbers = r.numbers.reduce((a, b) => a + b, 0); //réduire le tableau en appliquant la fonction => fait la somme de tous les éléments
  const numbersPart = r.hasX2 ? sumNumbers * 2 : sumNumbers; // x2 ne double pas les bonus
  let score = numbersPart + r.modifiersPlus;
  if (r.flip7) score += FLIP7_BONUS;
  return score;
}

function checkFlip7AndMaybeEndRound(player, roundCtx, logger, meta) {
  if (countUniqueNumbers(player) >= 7) {
    player.round.flip7 = true;
    roundCtx.ended = true;
    roundCtx.endReason = "FLIP7";
    roundCtx.flip7By = player.name;
    logger.log("FLIP7", { ...meta, player: player.name, bonus: FLIP7_BONUS });
  }
}

function applyModifierCard(player, card, logger, meta) {
  const r = player.round;
  if (card.kind === "X2") {
    r.hasX2 = true;
    logger.log("ADD_X2", { ...meta, player: player.name }); //utile ?
  } else {
    r.modifiersPlus += card.value;
    logger.log("ADD_PLUS", { ...meta, player: player.name, plus: card.value, totalPlus: r.modifiersPlus }); //utile ?
  }
}

function applyNumberCard(player, value, deck, logger, meta) {
  const r = player.round;

  if (r.numbers.includes(value)) {//si le joueur a déjà la carte
    if (r.hasSecondChance) {
        r.hasSecondChance = false;
        logger.log("SECOND_CHANCE_USED", { ...meta, player: player.name, duplicateValue: value });
        return { ok: true, duplicate: true, usedSecondChance: true };
    }
    //pas de seconde chance, joueur éliminé
    r.active = false;
    r.bustedByDuplicate = true;
    logger.log("BUST_DUPLICATE", { ...meta, player: player.name, value });
    return { ok: false, duplicate: true, usedSecondChance: false };
  }

  r.numbers.push(value);
  logger.log("ADD_NUMBER", { ...meta, player: player.name, value, numbers: [...r.numbers] });
  return { ok: true, duplicate: false, usedSecondChance: false };
}

function freezePlayer(target, logger, meta) {
  // Freeze = éliminé du tour + score du tour = 0
  target.round.active = false;
  target.round.stood = true;          // pour être sûr qu’il ne rejoue pas
  target.round.bustedByFreeze = true; // marqueur explicite

  logger.log("FREEZE_APPLIED", { ...meta, target: target.name });
}


async function chooseTarget(activePlayers, currentPlayer, rl) {//choisir une cible pour appliquer la carte 
  if (activePlayers.length === 1) return activePlayers[0];

  while (true) {
    console.log("\nChoisis une cible :");
    activePlayers.forEach((p, idx) => {
      console.log(`  ${idx + 1}) ${p.name}${p.id === currentPlayer.id ? " (toi)" : ""}`);
    });
    const ans = await rl.question("> numéro: ");
    const n = Number(ans);
    if (Number.isInteger(n) && n >= 1 && n <= activePlayers.length) return activePlayers[n - 1];
    console.log("Entrée invalide.");
  }
}

async function handleSecondChanceDraw(receiver, allPlayers, deck, logger, meta, rl) {
  const r = receiver.round;
  if (!r.hasSecondChance) {
    r.hasSecondChance = true;
    logger.log("SECOND_CHANCE_TAKEN", { ...meta, player: receiver.name });
    return;
  }

  const active = allPlayers.filter((p) => p.round.active && !p.round.stood);
  const eligible = active.filter((p) => !p.round.hasSecondChance);//tous ceux qui n'ont pas de seconde chance

  if (eligible.length === 0) {
    discardCard(deck, { type: "ACTION", name: "SECOND_CHANCE", virtual: false });
    logger.log("SECOND_CHANCE_DISCARDED_NO_ELIGIBLE", { ...meta, from: receiver.name });
    return;
  }

  console.log(`\n${receiver.name} a déjà une Second Chance. Il/elle doit la donner.`);
  const target = await chooseTarget(eligible, receiver, rl);
  target.round.hasSecondChance = true;
  logger.log("SECOND_CHANCE_GIVEN", { ...meta, from: receiver.name, to: target.name });
}

async function resolveActionCard(card, sourcePlayer, allPlayers, deck, roundCtx, rl, logger, meta, opts = {}) {
  const activePlayers = allPlayers.filter((p) => p.round.active && !p.round.stood);

  if (card.name === "FREEZE") {
    const target = await chooseTarget(activePlayers, sourcePlayer, rl);
    logger.log("ACTION_FREEZE_TARGET", { ...meta, from: sourcePlayer.name, to: target.name, duringFlipThree: !!opts.duringFlipThree });
    freezePlayer(target, logger, meta); //plein de lignes qu'on peut mettre direct dans les autres fonctions, non ?
    return;
  }

  if (card.name === "SECOND_CHANCE") {
    await handleSecondChanceDraw(sourcePlayer, allPlayers, deck, logger, meta, rl);
    logger.log("SECOND_CHANCE_REPLACEMENT_DRAW", { ...meta, player: sourcePlayer.name });
    await drawAndResolveForPlayer(sourcePlayer, allPlayers, deck, roundCtx, rl, logger, meta); //pioche à nouveau
    return;
  }

  if (card.name === "FLIP_THREE") {
    const target = await chooseTarget(activePlayers, sourcePlayer, rl);
    logger.log("ACTION_FLIP_THREE_TARGET", { ...meta, from: sourcePlayer.name, to: target.name, duringFlipThree: !!opts.duringFlipThree });
    await performFlipThree(target, allPlayers, deck, roundCtx, rl, logger, meta);
    return;
  }

  logger.log("ACTION_UNKNOWN", { ...meta, from: sourcePlayer.name, card });
}

async function performFlipThree(target, allPlayers, deck, roundCtx, rl, logger, meta) {
  console.log(`\n>>> FLIP THREE sur ${target.name}: il/elle doit piocher 3 cartes.`);
  const pendingActions = [];

  for (let i = 1; i <= 3; i++) {
    if (roundCtx.ended) break;

    const card = drawCard(deck, logger, meta); //pioche
    roundCtx.tableCards.push(card); 

    logger.log("FLIP_THREE_DRAW", { ...meta, target: target.name, index: i, card: { ...card } });
    console.log(`${target.name} (FlipThree) pioche ${i}/3: ${cardToString(card)}`);

    if (card.type === "NUMBER") {
      if (target.round.active) {
        const res = applyNumberCard(target, card.value, deck, logger, meta);
        if (res.ok) checkFlip7AndMaybeEndRound(target, roundCtx, logger, meta);
      } else {
        logger.log("IGNORED_CARD_TARGET_ALREADY_OUT", { ...meta, target: target.name, card });
      }
      continue;
    }

    if (card.type === "MODIFIER") {
      if (target.round.active) applyModifierCard(target, card, logger, meta);
      else logger.log("IGNORED_MODIFIER_TARGET_ALREADY_OUT", { ...meta, target: target.name, card });
      continue;
    }

    if (card.type === "ACTION") pendingActions.push(card);
  }

  if (roundCtx.ended) return;

  for (const actionCard of pendingActions) {
    if (roundCtx.ended) break;
    console.log(`\n>>> Résolution différée (FlipThree) pour ${target.name}: ${cardToString(actionCard)}`);
    await resolveActionCard(actionCard, target, allPlayers, deck, roundCtx, rl, logger, meta, { duringFlipThree: true });
  }
}

async function drawAndResolveForPlayer(player, allPlayers, deck, roundCtx, rl, logger, meta) {
  const card = drawCard(deck, logger, meta);
  roundCtx.tableCards.push(card); // on garde la carte sur la table pour la défausse de fin de tour

  logger.log("DRAW", { ...meta, player: player.name, card: { ...card } });
  console.log(`${player.name} pioche: ${cardToString(card)}`);

  if (card.type === "NUMBER") {
    const res = applyNumberCard(player, card.value, deck, logger, meta);
    if (res.ok) checkFlip7AndMaybeEndRound(player, roundCtx, logger, meta);
    return;
  }

  if (card.type === "MODIFIER") {
    applyModifierCard(player, card, logger, meta);
    return;
  }

  await resolveActionCard(card, player, allPlayers, deck, roundCtx, rl, logger, meta);
}

function showTable(players) {
  console.log("\n--- TABLE ---");
  for (const p of players) {
    const r = p.round;
    const nums = r.numbers.length ? r.numbers.join(",") : "-";
    const uniq = countUniqueNumbers(p);
    const mods = [];
    if (r.hasX2) mods.push("x2");
    if (r.modifiersPlus) mods.push(`+${r.modifiersPlus}`);
    if (r.hasSecondChance) mods.push("SecondChance");
    const status = !r.active ? "OUT" : r.stood ? "STAND" : "ACTIVE";
    console.log(`${p.name} [${status}] nums(${uniq} uniques): ${nums} | mods: ${mods.length ? mods.join(" | ") : "-"}`);
  }
  console.log("-------------\n");
}

async function roundLoop(players, deck, roundCtx, rl, logger, meta) {
  firstDeal = true;
  while (!roundCtx.ended) {
    const active = players.filter((p) => p.round.active && !p.round.stood);
    if (active.length === 0) {
      roundCtx.ended = true;
      roundCtx.endReason = "NO_ACTIVE";
      break;
    }

    for (const p of players) {
      if (roundCtx.ended) break;
      if (!p.round.active || p.round.stood) continue;
      
      if (!firstDeal){ // si ce n'est pas le premier tour, on laisse le choix au joueur de rester ou pas
        showTable(players);
        console.log(`${p.name}, total: ${p.total}. Potentiel tour: ${computeRoundScore(p)}`);
        console.log("Choix: (h)it = recevoir une carte, (s)tand = rester");
      }
      while (true) {
        let ans = "";
        if (!firstDeal){ // on demande au joueur
          ans = (await rl.question("> ")).trim().toLowerCase();
        } 
        if (firstDeal || ans === "h" || ans === "hit") {
          if (!firstDeal){
            logger.log("CHOICE", { ...meta, player: p.name, choice: "HIT" });
          }
          await drawAndResolveForPlayer(p, players, deck, roundCtx, rl, logger, meta);
          break;
        }
        if (ans === "s" || ans === "stand") {
          p.round.stood = true;
          logger.log("CHOICE", { ...meta, player: p.name, choice: "STAND" });
          break;
        }
        console.log("Entrée invalide. Tape h ou s.");
      }
    }
    firstDeal = false;
  }
}

function discardEndOfRoundSecondChance(players, deck, logger, meta) {
  for (const p of players) {
    if (p.round.hasSecondChance) {
      p.round.hasSecondChance = false;
      logger.log("SECOND_CHANCE_DISCARDED_END_ROUND", { ...meta, player: p.name });
    }
  }
}

async function playRound(players, deck, dealerIndex, rl, logger, gameId, roundNo) {
  players.forEach(initRoundState);

  const roundCtx = {
  roundNo,
  ended: false,
  endReason: null,
  flip7By: null,
  tableCards: [], // toutes les cartes révélées pendant ce tour
    };

  const meta = { gameId, round: roundNo };

  const dealer = players[dealerIndex];
  logger.log("ROUND_START", { ...meta, dealer: dealer.name });

  console.log(`\n====================\nTOUR #${roundNo} (donneur: ${dealer.name})\n====================\n`);

  //await initialDeal(players, deck, roundCtx, rl, logger, meta); // à enlever ?
  if (!roundCtx.ended) await roundLoop(players, deck, roundCtx, rl, logger, meta);

  discardEndOfRoundSecondChance(players, deck, logger, meta);

  const roundScores = {};
  for (const p of players) {
    const s = computeRoundScore(p);
    roundScores[p.name] = s;
    p.total += s;
  }
  // Fin de tour : toutes les cartes révélées vont en défausse
  for (const c of roundCtx.tableCards) {
    discardCard(deck, c);
  }
  logger.log("MOVE_TABLE_TO_DISCARD", { ...meta, count: roundCtx.tableCards.length });
  roundCtx.tableCards = [];

  logger.log("ROUND_END", {
    ...meta,
    reason: roundCtx.endReason,
    flip7By: roundCtx.flip7By,
    roundScores,
    totals: Object.fromEntries(players.map((p) => [p.name, p.total])),
  });

  console.log("\n--- FIN DE TOUR ---");
  console.log(`Raison: ${roundCtx.endReason}${roundCtx.flip7By ? ` (Flip7 par ${roundCtx.flip7By})` : ""}`);
  for (const p of players) console.log(`${p.name}: +${roundScores[p.name]} => total ${p.total}`);
  console.log("-------------------\n");

  return { roundCtx, roundScores };
}

function isGameOver(players) {
  return players.some((p) => p.total >= TARGET_SCORE);
}

function getWinner(players) {
  let best = players[0];
  for (const p of players) if (p.total > best.total) best = p;
  return best;
}

module.exports = {
  makePlayer,
  playRound,
  isGameOver,
  getWinner,
  TARGET_SCORE,
};
