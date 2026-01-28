// deck.js
function shuffle(arr) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [arr[i], arr[j]] = [arr[j], arr[i]];
  }
  return arr;
}

// Nombres: 12x12, 11x11, ... 1x1, + un 0
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
  const map = {
    PLUS_2: { type: "MODIFIER", kind: "PLUS", value: 2 },
    PLUS_4: { type: "MODIFIER", kind: "PLUS", value: 4 },
    PLUS_6: { type: "MODIFIER", kind: "PLUS", value: 6 },
    PLUS_8: { type: "MODIFIER", kind: "PLUS", value: 8 },
    PLUS_10: { type: "MODIFIER", kind: "PLUS", value: 10 },
    X2: { type: "MODIFIER", kind: "X2", value: 2 },
  };

  for (const [key, count] of Object.entries(config.modifiers || {})) {
    for (let i = 0; i < count; i++) cards.push({ ...map[key] });
  }

  return shuffle(cards);
}

function createDeckRuntime(config) {
  return {
    draw: buildDeck(config),
    discard: [],
  };
}

function discardCard(deck, card) {
  deck.discard.push(card);
}

function reshuffleIfNeeded(deck, logger, meta) {
  if (deck.draw.length > 0) return;
  if (deck.discard.length === 0) throw new Error("Plus de cartes (draw et discard vides).");

  deck.draw = shuffle(deck.discard);
  deck.discard = [];
  logger.log("RESHUFFLE", { ...meta, newDrawSize: deck.draw.length });
}

function drawCard(deck, logger, meta) {
  reshuffleIfNeeded(deck, logger, meta);
  return deck.draw.pop();
}

function cardToString(card) {
  if (card.type === "NUMBER") return `${card.value}`;
  if (card.type === "ACTION") return card.name === "FLIP_THREE" ? "FLIP THREE" : card.name;
  if (card.type === "MODIFIER") return card.kind === "X2" ? "x2" : `+${card.value}`;
  return "UNKNOWN";
}

module.exports = {
  createDeckRuntime,
  drawCard,
  discardCard,
  cardToString,
};
