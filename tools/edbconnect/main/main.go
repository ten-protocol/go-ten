package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/ten-protocol/go-ten/go/enclave/storage/init/edgelessdb"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
)

func main() {
	log.Println("Retrieving Edgeless DB credentials...")
	creds, found, err := edgelessdb.LoadCredentialsFromFile()
	if err != nil {
		log.Println("Error loading credentials from file:", err)
		panic(err)
	}
	if !found {
		panic("No existing EDB credentials found.")
	}
	log.Println("Found existing EDB credentials. Creating TLS config...")
	cfg, err := edgelessdb.CreateTLSCfg(creds)
	if err != nil {
		log.Println("Error creating TLS config from credentials:", err)
		panic(err)
	}
	log.Println("TLS config created. Connecting to Edgeless DB...")
	testlog.SetupSysOut()
	db, err := edgelessdb.ConnectToEdgelessDB("obscuronode-edgelessdb", cfg, testlog.Logger())
	if err != nil {
		log.Println("Error connecting to Edgeless DB:", err)
		panic(err)
	}
	log.Println("Connected to Edgeless DB.")

	startREPL(db)

	err = db.Close()
	if err != nil {
		log.Println("Error closing Edgeless DB connection:", err)
		panic(err)
	}
}

// Starts a loop that reads user input and runs queries against the Edgeless DB until user types "exit"
func startREPL(db *sql.DB) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	log.SetOutput(rl.Stderr())
	for {
		log.Println("\nEnter a query to run against the Edgeless DB (or type 'exit' to quit):")
		query, err := rl.Readline()
		if err != nil { // Handle EOF and Interrupt errors
			if errors.Is(err, readline.ErrInterrupt) {
				if len(query) == 0 {
					break
				} else {
					continue
				}
			} else if errors.Is(err, io.EOF) {
				break
			}
			log.Println("Error reading user input:", err)
			continue
		}
		// line break for readability
		log.Println()

		// Trim the newline character and surrounding whitespace
		query = strings.TrimSpace(query)

		if query == "" {
			continue
		}

		if query == "exit" {
			break
		}

		// Determine the type of query, so we can show appropriate output
		queryType := strings.ToUpper(strings.Split(query, " ")[0])
		switch queryType {
		case "SELECT", "SHOW", "DESCRIBE", "DESC", "EXPLAIN":
			// output rows
			runQuery(db, query)
		default:
			// output number of rows affected
			runExec(db, query)
		}
	}
	log.Println("Exiting...")
}

func runQuery(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("Error fetching columns:", err)
		return
	}

	// Print column headers
	for _, colName := range cols {
		log.Printf("%s\t", colName)
	}
	log.Println()

	// Prepare a slice to hold the values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Println("Error scanning row:", err)
			return
		}

		// Print the row values
		for _, val := range values {
			// Handle NULL values and convert byte slices to strings
			switch v := val.(type) {
			case nil:
				log.Print("NULL\t")
			case []byte:
				if isPrintableString(v) {
					log.Printf("%s\t", string(v))
				} else {
					log.Printf("%x\t", v) // Print binary data as hexadecimal
				}
			default:
				log.Printf("%v\t", v)
			}
		}
		log.Println()
	}

	if err = rows.Err(); err != nil {
		log.Println("Error during row iteration:", err)
	}
}

func runExec(db *sql.DB, query string) {
	result, err := db.Exec(query)
	if err != nil {
		log.Println("Error executing query against Edgeless DB:", err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting number of rows affected:", err)
		return
	}
	log.Println("Number of rows affected:", rowsAffected)
}

// isPrintableString checks if a byte slice contains only printable characters
func isPrintableString(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
}
