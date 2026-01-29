# JS Game

Ce projet contient un jeu développé en JavaScript pour ``Node.js``. Il ne dépend d’aucun module externe et peut être exécuté directement avec Node.

## Prérequis

- ``Node.js ``

## Installation

Aucune installation n’est nécessaire. Le projet ne requiert aucun module externe.

## Exécution

Lancer le jeu via la commande ``node game.js`` ou bien `npm start`.

## Comment jouer

Après avoir lancé le jeu, il faut rentrer un nombre de joueurs. On peut ensuite saisir le nom de chaque joueur un par un ou bien laisser vide pour un nom automatique (J1, J2, ...).

Le jeu distribue ensuite une carte à chaque joueur pour le premier tour. Si une action doit être réalisée, le jeu demande au joueur concerné à qui il veut l'appliquer et il doit choisir un numéro de joueur parmi ceux proposés. 

Après ce premier tour de distribution, le joueur "dealer" (celui qui distribue) demande à chacun s'il veut piocher une carte. Taper `h` ou `s` respectivement pour piocher (hit) ou s'arrêter (stand). 

A la fin de chaque tour, vous pouvez choisir d'arrêter le jeu. Le joueur ayant le plus de points est alors désigné gagnant. Sinon, le jeu s'arrête lorsqu'au moins un joueur atteint 200 points.

## Structure du projet

```
JS/
├── deck.js
├── engine.js
├── game.js
├── logger.js
├── package.json
├── README.md
└── logs/
```

## Description des fichiers

| Fichier      | Description |
|-------------|-------------|
| game.js   | Point d’entrée du jeu |
| engine.js | Moteur de jeu (logique principale) |
| deck.js   | Gestion du paquet / cartes |
| logger.js | Gestion des logs |
| logs/     | Répertoire contenant les fichiers de logs générés |
