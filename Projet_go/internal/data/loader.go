package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// Person représente une ligne du CSV avec nom et date
type Person struct {
	Name string
	Date string
}

// Charge le CSV et renvoie []Person (nom + date)
func LoadNamesAndDates(path string) ([]Person, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadNamesAndDatesFromReader(f)
}

// LoadFirstColumnFromReader lit un CSV depuis un flux (réseau, mémoire, fichier)
// et renvoie la 1ère colonne comme []string.
// skipHeader=true => ignore la première ligne (souvent un header).
func LoadNamesAndDatesFromReader(r io.Reader) ([]Person, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 1 {
		return nil, fmt.Errorf("CSV vide ou sans données")
	}

	out := make([]Person, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 3 {
			continue
		}
		name := strings.TrimSpace(row[0])
		date := strings.TrimSpace(row[2])
		if name == "" || date == "" {
			continue
		}
		out = append(out, Person{Name: name, Date: date})
	}
	return out, nil
}
