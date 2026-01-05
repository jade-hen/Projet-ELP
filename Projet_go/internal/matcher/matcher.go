package matcher

import (
	"runtime" // Pour connaître le nombre de CPU (pour choisir un nb de workers)
	"sort"
	"sync" // WaitGroup pour synchroniser les goroutines

	"levenshtein/internal/levenshtein" // Algo de distance de Levenshtein
)

// Match représente une paire (A, B) + sa distance.
type Match struct {
	A        string
	B        string
	Distance int
}

// subset limite la liste à "limit" éléments si limit est > 0
// - limit <= 0 : pas de limite
// - limit >= len(names) : pas de limite (on garde tout)
func subset(names []string, limit int) []string {
	if limit <= 0 || limit >= len(names) {
		return names
	}
	return names[:limit]
}

// FindMatchesSequential compare tous les couples (i, j) de manière séquentielle
// - threshold : distance max pour considérer une paire comme match
// - limit     : limite le nombre de noms considérés (utile pour tester rapidement)
func FindMatchesSequential(names []string, threshold, limit int) []Match {
	//on limite le volume de données
	names = subset(names, limit)

	// Slice des résultats (matches trouvés)
	matches := make([]Match, 0)

	// Compare chaque nom avec tous ceux qui le suivent (évite les doublons)
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			// Calcule la distance entre les deux chaînes
			d := levenshtein.Distance(names[i], names[j])

			// Si la distance est suffisamment faible, on garde la paire
			if d <= threshold {
				matches = append(matches, Match{A: names[i], B: names[j], Distance: d})
			}
		}
	}

	// Trie les matches pour avoir un ordre stable et exploitable
	sortMatches(matches)
	return matches
}

// FindMatchesConcurrent fait la même chose, mais en parallèle
// - workers : nb de goroutines de calcul ; si <= 0, on prend runtime.NumCPU()
func FindMatchesConcurrent(names []string, threshold, limit, workers int) []Match {
	// Optionnel : on limite le volume de données
	names = subset(names, limit)
	// Si aucun nombre de workers n'est fourni, on utilise le nb de coeurs CPU
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// job = une paire d'indices (i, j) à comparer
	type job struct{ i, j int }

	// jobs : file de travail (paires à comparer)
	// results : file des matches (paires qui passent le threshold)
	// Buffers 1024 pour réduire les blocages entre producteurs/consommateurs
	jobs := make(chan job, 1024)
	results := make(chan Match, 1024)

	// WaitGroup pour attendre que tous les workers aient fini
	var wg sync.WaitGroup
	wg.Add(workers)

	// Démarre "workers" goroutines qui lisent jobs et produisent éventuellement des results
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for jb := range jobs {
				a := names[jb.i]
				b := names[jb.j]
				// Calcul de distance
				d := levenshtein.Distance(a, b)
				// Si match, on l’envoie dans results
				if d <= threshold {
					results <- Match{A: a, B: b, Distance: d}
				}
			}
		}()
	}

	// Producteur de jobs : génère toutes les paires (i, j), puis ferme jobs
	go func() {
		for i := 0; i < len(names); i++ {
			for j := i + 1; j < len(names); j++ {
				jobs <- job{i: i, j: j}
			}
		}
		close(jobs)
	}()

	// Ferme results quand tous les workers ont fini (wg.Wait)
	go func() {
		wg.Wait()
		close(results)
	}()

	// Consomme results pour construire la slice finale
	matches := make([]Match, 0)
	for m := range results {
		matches = append(matches, m)
	}

	sortMatches(matches)
	return matches
}

// sortMatches trie les matches selon :
// 1) Distance croissante (meilleurs matches en premier)
// 2) A alphabétique
// 3) B alphabétique
func sortMatches(matches []Match) {
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Distance != matches[j].Distance {
			return matches[i].Distance < matches[j].Distance
		}
		if matches[i].A != matches[j].A {
			return matches[i].A < matches[j].A
		}
		return matches[i].B < matches[j].B
	})
}
