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
	names, err := data.LoadFirstColumn(filePath)
	if err != nil { // En cas d'erreur (fichier absent, CSV invalide, etc.), on arrête le programme et on affiche l'erreur.
		log.Fatal(err)
	}

	fmt.Println("Nb de noms chargés:", len(names)) // Affiche le nombre de noms chargés depuis le CSV

	if len(names) >= 2 { // Vérifie qu'on a au moins deux noms pour faire une comparaison
		// Calcule et affiche la distance de Levenshtein entre les deux premiers noms; exemple d’usage rapide pour valider que l’algo fonctionne sur les données.
		fmt.Println("Exemple distance:", names[0], "vs", names[1], "=>",
			levenshtein.Distance(names[0], names[1]))
	}

	//Comparaison des durées
	limit := 1000 //limiter le volume de données

	startSeq := time.Now()
	//Parcours pour trouver les matches en séquentiel
	matcher.FindMatchesSequential(names, 5, limit)
	elapsedSeq := time.Since(startSeq)

	startConc := time.Now()
	//Parcours pour trouver les matches en concurrent
	matcher.FindMatchesConcurrent(names, 5, limit, 0)
	elapsedConc := time.Since(startConc)

	fmt.Println("Temps d'exécution en séquentiel :", elapsedSeq)
	fmt.Println("Temps d'exécution en séquentiel :", elapsedConc)
}
