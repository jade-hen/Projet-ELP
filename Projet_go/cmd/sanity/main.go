package main

import (
	"fmt" // Affichage standard sur la sortie console
	"log" // Gestion simple des erreurs fatales (log + exit)
	"os"
	"time" // Chronométrage des fonctions pour comparaison des méthodes

	"levenshtein/internal/data"
	"levenshtein/internal/levenshtein"
	"levenshtein/internal/matcher"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : go run main.go <chemin_du_fichier>")
		return
	}

	filePath := os.Args[1]
	fmt.Println("Fichier utilisé :", filePath)
	// Charge la première colonne du fichier CSV (en sautant l’en-tête)
	// Renvoie une liste de chaînes (ex: des noms) et une erreur si échec
	persons, err := data.LoadNamesAndDates(filePath)
	if err != nil { // En cas d'erreur (fichier absent, CSV invalide, etc.), on arrête le programme et on affiche l'erreur.
		log.Fatal(err)
	}

	fmt.Println("Nb de noms chargés:", len(persons)) // Affiche le nombre de noms chargés depuis le CSV

	if len(persons) >= 2 { // Vérifie qu'on a au moins deux noms pour faire une comparaison
		// Calcule et affiche la distance de Levenshtein entre les deux premiers noms; exemple d’usage rapide pour valider que l’algo fonctionne sur les données.
		fmt.Println("Exemple distance:", persons[0], "vs", persons[1], "=>",
			levenshtein.Distance(persons[0].Name, persons[1].Name))
	}

	limit := 13395  //limiter le volume de données
	useDate := true // ou false si on veut ignorer les dates
	threshold := 2  // choisir le nombre de différences entre deux noms

	//Comparaison des résulats en utilisant ou non les dates
	fmt.Println("=== Matches sans utiliser les dates ===")
	matchesNoDate := matcher.FindMatchesConcurrent(persons, threshold, limit, 0, false)
	fmt.Println("Nb de matches (sans dates) :", len(matchesNoDate))

	fmt.Println("=== Matches en utilisant les dates ===")
	matchesWithDate := matcher.FindMatchesConcurrent(persons, threshold, limit, 0, true)
	fmt.Println("Nb de matches (avec dates) :", len(matchesWithDate))

	//Comparaison des résultats en séquentiel ou concurrence
	startSeq := time.Now()
	//Parcours pour trouver les matches en séquentiel
	matcher.FindMatchesSequential(persons, threshold, limit, useDate)
	elapsedSeq := time.Since(startSeq)

	startConc := time.Now()
	//Parcours pour trouver les matches en concurrent
	matcher.FindMatchesConcurrent(persons, threshold, limit, 0, useDate)
	elapsedConc := time.Since(startConc)

	fmt.Println("Temps d'exécution en séquentiel :", elapsedSeq)
	fmt.Println("Temps d'exécution en concurrence :", elapsedConc)
}
