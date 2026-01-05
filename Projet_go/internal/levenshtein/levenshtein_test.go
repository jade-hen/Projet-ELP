//Test de l'algorithme des distance de Levenshtein

package levenshtein

import "testing" //go test détecte automatiquement les fonctions qui commencent par Test et les execute ensuite

func TestDistance(t *testing.T) { // t *testing.T = l’objet qui te permet de signaler des erreurs
	//Chaque cas de test est une structure avec a = premiere chaine, b = 2eme chaine et want = le resultat attendu
	tests := []struct {
		a, b string
		want int
	}{
		// Les cas de test
		{"", "", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"kitten", "sitting", 3},
		{"medellin", "medelin", 1},
		{"Universo", "Universo", 0},
		{"Bogota", "Bogotá", 1}, // test accent (Unicode)
		{"niño", "nino", 1},     // test ñ (Unicode)
	}

	// La boucle sur les tests
	for _, tt := range tests { //parcourt chaque cas: "_" ignore l’index (on n’en a pas besoin) et "tt" reçoit la structure du test courant
		got := Distance(tt.a, tt.b)
		if got != tt.want {
			t.Fatalf("Distance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want) //F = “Fatal” : si ce test échoue, le test s’arrête immédiatement; %q affiche la chaîne entre guillemets et %d pour les entiers
		}
	}
}

// pour l'executer: go test ou go test -v pour voir les détails
