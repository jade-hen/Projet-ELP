// deck.js
function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [arr[i], arr[j]] = [arr[j], arr[i]];
  }
  return arr;
}

function buildDeck(config) {
  const cards = [];

  // 0
  cards.push({ type: "NUMBER", value: 0 });

  // 1..12 : v exemplaires de v
  for (let v = 1; v <= 12; v++) {
    for (let k = 0; k < v; k++) cards.push({ type: "NUMBER", value: v });
  }

  // Actions
  for (const [name, count] of Object.entries(config.actions || {})) {
    for (let i = 0; i < count; i++) cards.push({ type: "ACTION", name });
  }

  // Modificateurs
  for (const [k, n] of Object.entries(config.modifiers || {})){
    for (let i = 0; i < n; i++){
      cards.push({type: "MODIFIER", kind: k.startsWith("PLUS_") ? "PLUS" : "X2",
                  value: +(k.startsWith("PLUS_") ? k.slice(5):k.slice(1))});
    }
  }

  return shuffle(cards);
}

function createDeckRuntime(config) {
  return {
    draw: buildDeck(config),
    discard: [], //la défausse
  };
}

function discardCard(deck, card) {
  deck.discard.push(card);
}

function reshuffleIfNeeded(deck, logger, meta) {
  if (deck.draw.length > 0) return;
  if (deck.discard.length === 0)
    throw new Error("Plus de cartes (draw et discard vides).");

  deck.draw = shuffle(deck.discard);
  deck.discard = [];

  logger.log("RESHUFFLE", { ...meta, newDrawSize: deck.draw.length });
  console.log("-> Défausse utilisée et pioche mélangée");
}

function drawCard(deck, logger, meta) {
  reshuffleIfNeeded(deck, logger, meta);
  return deck.draw.pop();
}

function cardToString(card) {
  if (card.type === "NUMBER") return `${card.value}`;
  if (card.type === "ACTION") return card.name.split("_").join(" "); //enlever les _
  if (card.type === "MODIFIER") return card.kind === "X2" ? "x2" : `+${card.value}`;
  return "UNKNOWN";
}

module.exports = {
  createDeckRuntime,
  drawCard,
  discardCard,
  cardToString,
};
