# Concurrent Name Matching Server (Go)

## Description

Ce projet est une preuve de concept (PoC) développée en langage Go, dont l’objectif est de mettre en pratique les concepts de concurrence et de programmation réseau à travers un cas d’usage concret : la détection de doublons et de correspondances approximatives entre des noms issus de fichiers CSV.

L’application repose sur l’algorithme de distance de Levenshtein, permettant de mesurer la similarité entre deux chaînes de caractères, et explore deux approches :
- une approche séquentielle
- une approche concurrente, exploitant les goroutines et les channels

Le projet inclut également un serveur TCP concurrent, capable de traiter plusieurs clients simultanément.

## Contexte et motivation

Ce projet s’inspire d’un problème réel de qualité des données et d’identification des personnes, notamment mis en évidence dans le contexte des « faux positifs » en Colombie.

Dans ce scandale, des erreurs d’identification ont conduit à des confusions entre individus, souvent à cause de variations ou d’incohérences dans les noms (différences d’orthographe, accents, abréviations, erreurs de saisie ou sources multiples). Ces situations illustrent les limites des comparaisons strictes basées uniquement sur l’égalité exacte des chaînes de caractères.

L’objectif de ce projet est de montrer comment des techniques de rapprochement approximatif, comme la distance de Levenshtein, peuvent être utilisées pour détecter des doublons ou des correspondances potentielles dans une base de données de noms, tout en tenant compte des variations textuelles.

Le projet met également en évidence les enjeux de performance liés à ce type de traitement sur de grands volumes de données, et l’intérêt de la concurrence pour rendre ces analyses exploitables à l’échelle.


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
- Utilisation de channels pour la distribution des tâches
- Synchronisation des goroutines via des wait groups
- Serveur TCP capable de gérer plusieurs clients simultanément
- Mesure et analyse des performances


## Architecture générale

Le projet est structuré en plusieurs modules :
- data : chargement et parsing des fichiers CSV
- levenshtein : calcul de la distance entre deux chaînes
- matcher : détection des correspondances (séquentielle et concurrente)
- sanity : programme de validation et de benchmark
- client TCP : envoi de données au serveur
- serveur TCP : réception, traitement concurrent et réponse

Cette séparation permet une meilleure lisibilité et une réutilisation des composants dans différents contextes (ligne de commande, réseau).


## Chargement des données (module data)

Le module data est responsable du chargement des noms depuis un fichier CSV ou depuis un flux générique (io.Reader).

Deux fonctions principales sont fournies :
- LoadFirstColumn(path string) : charge un CSV depuis un fichier
- LoadFirstColumnFromReader(r io.Reader, skipHeader bool) : charge un CSV depuis un flux (utilisé notamment par le serveur TCP)

Seule la première colonne du CSV est conservée, correspondant aux noms à comparer.
Les lignes vides ou valeurs vides sont ignorées, et la première ligne est sautée lorsqu’elle correspond à un en-tête.

Le CSV est chargé intégralement en mémoire à l’aide de encoding/csv.ReadAll(), ce qui simplifie l’implémentation dans le cadre de ce projet pédagogique.


## Distance de Levenshtein (module levenshtein)

La distance de Levenshtein est implémentée via une approche de programmation dynamique.

Principes : 
- Les chaînes sont converties en []rune afin d’être Unicode-safe, ce qui permet de gérer correctement les accents et caractères non ASCII.
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

Le module matcher fournit deux implémentations du matching :

7.1 Approche séquentielle
La fonction FindMatchesSequential compare chaque paire de noms (i, j) avec j > i, afin d’éviter les doublons et les auto-comparaisons.

Pour chaque paire :
- la distance de Levenshtein est calculée
- si la distance est inférieure ou égale au seuil (threshold), la paire est conservée

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

Limitation du volume : un paramètre limit permet de restreindre le nombre de noms traités, afin de contrôler la complexité O(n²) inhérente au problème. 


## Programme de validation et benchmark (sanity)

Le programme sanity permet de :
- vérifier le bon fonctionnement de l’algorithme de Levenshtein sur des données réelles
- comparer les temps d’exécution des versions séquentielle et concurrente

