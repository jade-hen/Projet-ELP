# Concurrent Name Matching Server (Go)

## Description

Ce projet est une preuve de concept (PoC) développée en langage Go, dont l’objectif est de mettre en pratique les concepts de concurrence et de programmation réseau à travers un cas d’usage concret : la détection de doublons et de correspondances approximatives entre des noms issus de fichiers CSV.

L’application repose sur l’algorithme de distance de Levenshtein, permettant de mesurer la similarité entre deux chaînes de caractères, et explore deux approches :

une approche séquentielle

une approche concurrente, exploitant les goroutines et les channels

Le projet inclut également un serveur TCP concurrent, capable de traiter plusieurs clients simultanément.

## Contexte et motivation

Ce projet s’inspire d’un problème réel de qualité des données et d’identification des personnes, notamment mis en évidence dans le contexte des « faux positifs » en Colombie.

Dans ce scandale, des erreurs d’identification ont conduit à des confusions entre individus, souvent à cause de variations ou d’incohérences dans les noms (différences d’orthographe, accents, abréviations, erreurs de saisie ou sources multiples). Ces situations illustrent les limites des comparaisons strictes basées uniquement sur l’égalité exacte des chaînes de caractères.

L’objectif de ce projet est de montrer comment des techniques de rapprochement approximatif, comme la distance de Levenshtein, peuvent être utilisées pour détecter des doublons ou des correspondances potentielles dans une base de données de noms, tout en tenant compte des variations textuelles.

Le projet met également en évidence les enjeux de performance liés à ce type de traitement sur de grands volumes de données, et l’intérêt de la concurrence pour rendre ces analyses exploitables à l’échelle.

---

## Objectifs du projet

- Prendre en main le langage Go
- Implémenter un algorithme de matching approximatif (distance de Levenshtein)
- Exploiter la concurrence pour améliorer les performances
- Mettre en place un serveur TCP concurrent
- Comparer une approche séquentielle et une approche concurrente

---

## Fonctionnalités principales

- Calcul de la distance de Levenshtein entre deux chaînes de caractères
- Détection de doublons et de correspondances approximatives avec un seuil configurable
- Traitement concurrent des comparaisons grâce aux goroutines
- Utilisation de channels pour la distribution des tâches
- Synchronisation des goroutines via des wait groups
- Serveur TCP capable de gérer plusieurs clients simultanément
- Mesure et analyse des performances

---

## Architecture générale

L’application est structurée autour des composants suivants :

- **Algorithme** : calcul de la distance de Levenshtein
- **Concurrence** : worker pool basé sur des goroutines et des channels
- **Réseau** : serveur TCP concurrent (une goroutine par client)
- **Orchestration** : intégration du calcul concurrent dans le serveur

---

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

---

## Équipe

Ce projet est réalisé en groupe par :

- **Anna Grataloup**
- **Emma Payrard**
- **Jade Henninot**

---

## Contexte académique

Ce projet est réalisé dans un cadre pédagogique afin de mettre en pratique les concepts de concurrence et de programmation réseau en Go à travers une application concrète, mesurable et ancrée dans un problème réel de qualité des données.

---
