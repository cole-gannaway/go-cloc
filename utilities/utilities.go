package utilities

import (
	"flag"
	"go-cloc/logger"
	"go-cloc/scanner"
	"os"
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
	LogLevel              string
	Mode                  string
	LocalScanFilePath     string
	AccessToken           string
	Organization          string
	IgnorePatterns        []string
	ExcludeRepositories   []string
	IncludeRepositories   []string
	CloneRepoUsingZip     bool
	DumpCSVs              bool
	ResultsDirectoryPath  string
	UseHTTPS              bool
	DevopsBaseUrlOverride string
}

func ParseArgsFromCLI() CLIArgs {

	// print out arguments
	printLanguagesArg := flag.Bool("print-languages", false, "Prints out the supported languages, file suffixes, and comment configurations. Does not run the tool.")
	// mandatory arguments
	modeArg := flag.String("devops", LOCAL, "GitHub, AzureDevOps, Bitbucket, GitLab, or Local")
	accessTokenArg := flag.String("accessToken", "", "Your DevOps personal access token used for discovering and downloading repositories in your organization")
	organizationArg := flag.String("organization", "", "Your DevOps organization name")
	// optional arguments
	logLevelArg := flag.String("log-level", "INFO", "Log level - DEBUG, INFO, WARN, ERROR")
	localScanFilePathArg := flag.String("local-file-path", "", "Path to youthe local file or directory that you wish to scan")
	ignoreFilePathArg := flag.String("ignore-file", "", "Path to your ignore file. Defines directories and files to exclude when scanning. Please see the README.md for how to format your ignore configuration")
	excludeRepositoriesFilePathArg := flag.String("exclude-repositories-file", "", "Path to your exclude repositories file. Defines repositories to exclude, all others will be included. Please see the README.md for how to format your exclude repositories configuration")
	includeRepositoriesFilePathArg := flag.String("include-repositories-file", "", "Path to your include repositories file. Defines repositories to include, all others will be excluded. Please see the README.md for how to format your include repositories configuration")
	cloneRepoUsingZipArg := flag.Bool("clone-repo-using-zip", false, "When true, repositories are downloaded as zip files instead of git clone to drastically improve performance. Default is false. This is a BETA feature and has not been extensively tested")
	dumpCSVsArg := flag.Bool("dump-csvs", true, "When false, disables csv file dumps. DEBUG logging available to still see csv results in logs.")
	resultsDirectoryPathArg := flag.String("results-directory-path", "", "Path to a new directory for storing the results. Default the tool will create one based on the start time")
	useHttpsArg := flag.Bool("use-https", true, "When false, uses http instead of https for all HTTP calls.")
	devopsBaseUrlOverrideArg := flag.String("devops-base-url-override", "", "Overrides the base URL for the DevOps provider. Defaults will be github.com, dev.azure.com, bitbucket.org, or gitlab.com. However, you can override this with your own self-hosted ip or domain")

	// parse the CLI arguments
	flag.Parse()

	// dereference all CLI args to make it easier to use
	printLanguages := *printLanguagesArg
	logLevel := *logLevelArg
	mode := *modeArg
	localScanFilePath := *localScanFilePathArg
	accessToken := *accessTokenArg
	organization := *organizationArg
	ignoreFilePath := *ignoreFilePathArg
	excludeRepositoriesFilePath := *excludeRepositoriesFilePathArg
	includeRepositoriesFilePath := *includeRepositoriesFilePathArg
	cloneRepoUsingZip := *cloneRepoUsingZipArg
	dumpCSVs := *dumpCSVsArg
	resultsDirectoryPath := *resultsDirectoryPathArg
	useHttps := *useHttpsArg
	devopsBaseUrlOverride := *devopsBaseUrlOverrideArg

	// set log level
	logger.SetLogLevel(logger.ConvertStringToLogLevel(logLevel))
	logger.SetOutput(os.Stdout)

	logger.Info("Setting Log Level to " + logLevel)
	logger.Info("Parsing CLI arguments")

	// print out arguments
	logger.Debug("Mode: ", mode)
	logger.Debug("clone-repo-using-zip: ", cloneRepoUsingZip)
	logger.Debug("dump-csvs: ", dumpCSVs)

	// print out languages
	if printLanguages {
		scanner.PrintLanguages()
		os.Exit(0)
	}

	// validate mandatory arguments
	logger.Debug("Validating mandatory arguments")
	if mode == LOCAL {
		if localScanFilePath == "" {
			logger.Error("Mode ", mode, " requires : --local-file-path")
			os.Exit(-1)
		}
	} else {
		if organization == "" || accessToken == "" {
			logger.Error("Mode ", mode, " requires : --organization & --accessToken")
			os.Exit(-1)
		}
	}

	// validate optional arguments

	// parse ignore patterns
	ignorePatterns := []string{}
	if ignoreFilePath != "" {
		logger.Debug("Parsing ignore-file ", ignoreFilePath)
		ignorePatterns = scanner.ReadIgnoreFile(ignoreFilePath)
		logger.Debug("Successfully read in the ignore-file ", ignoreFilePath)
		logger.Debug("Ignore Patterns: ", ignorePatterns)
	}

	// parse exclude repositories
	excludeRepositories := []string{}
	if excludeRepositoriesFilePath != "" {
		logger.Debug("Parsing exclude-repositories-file ", excludeRepositoriesFilePath)
		excludeRepositories = scanner.ReadIgnoreFile(excludeRepositoriesFilePath)
		logger.Debug("Successfully read in the exclude-repositories-file ", excludeRepositoriesFilePath)
		logger.Debug("Exclude Repositories: ", excludeRepositories)
	}

	// parse include repositories
	includeRepositories := []string{}
	if includeRepositoriesFilePath != "" {
		logger.Debug("Parsing include-repositories-file ", includeRepositoriesFilePath)
		includeRepositories = scanner.ReadIgnoreFile(includeRepositoriesFilePath)
		logger.Debug("Successfully read in the include-repositories-file ", includeRepositoriesFilePath)
		logger.Debug("Include Repositories: ", includeRepositories)
	}

	if !dumpCSVs && resultsDirectoryPath != "" {
		logger.Error("Cannot simultaneously set --results-directory-path and --dump-csvs=false")
		logger.LogStackTraceAndExit(nil)
	}

	// set results directory if dumpCSVs is true
	if resultsDirectoryPath == "" && dumpCSVs {
		resultsDirectoryPath = time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	}
	logger.Debug("Results Directory Path: ", resultsDirectoryPath)

	args := CLIArgs{
		LogLevel:              logLevel,
		Mode:                  mode,
		LocalScanFilePath:     localScanFilePath,
		AccessToken:           accessToken,
		Organization:          organization,
		IgnorePatterns:        ignorePatterns,
		ExcludeRepositories:   excludeRepositories,
		IncludeRepositories:   includeRepositories,
		CloneRepoUsingZip:     cloneRepoUsingZip,
		DumpCSVs:              dumpCSVs,
		ResultsDirectoryPath:  resultsDirectoryPath,
		UseHTTPS:              useHttps,
		DevopsBaseUrlOverride: devopsBaseUrlOverride,
	}

	return args
}

func GetHttpProtocolSetting(useHttps bool) string {
	if useHttps {
		return "https"
	}
	return "http"
}
