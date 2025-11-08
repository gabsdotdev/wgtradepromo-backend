package main

import (
	"fmt"
	"log"
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

	// List all CSV files in the current directory
	files, err := filepath.Glob("./*.csv")
	if err != nil {
		log.Fatalf("Error listing CSV files: %v", err)
	}

	if len(files) == 0 {
		log.Fatal("No CSV files found in the current directory")
	}

	// Process each CSV file
	var totalImported, totalSkipped int

	for fileIndex, csvPath := range files {
		fmt.Printf("\nProcessing file %d/%d: %s\n", fileIndex+1, len(files), filepath.Base(csvPath))

		// Parse CSV file
		stmt, err := parser.ParseCSV(csvPath)
		if err != nil {
			log.Printf("Error parsing CSV %s: %v - skipping file\n", csvPath, err)
			continue
		}

		// Validate account number
		if strings.TrimSpace(stmt.AccountNumber) == "" {
			log.Printf("Account number must be informed in the CSV %s (field Conta) - skipping file\n", csvPath)
			continue
		}

		// Get conta ID from database
		contaID, err := database.GetContaIDByNumero(stmt.AccountNumber)
		if err != nil {
			log.Printf("Error finding conta ID for %s: %v - skipping file\n", csvPath, err)
			continue
		}

		// Import transactions
		fmt.Printf("Importing %d transactions for conta %s...\n", len(stmt.Transactions), stmt.AccountNumber)

		var fileImported, fileSkipped int

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
					fileSkipped++
					continue
				}
				log.Printf("Error inserting transaction from %s: %v - skipping transaction\n", csvPath, err)
				continue
			}

			fileImported++
		}

		totalImported += fileImported
		totalSkipped += fileSkipped

		fmt.Printf("File import completed!\n")
		fmt.Printf("Imported: %d transactions\n", fileImported)
		fmt.Printf("Skipped (duplicates): %d transactions\n", fileSkipped)
		fmt.Printf("Total processed in this file: %d transactions\n", len(stmt.Transactions))
	}

	fmt.Printf("\n=== Final Import Summary ===\n")
	fmt.Printf("Total files processed: %d\n", len(files))
	fmt.Printf("Total transactions imported: %d\n", totalImported)
	fmt.Printf("Total transactions skipped: %d\n", totalSkipped)
	fmt.Printf("Total transactions processed: %d\n", totalImported+totalSkipped)
}

func isUniqueViolation(err error) bool {
	return err != nil && err.Error() == "duplicate key value violates unique constraint"
}
