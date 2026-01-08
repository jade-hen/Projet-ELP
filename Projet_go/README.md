# Concurrent Name Matching Server (Go)

## Description

Ce projet est une preuve de concept (PoC) développée en langage Go, dont l’objectif est de mettre en pratique les concepts de concurrence et de programmation réseau à travers un cas concret : la détection de doublons et de correspondances approximatives entre des noms issus de fichiers CSV.

L’application repose sur l’algorithme de distance de Levenshtein, permettant de mesurer la similarité entre deux chaînes de caractères, et explore deux approches :
- une approche séquentielle
- une approche concurrente, exploitant les goroutines et les channels

Le projet inclut également un serveur TCP concurrent, capable de traiter plusieurs clients simultanément.

> Note : l'utilisation des fichiers et les commandes à utiliser sont données dans la suite pour chaque module exécutable. 

## Équipe

Ce projet est réalisé en groupe par :

- **Anna Grataloup**
- **Emma Payrard**
- **Jade Henninot**

## Contexte

Ce projet s’inspire d’un problème réel de qualité des données et d’identification des personnes, notamment mis en évidence dans le contexte des « faux positifs » en Colombie.

Dans ce scandale, des erreurs d’identification ont conduit à des confusions entre individus, souvent à cause de variations ou d’incohérences dans les noms (différences d’orthographe, accents, abréviations, erreurs de saisie ou sources multiples). Ces situations illustrent les limites des comparaisons strictes basées uniquement sur l’égalité exacte des chaînes de caractères.

L’objectif de ce projet est de montrer comment des techniques peuvent être utilisées pour détecter des doublons ou des correspondances potentielles dans une base de données de noms, tout en tenant compte des variations textuelles.

Le projet met également en évidence les enjeux de performance liés à ce type de traitement sur de grands volumes de données, et l’intérêt de la concurrence pour rendre ces analyses exploitables.


## Objectifs du projet

- Prendre en main le langage Go
- Implémenter un algorithme de matching approximatif (distance de Levenshtein)
- Exploiter la concurrence pour améliorer les performances
- Mettre en place un serveur TCP concurrent
- Comparer une approche séquentielle et une approche concurrente

## Fonctionnalités principales

- Calcul de la distance de Levenshtein entre deux chaînes de caractères
- Détection de doublons et de correspondances approximatives avec un seuil configurable
- Traitement concurrent des comparaisons grâce aux goroutines
- Serveur TCP capable de gérer plusieurs clients simultanément
- Mesure et analyse des performances


## Architecture générale

Le projet est structuré en plusieurs modules :
- data : chargement et parsing des fichiers CSV
- levenshtein : calcul de la distance entre deux chaînes
- matcher : détection des correspondances (séquentielle et concurrente)
- sanity : programme de validation et de comparaison des performances
- client TCP : envoi de données au serveur
- serveur TCP : réception, traitement concurrent et réponse

Cette séparation permet une meilleure lisibilité et une réutilisation des composants dans différents contextes (ligne de commande, réseau).

## Chargement des données (module data)

Le module data est responsable du chargement des noms depuis un fichier CSV ou depuis un flux générique (io.Reader).

Deux fonctions principales sont fournies :
- LoadNamesAndDates(path string) : charge un CSV depuis un fichier
- LoadNamesAndDatesFromReader(r io.Reader) : charge un CSV depuis un flux (utilisé notamment par le serveur TCP)

Les lignes vides ou valeurs vides sont ignorées, et la première ligne est sautée lorsqu’elle correspond à un en-tête.

Le CSV est chargé intégralement en mémoire à l’aide de encoding/csv.ReadAll(), ce qui simplifie l’implémentation dans le cadre de ce projet (pas besoin de lire les lignes une par une directement dans le csv).

Les fonctions retournent une liste de structures Person, qui correspondent à un nom et une date. 

## Distance de Levenshtein (module levenshtein)

La distance de Levenshtein est implémentée via une approche de programmation dynamique.

Principes : 
- Les chaînes sont converties en []rune afin de gérer le Unicode (accents et caractères non ASCII). Ce n'est pas utile dans le cas du fichier fourni (les noms sont en majuscules non accentuées), mais ça pourrait le devenir suivant les fichiers que les clients envoient au serveur. 
- Les opérations autorisées sont :
   - insertion
   - suppression
   - substitution
- Chaque opération a un coût de 1.

Optimisation mémoire : l'algorithme ne conserve que deux lignes de la matrice de calcul (prev et curr), ce qui réduit la complexité mémoire à O(m), où m est la longueur de la seconde chaîne.

Complexité : 
- Temps : O(n × m)
- Mémoire : O(m)

