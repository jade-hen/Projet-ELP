package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"levenshtein/internal/data"
	"levenshtein/internal/matcher"
)

type request struct {
	threshold int
	limit     int
	csvBytes  int
	useDate   bool
}

func main() {
	addr := flag.String("addr", ":8080", "Adresse d'écoute (ex: :8080)")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Serveur TCP en écoute sur", *addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn) // 1 goroutine par client => serveur TCP concurrent
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// Évite les blocages infinis
	_ = conn.SetDeadline(time.Now().Add(60 * time.Second))

	r := bufio.NewReader(conn)

	// 1) Lire l'en-tête (une ligne)
	headerBytes, err := r.ReadBytes('\n')
	if err != nil {
		fmt.Fprintln(conn, "ERR header read:", err)
		return
	}
	header := string(headerBytes)

	req := parseHeader(strings.TrimSpace(header))
	if req.csvBytes <= 0 {
		fmt.Fprintln(conn, "ERR csvbytes must be > 0")
		return
	}

	// Protection simple (évite d'avaler des tailles absurdes)
	const maxCSV = 50 * 1024 * 1024 // 50MB
	if req.csvBytes > maxCSV {
		fmt.Fprintln(conn, "ERR csvbytes too large (max 50MB)")
		return
	}

	// 2) Lire exactement csvBytes octets (CSV brut)
	raw := make([]byte, req.csvBytes)
	if _, err := io.ReadFull(r, raw); err != nil {
		fmt.Fprintln(conn, "ERR csv read:", err)
		return
	}

	// 3) Parser CSV -> []string
	persons, err := data.LoadNamesAndDatesFromReader(bytes.NewReader(raw))
	if err != nil {
		fmt.Fprintln(conn, "ERR csv parse:", err)
		return
	}
	if len(persons) < 2 {
		fmt.Fprintln(conn, "ERR not enough persons (need >=2)")
		return
	}
	fmt.Println(bool(req.useDate))
	// 4) Matching concurrent (workers = NumCPU)
	workers := runtime.NumCPU()
	start := time.Now()
	matches := matcher.FindMatchesConcurrent(persons, req.threshold, req.limit, workers, bool(req.useDate))
	elapsed := time.Since(start)

	// 5) Réponse
	fmt.Fprintf(conn, "OK persons=%d threshold=%d limit=%d workers=%d\n", len(persons), req.threshold, req.limit, workers)
	fmt.Fprintf(conn, "matches=%d elapsed_ms=%d\n", len(matches), elapsed.Milliseconds())

	maxShow := 20 //affichage de 20 max
	if maxShow > len(matches) {
		maxShow = len(matches)
	}
	for i := 0; i < maxShow; i++ {
		m := matches[i]
		fmt.Fprintf(conn, "d=%d | %s <-> %s\n", m.Distance, m.A, m.B)
	}
	conn.Close()
}

func parseHeader(line string) request {
	out := request{
		threshold: 2,
		limit:     500,
		csvBytes:  0,
		useDate:   false,
	}

	for _, p := range strings.Fields(line) {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])

		switch key {
		case "threshold":
			if x, err := strconv.Atoi(val); err == nil {
				out.threshold = x
			}
		case "limit":
			if x, err := strconv.Atoi(val); err == nil {
				out.limit = x
			}
		case "csvbytes":
			if x, err := strconv.Atoi(val); err == nil {
				out.csvBytes = x
			}
		case "usedate":
			if x, err := strconv.Atoi(val); err == nil && x != 0 {
				out.useDate = true
			}
		}
	}
	return out
}
