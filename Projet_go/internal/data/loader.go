package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadFirstColumn lit un CSV depuis un fichier et renvoie la 1ère colonne.
// Ici on saute l'en-tête (ligne 1) par défaut.
func LoadFirstColumn(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadFirstColumnFromReader(f, true)
}

// LoadFirstColumnFromReader lit un CSV depuis un flux (réseau, mémoire, fichier)
// et renvoie la 1ère colonne comme []string.
// skipHeader=true => ignore la première ligne (souvent un header).
func LoadFirstColumnFromReader(r io.Reader, skipHeader bool) ([]string, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty csv")
	}

	start := 0
	if skipHeader && len(records) > 0 {
		start = 1
	}

	out := make([]string, 0, len(records)-start)
	for i := start; i < len(records); i++ {
		row := records[i]
		if len(row) == 0 {
			continue
		}
		val := strings.TrimSpace(row[0])
		if val == "" {
			continue
		}
		out = append(out, val)
	}
	return out, nil
}
