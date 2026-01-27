package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {

	// créer des options pour que le client puisse choisir ses paramètres
	addr := flag.String("addr", "127.0.0.1:8080", "Adresse serveur (ex: 127.0.0.1:8080)")
	csvPath := flag.String("csv", "data/data.csv", "Chemin du CSV à envoyer")
	threshold := flag.Int("threshold", 2, "Seuil Levenshtein")
	limit := flag.Int("limit", 500, "Limit (0 = tout, attention O(n^2))")
	useDate := flag.Int("usedate", 1, "Comparer les dates ou non (0/1)")

	flag.Parse()

	// Lecture du fichier CSV en mémoire
	csvBytes, err := os.ReadFile(*csvPath)
	if err != nil {
		panic(err)
	}

	// Connexion TCP au serveur (avec timeout)
	conn, err := net.DialTimeout("tcp", *addr, 5*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 1) En-tête + taille CSV
	header := fmt.Sprintf("threshold=%d limit=%d useDate=%d csvbytes=%d\n", *threshold, *limit, *useDate, len(csvBytes))
	if _, err := conn.Write([]byte(header)); err != nil {
		panic(err)
	}

	// 2) CSV brut
	if _, err := conn.Write(csvBytes); err != nil {
		panic(err)
	}

	// 3) Lire la réponse (le serveur ferme ensuite la connexion)
	resp, err := io.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(resp))
}
