//Algorithme des distance de Levenshtein

package levenshtein //Déclare que ce fichier appartient au package Go nommé levenshtein

// Distance calcule la distance de Levenshtein entre deux chaînes.
// Version plus lisible et robuste (Unicode-safe via []rune).
func Distance(a, b string) int { //fonction avec en entree: 2 chaines a et b et en sortie un int = ;a distance de levenshtein
	ra := []rune(a) //Convertit la chaîne a en tableau de runes (car par exemple: "你" = 3 octets en UTF-8 et "你" = 1 rune)
	rb := []rune(b)

	// Si l'une des chaînes est vide, la distance = longueur de l'autre.
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	// prev = ligne précédente de la programmation dynamique, curr = ligne courante
	// optimisation en mémoire : on ne garde que la ligne précédente et la ligne courante
	prev := make([]int, len(rb)+1) //Crée un tableau d’entiers de taille len(rb)+1 ; Représente les distances entre un préfixe de a et tous les préfixes de b (y compris le préfixe vide)
	curr := make([]int, len(rb)+1)

	// Initialisation : transformer chaîne vide "" en rb[:j] nécessite j insertions
	for j := 0; j <= len(rb); j++ {
		prev[j] = j //Remplit la ligne 0 : distance("", rb[:j]) = j
	}

	for i := 1; i <= len(ra); i++ {
		// transformer ra[:i] en "" nécessite i suppressions
		curr[0] = i

		for j := 1; j <= len(rb); j++ {
			cost := 0               //Coût de substitution : par défaut 0 (si caractères identiques)
			if ra[i-1] != rb[j-1] { //Compare le dernier caractère du préfixe ra[:i] et du préfixe rb[:j]
				cost = 1 //Si les caractères diffèrent, une substitution coûte 1
			}

			//Calcul des trois opérations possibles
			del := prev[j] + 1      // suppression
			ins := curr[j-1] + 1    // insertion
			sub := prev[j-1] + cost // substitution

			curr[j] = min3(del, ins, sub) //La distance DP pour (i, j) est le minimum des 3 opérations
		}

		// La ligne courante devient la ligne précédente
		prev, curr = curr, prev
	}

	return prev[len(rb)]
}

func min3(a, b, c int) int { //fonction qui renvoie le minimum de 3 entiers
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
