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
	fmt.Println()

	// ---------------- Table 1 : nb de matches pour limit=5000 et threshold 1..5 ----------------
	limit := 5000
	thresholds := []int{1, 2, 3, 4, 5}

	fmt.Println("=== Performances : nombre de matches pour limit=5000 et threshold 1..5 ===")
	fmt.Println("threshold\tnbMatches")
	for _, th := range thresholds {
		matches := matcher.FindMatchesConcurrent(persons, th, limit, 0, false)
		fmt.Printf("%d\t\t%d\n", th, len(matches))
	}
	fmt.Println()

	// ---------------- Table 2 : nb de matches pour plusieurs limits, avec et sans date ----------------
	limits := []int{200, 500, 1000, 5000, 13395}
	threshold := 2

	fmt.Println("=== Performances : nombre de matches selon limit et utilisation des dates ===")
	fmt.Println("limit\tuseDate\tnbMatches")
	for _, lim := range limits {
		// Sans date
		matchesNoDate := matcher.FindMatchesConcurrent(persons, threshold, lim, 0, false)
		fmt.Printf("%d\t%v\t%d\n", lim, false, len(matchesNoDate))

		// Avec date
		matchesWithDate := matcher.FindMatchesConcurrent(persons, threshold, lim, 0, true)
		fmt.Printf("%d\t%v\t%d\n", lim, true, len(matchesWithDate))
	}
	fmt.Println()

	// ---------------- Table 3 : comparaison des performances temporelles ----------------
	fmt.Println("=== Comparaison performance temporelle (ms) ===")
	fmt.Println("limit\tthreshold\tsequential(ms)\tconcurrent(ms)")
	for _, lim := range limits {
		for _, th := range thresholds {
			// Séquentiel
			startSeq := time.Now()
			_ = matcher.FindMatchesSequential(persons, th, lim, true)
			elapsedSeq := time.Since(startSeq)

			// Concurrent
			startConc := time.Now()
			_ = matcher.FindMatchesConcurrent(persons, th, lim, 0, true)
			elapsedConc := time.Since(startConc)

			fmt.Printf("%d\t%d\t%d\t%d\n",
				lim, th,
				elapsedSeq.Milliseconds(),
				elapsedConc.Milliseconds(),
			)
		}
	}
}
