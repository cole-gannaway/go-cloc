package utilities

import (
	"flag"
	"go-cloc/logger"
	"go-cloc/scanner"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Modes
const (
	LOCAL       string = "Local"
	GITHUB      string = "GitHub"
	AZUREDEVOPS string = "AzureDevOps"
	GITLAB      string = "GitLab"
	BITBUCKET   string = "Bitbucket"
)

type CLIArgs struct {
	LogLevel                        string
	LocalScanFilePath               string
	IgnorePatterns                  []string
	DumpCSVs                        bool
	ResultsDirectoryPath            string
	ScanId                          string
	OverrideLanguagesConfigFilePath string
}

func CleanLocalFilePath(targetPath string) string {
	logger.Debug("CleanLocalFilePath targetPath before: '", targetPath, "'")
	targetPath = filepath.Clean(targetPath)
	// On windows this may be needed if spaces are in the file path
	targetPath = strings.TrimSuffix(targetPath, "\"")
	logger.Debug("CleanLocalFilePath targetPath after: '", targetPath, "'")
	return targetPath
}

func ParseArgsFromCLI() CLIArgs {
	// print out arguments
	printLanguagesArg := flag.Bool("print-languages", false, "Prints out the supported languages, file suffixes, and comment configurations. Does not run the tool.")

	// optional arguments
	logLevelArg := flag.String("log-level", "INFO", "Log level - DEBUG, INFO, WARN, ERROR")
	scanIdArg := flag.String("scan-id", "", "Identifier for the scan. For reference in a csv file later")
	ignoreFilePathArg := flag.String("ignore-file", "", "Path to your ignore file. Defines directories and files to exclude when scanning. Please see the README.md for how to format your ignore configuration")
	dumpCSVsArg := flag.Bool("dump-csv", false, "When true, dumps results to a csv file, otherwise gives results in logs")
	resultsDirectoryPathArg := flag.String("results-directory-path", "", "Path to a new directory for storing the results. Default the tool will create one based on the start time")
	overrideLanguageConfigFilePathArg := flag.String("override-languages-path", "", "Path to languages configuration to override the default configuration.")

	// parse the CLI arguments
	flag.Parse()

	// dereference all CLI args to make it easier to use
	printLanguages := *printLanguagesArg

	// print out languages
	if printLanguages {
		scanner.PrintLanguages()
		os.Exit(0)
	}

	// Collect the remaining arguments
	cliArgs := flag.Args()

	// Ensure at least one argument
	if len(cliArgs) < 1 {
		logger.Error("Requires a path to the file or directory to scan as the first command line argument, ex: 'go-cloc file1.js'")
		os.Exit(-1)
	}

	// Parse any remaining flags after the first non-flag argument
	flag.CommandLine.Parse(cliArgs[1:])

	// dereference all CLI args to make it easier to use
	logLevel := *logLevelArg
	ignoreFilePath := *ignoreFilePathArg
	dumpCSVs := *dumpCSVsArg
	resultsDirectoryPath := *resultsDirectoryPathArg
	scanId := *scanIdArg
	overrideLanguageConfigFilePath := *overrideLanguageConfigFilePathArg

	// set log level
	logger.SetLogLevel(logger.ConvertStringToLogLevel(logLevel))
	logger.SetOutput(os.Stdout)

	logger.Info("Setting Log Level to " + logLevel)
	logger.Info("Parsing CLI arguments")

	// print out arguments
	logger.Debug("dump-csvs: ", dumpCSVs)

	logger.Debug("Validating mandatory arguments")

	// validate mandatory arguments
	if dumpCSVs && scanId == "" {
		logger.Error("Requires : --scan-id for --dump-csvs")
		os.Exit(-1)
	}

	// Set file path to scan
	localScanFilePath := CleanLocalFilePath(cliArgs[0])

	// validate optional arguments

	// parse ignore patterns
	ignorePatterns := []string{}
	if ignoreFilePath != "" {
		logger.Debug("Parsing ignore-file ", ignoreFilePath)
		ignorePatterns = scanner.ReadIgnoreFile(ignoreFilePath)
		logger.Debug("Successfully read in the ignore-file ", ignoreFilePath)
		logger.Debug("Ignore Patterns: ", ignorePatterns)
	}

	if !dumpCSVs && resultsDirectoryPath != "" {
		logger.Error("Cannot simultaneously set --results-directory-path and --dump-csvs=false")
		logger.LogStackTraceAndExit(nil)
	}

	// set results directory if dumpCSVs is true
	if resultsDirectoryPath == "" && dumpCSVs {
		resultsDirectoryPath = time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
		logger.Debug("Results Directory Path: ", resultsDirectoryPath)
	}

	// override languages config
	if overrideLanguageConfigFilePath != "" {
		logger.Debug("Overriding default languages with ", overrideLanguageConfigFilePath)
		scanner.LoadLanguages(overrideLanguageConfigFilePath)
	}

	args := CLIArgs{
		LogLevel:                        logLevel,
		LocalScanFilePath:               localScanFilePath,
		IgnorePatterns:                  ignorePatterns,
		DumpCSVs:                        dumpCSVs,
		ResultsDirectoryPath:            resultsDirectoryPath,
		ScanId:                          scanId,
		OverrideLanguagesConfigFilePath: overrideLanguageConfigFilePath,
	}

	return args
}