Son fonctionnement est le suivant :
- Chargement des noms depuis un fichier CSV
- Calcul d’un exemple de distance entre deux noms
- Exécution du matching séquentiel
- Exécution du matching concurrent
- Affichage des temps d’exécution

Un paramètre limit est utilisé pour éviter une explosion du temps de calcul lors des tests.


## Serveur TCP concurrent
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

- Un ReadDeadline est appliqué afin d’éviter les blocages en cas de client inactif.
- La lecture du CSV est effectuée avec io.ReadFull, garantissant la lecture complète des données annoncées.
- Une taille maximale de CSV (50 MB) est imposée pour éviter les abus.

## Client TCP

Le client TCP permet de tester facilement le serveur :
- lecture d’un fichier CSV local
- connexion au serveur via TCP
- envoi de l’en-tête et du CSV brut
- lecture et affichage de la réponse du serveur

Il constitue un outil simple pour valider le bon fonctionnement du serveur et mesurer les performances côté serveur.

## Évaluation des performances

Les performances sont évaluées en comparant :
- une implémentation séquentielle
- une implémentation concurrente

Les résultats montrent que :
- pour de petits volumes, l’approche séquentielle est suffisante
- lorsque le volume augmente, la version concurrente permet de réduire le temps d’exécution en exploitant plusieurs cœurs CPU
- le coût intrinsèque du problème (O(n²)) reste le facteur limitant principal

## Limites du projet

- La complexité quadratique limite la scalabilité sur de très grands jeux de données
- Le CSV est chargé entièrement en mémoire
- La distance de Levenshtein est calculée intégralement pour chaque paire, sans optimisation par seuil

Ces limites sont acceptables dans le cadre pédagogique du projet.


## Conclusion

Ce projet a permis de mettre en œuvre de manière concrète les concepts fondamentaux de concurrence en Go et de programmation réseau.
Il illustre comment un problème réel de qualité des données peut être traité à l’aide d’algorithmes de matching approximatif, tout en soulignant les enjeux de performance associés.

L’utilisation combinée des goroutines, des channels et des WaitGroups a permis de construire une application concurrente claire, modulaire et fonctionnelle, adaptée à un contexte d’apprentissage.

## Équipe

Ce projet est réalisé en groupe par :

- **Anna Grataloup**
- **Emma Payrard**
- **Jade Henninot**



## Contexte académique

Ce projet est réalisé dans un cadre pédagogique afin de mettre en pratique les concepts de concurrence et de programmation réseau en Go à travers une application concrète, mesurable et ancrée dans un problème réel de qualité des données.


FIN












## Concurrence

Le projet met en œuvre deux niveaux de concurrence :

1. **Concurrence de calcul (CPU-bound)**  
   Les comparaisons de noms sont réparties entre plusieurs goroutines à l’aide d’un worker pool afin d’exploiter les cœurs CPU disponibles.

2. **Concurrence réseau (I/O-bound)**  
   Chaque client TCP est géré dans une goroutine dédiée, ce qui permet au serveur de rester réactif même en présence de plusieurs connexions simultanées.

---

## Protocole TCP (simplifié)

Les clients se connectent au serveur via TCP et envoient une requête contenant :
- une liste de noms à comparer
- un seuil de distance pour le matching

Le serveur retourne les correspondances trouvées ainsi que la distance associée.

---

## Évaluation des performances

Le projet compare les performances :
- d’une implémentation séquentielle du matching
- d’une implémentation concurrente

Les temps d’exécution sont mesurés afin d’analyser les gains apportés par la concurrence et d’en discuter les limites.

---

## Technologies utilisées

- Go
- Goroutines
- Channels
- sync.WaitGroup
- Sockets TCP



## Équipe

Ce projet est réalisé en groupe par :

- **Anna Grataloup**
- **Emma Payrard**
- **Jade Henninot**



## Contexte académique

Ce projet est réalisé dans un cadre pédagogique afin de mettre en pratique les concepts de concurrence et de programmation réseau en Go à travers une application concrète, mesurable et ancrée dans un problème réel de qualité des données.


