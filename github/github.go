package github

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"go-cloc/utilities"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Define a struct with only the fields you care about
type item struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

type repo struct {
	DefaultBranch string `json:"default_branch"`
}

const defaultBaseUrl = "github.com"

// Example URL: https://oauth2:accesstoken@github.com/organization/repoName.git
func CreateCloneURLGithub(accessToken string, organization string, repoName string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://oauth2:" + accessToken + "@" + baseUrl + "/" + organization + "/" + repoName + ".git"
}

// Example URL: https://github.com/organization/repoName/archive/refs/heads/defaultBranch.zip
func CreateZipURLGithub(organization string, repoName string, defaultBranch string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/" + repoName + "/archive/refs/heads/" + defaultBranch + ".zip"
}

// Example URL: https://api.github.com/repos/organization/repoName
func CreateGetDefaultBranchURLGitHub(organization string, repoName string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://api." + baseUrl + "/repos/" + organization + "/" + repoName
}

// Example URL: https://api.github.com/orgs/organization/repos?per_page=100&page=1
func CreateDiscoverReposURLGitHub(organization string, pageNum int, pageSize int, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://api." + baseUrl + "/orgs/" + organization + "/repos?per_page=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}

func DiscoverReposGithub(organization string, accessToken string, devopsBaseUrlOverride string, useHttps bool) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}

	// pageNum -1 means there are no more pages to discover
	for pageNum != -1 {
		apiURL := CreateDiscoverReposURLGitHub(organization, pageNum, pageSize, devopsBaseUrlOverride, useHttps)
		logger.Debug("GET: " + apiURL)

		// Create a new HTTP request
		req, _ := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		// Perform the request using the default HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.LogStackTraceAndExit(err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Failed to read response body: ", body)
			logger.LogStackTraceAndExit(err)
		}

		// Check if the status code is 200
		if resp.StatusCode != http.StatusOK {
			logger.Error("Response status code: ", resp.StatusCode, " expected 200")
			logger.LogStackTraceAndExit(nil)
		}

		// Parse the JSON response into a slice of Item
		var result []item
		if err := json.Unmarshal(body, &result); err != nil {
			logger.Error("Failed to parse JSON: ", err)
			logger.LogStackTraceAndExit(err)
		}

		// Print the parsed data
		logger.Debug("Default branch is: ", result)

		for _, item := range result {
			repoInfo := devops.NewRepoInfo(organization, "", item.Name, item.DefaultBranch)
			repoNames = append(repoNames, repoInfo)
		}

		// Get the next page URL
		link := resp.Header.Get("Link")
		logger.Debug("Link header: ", link)
		// If there is no next page, stop the loop
		if link == "" || !strings.Contains(link, `rel="next"`) {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}

func DiscoverDefaultBranchForRepoGithub(organization string, repoName string, accessToken string, devopsBaseUrlOverride string, useHttps bool) string {
	logger.Debug("Getting default branch for ", organization, "/", repoName)

	url := CreateGetDefaultBranchURLGitHub(organization, repoName, devopsBaseUrlOverride, useHttps)
	logger.Debug("GET: " + url)

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogStackTraceAndExit(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body: ", body)
		logger.LogStackTraceAndExit(err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Error("Response status code: ", resp.StatusCode, " expected 200")
		logger.LogStackTraceAndExit(nil)
	}

	// Parse the JSON response into a slice of Item
	var repoResult repo
	if err := json.Unmarshal(body, &repoResult); err != nil {
		logger.Error("Failed to parse JSON: ", err)
		logger.LogStackTraceAndExit(err)
	}

	// Print the parsed data
	logger.Debug("Default branch is: ", repoResult.DefaultBranch)

	return repoResult.DefaultBranch
}
