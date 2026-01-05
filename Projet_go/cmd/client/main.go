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
	addr := flag.String("addr", "127.0.0.1:8080", "Adresse serveur (ex: 127.0.0.1:8080)")
	csvPath := flag.String("csv", "data/UniversoGITT_Medellin.csv", "Chemin du CSV à envoyer")
	threshold := flag.Int("threshold", 2, "Seuil Levenshtein")
	limit := flag.Int("limit", 500, "Limit (0 = tout, attention O(n^2))")
	flag.Parse()

	csvBytes, err := os.ReadFile(*csvPath)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTimeout("tcp", *addr, 5*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 1) En-tête + taille CSV
	header := fmt.Sprintf("threshold=%d limit=%d csvbytes=%d\n", *threshold, *limit, len(csvBytes))
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