## Détection des correspondances (module matcher)

Utilisation de la date : le paramètre useDate permet de choisir si l'on souhaite comparer les dates ou non. En effet, deux personnes pourraient avoir un nom similaire, l'une d'entre elles serait alors supprimée à tort. La comparaison des dates évite ces "faux positifs". Note : on regarde ici si les dates sont exactement égales. On part du principe qu'il n'y a pas à la fois une erreur sur le nom et sur la date. 

Le module matcher fournit deux implémentations du matching :

7.1 Approche séquentielle
La fonction FindMatchesSequential compare chaque paire de noms (i, j) avec j > i, afin d’éviter les doublons et les auto-comparaisons.

Pour chaque paire :
- la distance de Levenshtein est calculée
- si la distance est inférieure ou égale au seuil (threshold), la paire est conservée dans les matches
- si useDate est à vrai, on compare aussi les dates

Les résultats sont ensuite triés par :
- distance croissante
- ordre alphabétique des noms

7.2 Approche concurrente
La fonction FindMatchesConcurrent implémente un worker pool :
- Un channel jobs distribue les paires (i, j) à comparer
- Un ensemble de goroutines (workers) consomme ces jobs
- Chaque worker calcule la distance et envoie les matches valides dans un channel results
- Un sync.WaitGroup garantit la synchronisation et la fermeture correcte des channels

Le nombre de workers est configurable ; s’il n’est pas fourni, il est automatiquement fixé au nombre de cœurs CPU (runtime.NumCPU()).

Limitation du volume : un paramètre limit permet de restreindre le nombre de noms traités, afin de contrôler la complexité O(n²) inhérente au problème (pour les tests notamment)

## Programme de validation et tests (sanity)

> Utilisation (dans le dossier Projet_go) : `go run cmd/sanity/main.go <chemin vers le fichier .csv de données>`

Le programme sanity permet de :
- vérifier le bon fonctionnement de l’algorithme de Levenshtein sur des données réelles
- comparer les temps d’exécution des versions séquentielle et concurrente
- comparer le nombre de matches suivant si on prend uniquement les noms en compte ou bien si l'on regarde également les dates
- comparer l'impact du paramètre threshold

Son fonctionnement est le suivant :
- Chargement des noms depuis un fichier CSV
- Calcul d’un exemple de distance entre deux noms
- Analyse de l'impact du threshold
- Comparaison avec ou sans date en fonction du nombre de noms
- Affichage des temps d’exécution 

Le paramètre limit est en partie utilisé pour éviter une explosion du temps de calcul lors des tests.

## Serveur TCP concurrent

> Utilisation (dans le dossier Projet_go) `go run cmd/server/main.go [--addr :PORT]`
>> Note : Il faut impérativement lancer le serveur avant le client.


9.1 Fonctionnement général

Le serveur TCP écoute sur un port donné (par défaut :8080) et accepte plusieurs connexions simultanément.
Chaque client est traité dans une goroutine dédiée, ce qui garantit que le serveur reste réactif.

9.2 Protocole TCP

Le protocole est volontairement simple :
- Le client envoie une ligne d’en-tête terminée par \n, contenant :
threshold=<int> limit=<int> csvbytes=<int>
- Le client envoie ensuite exactement csvbytes octets correspondant au contenu brut du CSV.
- Le serveur traite la requête et renvoie les résultats, puis ferme la connexion.

9.3 Robustesse réseau

- La lecture du CSV est effectuée avec io.ReadFull, garantissant la lecture complète des données annoncées.
- Une taille maximale de CSV (50 MB) est imposée pour éviter les abus.

## Client TCP

>  Utilisation (dans le dossier Projet_go) `go run cmd/client/main.go [--addr ADDR:PORT] [--csv PATH] [--threshold N] [--limit N] [--usedate 0|1]`
>> Notes : 
>>- ne pas oublier de lancer le serveur avant le (les) clients. 
>>- Attention au chemin vers le csv. 

Le client TCP permet de tester facilement le serveur :
- lecture d’un fichier CSV local
- connexion au serveur via TCP
- envoi de l’en-tête et du CSV brut
- lecture et affichage de la réponse du serveur

Il constitue un outil simple pour valider le bon fonctionnement du serveur et mesurer les performances côté serveur.

## Évaluation des performances
### performances temporelles

Les performances sont évaluées en comparant :
- une implémentation séquentielle
- une implémentation concurrente

Limit correspond à la taille du tableau de noms; Threshold correspond à la distance maximale entre deux noms; sequential et concurrent correspondent aux temps d'exécution des fonctions 

