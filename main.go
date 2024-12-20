package main

import (
	"fmt"
	"go-cloc/logger"
	"go-cloc/report"
	"go-cloc/scanner"
	"go-cloc/utilities"
	"os"
	"path/filepath"
)

func main() {
	// parse CLI arguments and store them in a struct
	args := utilities.ParseArgsFromCLI()

	// create output folder
	if args.DumpCSVs {
		// only create the folder if the folder does not exist
		if _, err := os.Stat(args.ResultsDirectoryPath); err == nil {
			logger.Debug("Folder ", args.ResultsDirectoryPath, " already exists")
		} else if os.IsNotExist(err) {
			os.Mkdir(args.ResultsDirectoryPath, 0777)
		}
	}

	// scan LOC for the directory
	logger.Info("Scanning ", args.LocalScanFilePath, "...")
	filePaths := scanner.WalkDirectory(args.LocalScanFilePath, args.IgnorePatterns)
	fileScanResultsArr := []scanner.FileScanResults{}
	for _, filePath := range filePaths {
		fileScanResultsArr = append(fileScanResultsArr, scanner.ScanFile(filePath))
	}

	logger.Debug("Calculating total LOC ...")

	// sort and calculate total LOC
	fileScanResultsArr = report.SortFileScanResults(fileScanResultsArr)
	repoTotalResult := report.CalculateTotalLineOfCode(fileScanResultsArr)

	// convert results into records for CSV or command line output
	records := report.ConvertFileResultsIntoRecords(fileScanResultsArr, repoTotalResult)

	// Dump results by file in a csv
	if args.DumpCSVs {
		outputCsvFilePath := filepath.Join(args.ResultsDirectoryPath, args.ScanId+".csv")
		logger.Debug("Dumping results by file to ", outputCsvFilePath)
		report.WriteCsv(outputCsvFilePath, records)
		logger.Info("Done! Results for ", args.ScanId, " can be found ", outputCsvFilePath)
	} else {
		// print results to the command line
		logger.Info("Results by file for ", args.LocalScanFilePath, ":")
		report.PrintCsv(records)
	}

	logger.Info("Total LOC for ", args.LocalScanFilePath, " is ", repoTotalResult.CodeLineCount)

	// Print the total LOC to standard output to make it easy for external tools to parse
	fmt.Println(repoTotalResult.CodeLineCount)
}
