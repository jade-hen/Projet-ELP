// “Adaptateur” CSV → Go (pour les données)

package data

import (
	"encoding/csv"
	"fmt"
	"os" // pour ouvrir un fichier
)

func LoadFirstColumn(path string) ([]string, error) { //Entrée : path = chemin du fichier CSV; Sortie : []string = liste de valeurs (la 1ère colonne) et error = erreur si quelque chose se passe mal
	// Ouvre le fichier CSV en lecture
	f, err := os.Open(path)
	if err != nil {
		return nil, err // Si le fichier n'existe pas ou n'est pas accessible
	}
	defer f.Close() // Ferme le fichier à la fin de la fonction, quoi qu'il arrive

	// Lit tout le contenu du CSV en mémoire :
	// records = tableau de lignes, chaque ligne = tableau de colonnes ([][]string)
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err // CSV invalide, problème d'encodage, etc.
	}

	// Si aucune ligne n'est présente, le CSV est vide (même pas d'en-tête)
	if len(records) == 0 {
		return nil, fmt.Errorf("empty csv")
	}

	// Slice de sortie : une liste de strings (valeurs de la première colonne)
	var out []string

	// Parcourt toutes les lignes de données (on saute la ligne 0 : l'en-tête)
	for _, row := range records[1:] {
		// Ignore les lignes vides ou les lignes sans première colonne exploitable
		if len(row) == 0 || row[0] == "" {
			continue
		}

		// Ajoute la valeur de la première colonne à la sortie
		out = append(out, row[0])
	}

	// Renvoie la liste des valeurs (et nil pour dire "pas d'erreur")
	return out, nil
}
