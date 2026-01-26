# GuessIt — Projet Elm

## Présentation  
Il s’agit d’un jeu de devinette : l’utilisateur doit retrouver un mot à partir de définitions fournies par une API de dictionnaire.

## Lancement du projet

### Étapes pour exécuter l’application

1. Ouvrir un terminal et se placer dans le dossier du projet : `cd Projet_elm`

2. Lancer Elm Reactor : ``elm reactor``

3. Ouvrir un navigateur web et aller à l’adresse suivante : `http://localhost:8000`

4. Ouvrir le dossier src, puis cliquer sur Main.elm pour lancer l’application.

### Fonctionnement de l’application
Au démarrage, un mot est sélectionné aléatoirement depuis la liste. Les définitions du mot sont récupérées via une API externe. Ensuite, les définitions sont affichées.

L’utilisateur saisit une proposition et l’application indique si la réponse est correcte et attribue des points.

**Attention !** Le joueur n'obtient des points que s'il n'a pas affiché la réponse (impossible de tricher !). 

## Structure du projet
```
Projet_elm/
├─ src/
│ ├─ Main.elm
│ ├─ Types.elm
│ ├─ Dictionary.elm
│ └─ Words.elm
│
├─ static/
│ └─ words.txt
│
└─ elm.json
```
### Description des fichiers

- **Main.elm**  
  Point d’entrée de l’application. Contient la fonction `main`, ainsi que la logique principale (`init`, `update`, `view`).

- **Types.elm**  
  Définit les types principaux de l’application (`Model`, `Msg`, `Meaning`, etc.).

- **Dictionary.elm**  
  Gère les appels à l’API de dictionnaire et le décodage des réponses JSON.

- **Words.elm**  
  Gère la liste de mots, leur chargement depuis un fichier texte et le tirage aléatoire.

- **static/words.txt**  
  Fichier contenant les mots utilisés par le jeu.

