package matcher

import (
	"fmt"
	"runtime" // Pour connaître le nombre de CPU (pour choisir un nb de workers)
	"sort"
	"sync" // WaitGroup pour synchroniser les goroutines

	"levenshtein/internal/data"
	"levenshtein/internal/levenshtein" // Algo de distance de Levenshtein
)

// Match représente une paire (A, B) + sa distance.
type Match struct {
	A        string
	B        string
	Distance int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func subsetPersons(persons []data.Person, limit int) []data.Person {
	if limit <= 0 || limit >= len(persons) {
		return persons
	}
	return persons[:limit]
}

// FindMatchesSequential compare tous les couples (i, j) de manière séquentielle
// - threshold : distance max pour considérer une paire comme match
// - limit     : limite le nombre de noms considérés (utile pour tester rapidement)
func FindMatchesSequential(persons []data.Person, threshold, limit int, useDate bool) []Match {
	//on limite le volume de données
	persons = subsetPersons(persons, limit)

	// Slice des résultats (matches trouvés)
	matches := make([]Match, 0)

	// Compare chaque nom avec tous ceux qui le suivent (évite les doublons)
	for i := 0; i < len(persons); i++ {
		for j := i + 1; j < len(persons); j++ {

			// Optimisation : si la différence de longueur dépasse le seuil, on ne calcule pas
			if abs(len(persons[i].Name)-len(persons[j].Name)) > threshold {
				continue
			}

			// Calcule la distance entre les deux chaînes
			d := levenshtein.Distance(persons[i].Name, persons[j].Name)

			// Si la distance est suffisamment faible, on garde la paire
			if d <= threshold {
				if !useDate || persons[i].Date == persons[j].Date {
					matches = append(matches, Match{
						A:        persons[i].Name,
						B:        persons[j].Name,
						Distance: d,
					})
				}
			}
		}
	}

	// Trie les matches pour avoir un ordre stable et exploitable
	sortMatches(matches)
	return matches
}

// FindMatchesConcurrent fait la même chose, mais en parallèle
// - workers : nb de goroutines de calcul ; si <= 0, on prend runtime.NumCPU()
func FindMatchesConcurrent(persons []data.Person, threshold, limit, workers int, useDate bool) []Match {
	// Optionnel : on limite le volume de données
	persons = subsetPersons(persons, limit)
	// Si aucun nombre de workers n'est fourni, on utilise le nb de coeurs CPU
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// job = une paire d'indices (i, j) à comparer
	type job struct{ i, j int }

	// jobs : file de travail (paires à comparer)
	// results : file des matches (paires qui passent le threshold)
	jobs := make(chan job, 10000)
	results := make(chan Match, 10000)

	// WaitGroup pour attendre que tous les workers aient fini
	var wg sync.WaitGroup
	wg.Add(workers)

	// Démarre "workers" goroutines qui lisent jobs et produisent éventuellement des results
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for jb := range jobs {
				a := persons[jb.i]
				b := persons[jb.j]

				// Optimisation : si la différence de longueur dépasse le seuil, on ne calcule pas
				if abs(len(a.Name)-len(b.Name)) > threshold {
					continue
				}

				d := levenshtein.Distance(a.Name, b.Name)
				// Si match, on l’envoie dans results
				if d <= threshold {
					if !useDate || a.Date == b.Date {
						results <- Match{A: a.Name, B: b.Name, Distance: d}
					}
				}
			}
		}()
	}

	// Producteur de jobs : génère toutes les paires (i, j), puis ferme jobs
	go func() {
		for i := 0; i < len(persons); i++ {
			if i%1000 == 0 {
				fmt.Println("progress:", i, "/", len(persons))
			}
			for j := i + 1; j < len(persons); j++ {
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