| limit | threshold | sequential (ms) | concurrent (ms) |
|-------|-----------|----------------|----------------|
| 200   | 2         | 31             | 14             |
| 200   | 3         | 41             | 17             |
| 500   | 2         | 198            | 99             |
| 500   | 3         | 269            | 122            |
| 1000  | 2         | 809            | 351            |
| 1000  | 3         | 1090           | 447            |
| 5000  | 2         | 25506          | 7728           |
| 5000  | 3         | 59705          | 8907           |


Les résultats montrent que :
- pour de petits volumes, l’approche séquentielle est suffisante
- lorsque le volume augmente, la version concurrente permet de réduire le temps d’exécution en exploitant plusieurs cœurs CPU
- le coût intrinsèque du problème (O(n²)) reste le facteur limitant principal

Nous cherchons aussi à déterminer l'impact du nombre de go routines sur le temps d'exécution du programme concurrentiel. Pour cela, nous comparons les durées en fonction du nombre de go routines (paramètre `workers`). Nous mesurons ces valeurs pour un `threshold` de 2 et une `limit` de 5000 (pour éviter une attente trop longue). 

| Nombre de workers | durée (ms) |
| ----------------- | ---------- |
| 1                 | 26494      |
| 3                 | 11550      |
| 6                 | 9053       |
| 12                | 8158       |
| 24                | 8184       |
| 48                | 8232       |

Le minimum est obtenu ici pour 12 go routines, même si la durée pour 24 et 48 est sensiblement identique. 12 correspond ici au nombre de CPU de la machine sur laquelle a été exécuté le programme. Comme on peut s'y attendre, il s'agit du nombre minimal de workers pour la résolution la plus rapide possible de ce problème. 

Nous remarquons également que la durée pour un worker est assez proche (à une seconde près) de la durée d'exécution du programme en séquentiel, ce qui est encore une fois attendu.

### impact du paramètre threshold
Nous comparons ici le nombre de matches pour différents thresholds, avec une limit de 5000 (pour ne pas que le programme mette trop de temps à s'exécuter) et useDate à faux (on ne compare pas en plus les dates)

| threshold | nbMatches |
| --------- | --------- |
| 1         | 1         |
| 2         | 5         |
| 3         | 20        |
| 4         | 47        |
| 5         | 107       |

Les résultats montrent l'impact du paramètre threshold sur le nombre de matches trouvés (le nombre est multiplié par 5 au maximum, et au moins par 2 quand on augmente le threshold de 1). Le choix du threshold est donc crucial pour trouver les bons matches. 

### impact de l'utilisation des dates
Nous comparons ici :
- une implémentation en comparant uniquement les noms
- une implémentation qui compare également les dates

Nous utilisons par défaut un threshold de 2; useDate est à true si les dates ont été comparées, à false sinon. 

| limit  | useDate | nbMatches |
|--------|---------|-----------|
| 1000   | false   | 1         |
| 1000   | true    | 1         |
| 5000   | false   | 5         |
| 5000   | true    | 1         |
| 10000  | false   | 50        |
| 10000  | true    | 8         |
| 13395  | false   | 126       |
| 13395  | true    | 15        |

Les résultats montrent qu'il existe de nombreuses lignes dont les noms matchent, mais pas les dates. Ces lignes en moins peuvent correspondre à des "faux positifs" : les personnes sont simplement des quasi homonymes.
Cependant, il pourrait aussi y avoir des erreurs sur les dates (non pris en compte ici). 

## Limites du projet

- La complexité quadratique limite la scalabilité sur de très grands jeux de données
- Le CSV est chargé entièrement en mémoire
- La distance de Levenshtein est calculée intégralement pour chaque paire, sans optimisation par seuil
- Il faut bien choisir les valeurs des paramètres (threshold en particulier)
- Nous avons regardé des matches parfaits pour les dates, mais il pourrait très bien y avoir des erreurs à la fois dans les noms et dans les dates, ou uniquement dans les dates. Celles-ci ne sont pas prises en compte par notre algorithme
- Nous aurions pu également comparer les sexes pour limiter les "faux positifs" (cas des prénoms unisexe)
- Nos algorithmes retournent les noms doublons. Pour une utilisation réelle avec des clients, il aurait probablement fallu retourner le nouveau fichier csv sans les doublons. 

Ces limites sont acceptables dans le cadre pédagogique du projet.


## Conclusion

Ce projet a permis de mettre en œuvre de manière concrète les concepts fondamentaux de concurrence en Go et de programmation réseau.
Il illustre comment un problème réel de qualité des données peut être traité à l’aide d’algorithmes de matching approximatif, tout en soulignant les enjeux de performance associés, et l'importance du bon choix des paramètres
