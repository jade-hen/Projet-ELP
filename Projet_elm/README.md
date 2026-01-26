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



## Explication du flux du jeu

1. **Démarrage** : `init` lance `Words.loadWords` pour charger `static/words.txt`.
2. **Mots chargés** : `GotWords (Ok txt)` → `Words.parseWords txt` remplit `model.words`, puis `Words.chooseRandomIndex` tire un index.
3. **Mot choisi** : `PickedIndex i` → récupère le mot dans `model.words`, le met dans `target`, puis lance `Dictionary.fetchMeanings`.
4. **Définitions reçues** : `GotDefs (Ok meanings)` → stocke `meanings` et passe l’état à `Ready`.
5. **Jeu** : `GuessChanged` met à jour `guess` et si `normalize guess == normalize target` → état `Won` (+1 point si la solution n’a pas été affichée).  
   `NewGame` relance un tirage, `ToggleSolution` affiche/cache la solution.


## Erreurs possibles

- **`BadStatus: 404`** : fichier `words.txt` introuvable ou mot non trouvé par l’API.
- **`NetworkError`** : problème de connexion réseau / accès API.
- **`Timeout`** : la requête met trop de temps.
- **`BadBody: ...`** : réponse JSON inattendue / décodage impossible.



## Auteurs
Ce projet a été réalisé en groupe par :
- **Anna Grataloup**
- **Emma Payrard**
- **Jade Henninot**