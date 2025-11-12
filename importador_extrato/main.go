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

	// Create parser factory
	factory := parser.NewParserFactory()

	fmt.Println("=== Importador de Extratos ===")
	fmt.Printf("Parsers disponíveis: %s\n\n", strings.Join(factory.ListSupportedParsers(), ", "))

	// List all supported files in the rawdata/extrato directory
	files, err := findImportableFiles("./rawdata/extrato")
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	if len(files) == 0 {
		log.Fatal("No importable files found in ./rawdata/extrato")
	}

	fmt.Printf("Found %d file(s) to process\n\n", len(files))

	// Process each file
	var totalImported, totalSkipped, totalErrors int

	for fileIndex, filePath := range files {
		fmt.Printf("\n=== Processing file %d/%d: %s ===\n", fileIndex+1, len(files), filepath.Base(filePath))

		// Get appropriate parser for this file
		p, err := factory.GetParser(filePath)
		if err != nil {
			log.Printf("Skipping file %s: %v\n", filepath.Base(filePath), err)
			totalErrors++
			continue
		}

		fmt.Printf("Using parser: %s\n", p.GetName())

		// Parse file
		stmt, err := p.Parse(filePath)
		if err != nil {
			log.Printf("Error parsing file %s: %v - skipping file\n", filepath.Base(filePath), err)
			totalErrors++
			continue
		}

		// Validate account number
		if strings.TrimSpace(stmt.AccountNumber) == "" {
			log.Printf("Account number must be informed in the file %s - skipping file\n", filepath.Base(filePath))
			totalErrors++
			continue
		}

		var fileImported, fileSkipped int

		switch stmt.AccountNumber {
		case "simples-nacional":
			// Get conta ID from database
			empresaID, err := database.GetEmpresaIDByCNPJ(stmt.DasDocumento.CNPJ)
			if err != nil {
				log.Printf("Error finding conta ID for %s: %v - skipping file\n", stmt.AccountNumber, err)
				totalErrors++
				continue
			}

			// Import das ducumento
			fmt.Printf("Importing DAS documento for empresa %s...\n", empresaID)

			err = database.InsertDasDocumento(
				empresaID,
				stmt.DasDocumento.PeriodoApuracao,
				stmt.DasDocumento.DataVencimento,
				stmt.DasDocumento.NumeroDocumento,
				stmt.DasDocumento.ValorTotal,
			)

			if err != nil {
				if isUniqueViolation(err) {
					fileSkipped++
					continue
				}
				log.Printf("Error inserting das documento: %v - skipping \n", err)
				continue
			}

			fileImported++

		default:
			// Get conta ID from database
			contaID, err := database.GetContaIDByNumero(stmt.AccountNumber)
			if err != nil {
				log.Printf("Error finding conta ID for %s: %v - skipping file\n", stmt.AccountNumber, err)
				totalErrors++
				continue
			}

			// Import transactions
			fmt.Printf("Importing %d transaction(s) for account %s...\n", len(stmt.Transactions), stmt.AccountNumber)

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
					log.Printf("Error inserting transaction: %v - skipping transaction\n", err)
					continue
				}

				fileImported++
			}

			totalImported += fileImported
			totalSkipped += fileSkipped

		}

		fmt.Printf("✓ File import completed!\n")
		fmt.Printf("  Imported: %d transaction(s)\n", fileImported)
		fmt.Printf("  Skipped (duplicates): %d transaction(s)\n", fileSkipped)
		fmt.Printf("  Total processed: %d transaction(s)\n", len(stmt.Transactions))
	}

	fmt.Printf("\n=== Final Import Summary ===\n")
	fmt.Printf("Files processed: %d\n", len(files))
	fmt.Printf("Files with errors: %d\n", totalErrors)
	fmt.Printf("Transactions imported: %d\n", totalImported)
	fmt.Printf("Transactions skipped: %d\n", totalSkipped)
	fmt.Printf("Total transactions processed: %d\n", totalImported+totalSkipped)
}

// findImportableFiles busca recursivamente por arquivos suportados
func findImportableFiles(rootDir string) ([]string, error) {
	var files []string

	// Busca por CSVs (Inter e Nubank)
	csvFiles, err := filepath.Glob(filepath.Join(rootDir, "**/*.csv"))
	if err != nil {
		return nil, fmt.Errorf("error listing CSV files: %v", err)
	}
	files = append(files, csvFiles...)

	// Busca por PDFs (Simples Nacional)
	pdfFiles, err := filepath.Glob(filepath.Join(rootDir, "**/*.pdf"))
	if err != nil {
		return nil, fmt.Errorf("error listing PDF files: %v", err)
	}
	files = append(files, pdfFiles...)

	// Como glob com ** não funciona sempre, vamos fazer busca manual
	if len(files) == 0 {
		files, err = findFilesRecursive(rootDir, []string{".csv", ".pdf"})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

// findFilesRecursive busca arquivos recursivamente
func findFilesRecursive(rootDir string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
			if ext == validExt {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return files, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && err.Error() == "duplicate key value violates unique constraint"
}
