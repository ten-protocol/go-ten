package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/ten-protocol/go-ten/go/enclave/storage/init/edgelessdb"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
)

func main() {
	// get edbHost from first command line arg
	var edbHost string
	if len(os.Args) > 1 {
		edbHost = os.Args[1]
	} else {
		fmt.Println("Usage: edbconnect <edb-host>")
		fmt.Println("Ensure you have the latest copy of the ./edb-connect.sh launch script if you see this error.")
		os.Exit(1)
	}

	fmt.Println("Retrieving Edgeless DB credentials...")
	creds, found, err := edgelessdb.LoadCredentialsFromFile()
	if err != nil {
		fmt.Println("Error loading credentials from file:", err)
		panic(err)
	}
	if !found {
		panic("No existing EDB credentials found.")
	}
	fmt.Println("Found existing EDB credentials. Creating TLS config...")
	cfg, err := edgelessdb.CreateTLSCfg(creds)
	if err != nil {
		fmt.Println("Error creating TLS config from credentials:", err)
		panic(err)
	}
	fmt.Println("TLS config created. Connecting to Edgeless DB...")
	testlog.SetupSysOut()
	db, err := edgelessdb.ConnectToEdgelessDB(edbHost, cfg, testlog.Logger())
	if err != nil {
		fmt.Printf("Error connecting to Edgeless DB at %s: %v\n", edbHost, err)
		panic(err)
	}
	fmt.Println("Connected to Edgeless DB.")

	startREPL(db.DB)

	err = db.Close()
	if err != nil {
		fmt.Println("Error closing Edgeless DB connection:", err)
		panic(err)
	}
}

// Starts a loop that reads user input and runs queries against the Edgeless DB until user types "exit"
func startREPL(db *sql.DB) {
	for {
		fmt.Println("\nEnter a query to run against the Edgeless DB (or type 'exit' to quit):")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>> ") // Display the prompt

		query, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading user input:", err)
			continue
		}
		// Trim the newline character and surrounding whitespace
		query = strings.TrimSpace(query)

		// line break for readability
		fmt.Println()

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
	fmt.Println("Exiting...")
}

func runQuery(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Error fetching columns:", err)
		return
	}

	// Print column headers
	for _, colName := range cols {
		fmt.Printf("%s\t", colName)
	}
	fmt.Println()

	// Prepare a slice to hold the values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		// Print the row values
		for _, val := range values {
			// Handle NULL values and convert byte slices to strings
			switch v := val.(type) {
			case nil:
				fmt.Print("NULL\t")
			case []byte:
				if isPrintableString(v) {
					fmt.Printf("%s\t", string(v))
				} else {
					fmt.Printf("%x\t", v) // Print binary data as hexadecimal
				}
			default:
				fmt.Printf("%v\t", v)
			}
		}
		fmt.Println()
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error during row iteration:", err)
	}
}

func runExec(db *sql.DB, query string) {
	result, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error executing query against Edgeless DB:", err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting number of rows affected:", err)
		return
	}
	fmt.Println("Number of rows affected:", rowsAffected)
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
