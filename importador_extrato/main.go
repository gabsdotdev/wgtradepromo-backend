package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabsdotdev/wgtradepromo-backend/importador_extrato/config"
	"github.com/gabsdotdev/wgtradepromo-backend/importador_extrato/db"
	"github.com/gabsdotdev/wgtradepromo-backend/importador_extrato/parser"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	database, err := db.NewConnection(cfg.GetConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	// Get CSV file path
	csvPath := filepath.Join(".", "Extrato-01-01-2024-a-07-11-2025-CSV.csv")

	// Parse CSV file
	stmt, err := parser.ParseCSV(csvPath)
	if err != nil {
		log.Fatalf("Error parsing CSV: %v", err)
	}

	// Validate account number
	if strings.TrimSpace(stmt.AccountNumber) == "" {
		log.Fatal("Account number must be informed in the CSV (field Conta).")
	}

	// Get conta ID from database
	contaID, err := database.GetContaIDByNumero(stmt.AccountNumber)
	if err != nil {
		log.Fatalf("Error finding conta ID: %v", err)
	}

	// Import transactions
	fmt.Printf("Importing %d transactions for conta %s...\n", len(stmt.Transactions), stmt.AccountNumber)

	var importedCount, skippedCount int

	for _, tx := range stmt.Transactions {
		err := database.InsertTransaction(
			contaID,
			tx.Date,
			tx.Description,
			tx.Details,
			tx.Amount,
		)

		if err != nil {
			if isUniqueViolation(err) {
				skippedCount++
				continue
			}
			log.Fatalf("Error inserting transaction: %v", err)
		}

		importedCount++
	}

	fmt.Printf("Import completed!\n")
	fmt.Printf("Imported: %d transactions\n", importedCount)
	fmt.Printf("Skipped (duplicates): %d transactions\n", skippedCount)
	fmt.Printf("Total processed: %d transactions\n", len(stmt.Transactions))
}

func isUniqueViolation(err error) bool {
	return err != nil && err.Error() == "duplicate key value violates unique constraint"
}

func init() {
	// Ensure CSV file exists
	if _, err := os.Stat("Extrato-01-01-2024-a-07-11-2025-CSV.csv"); os.IsNotExist(err) {
		log.Fatal("CSV file not found: Extrato-01-01-2024-a-07-11-2025-CSV.csv")
	}
}
