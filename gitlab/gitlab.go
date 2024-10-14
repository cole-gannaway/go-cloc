package gitlab

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"go-cloc/utilities"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const defaultBaseUrl = "gitlab.com"

// Example URL: https://oauth2:accessToken@gitlab.com/organization/repoName.git
func CreateCloneURLGitLab(accessToken string, organization string, respository string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	// Create the URL
	return httpProtocolSetting + "://oauth2:" + accessToken + "@" + baseUrl + "/" + organization + "/" + respository + ".git"
}

// Example URL: https://accesstoken@gitlab.com/api/v4/groups/organization/projects?per_page=100&page=1
func CreateDiscoverURLGitLab(accessToken string, organization string, pageNum int, pageSize int, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + accessToken + "@" + baseUrl + "/api/v4/groups/" + organization + "/projects?per_page=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}

// Example URL: https://gitlab.com/organization/repoName/-/archive/defaultBranch/repoName-defaultBranch.zip
func CreateZipURLGitLab(organization string, repoName string, defaultBranch string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/" + repoName + "/-/archive/" + defaultBranch + "/" + repoName + "-" + defaultBranch + ".zip"
}

// Define the nested struct types
type item struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

func DiscoverReposGitlab(organization string, accessToken string, devopsBaseUrlOverride string, useHttps bool) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}
	// pageNum -1 means there are no more pages to discover
	for pageNum != -1 {
		apiURL := CreateDiscoverURLGitLab(accessToken, organization, pageNum, pageSize, devopsBaseUrlOverride, useHttps)
		logger.Debug("Discovering repos using url: ", apiURL)

		// Create a new HTTP request
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf("Failed to create HTTP request: %v", err)
		}

		// Add the Authorization header
		req.Header.Set("Authorization", "Bearer "+accessToken)

		// Perform the request using the default HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to fetch data from API: %v", err)
		}
		defer resp.Body.Close()

		// Check if the status code is 200
		if resp.StatusCode != http.StatusOK {
			logger.Error("Failed to fetch data from API: ", resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// Parse the JSON response into a slice of Item
		var result []item
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Failed to parse JSON: %v", err)
		}

		// Print the parsed data
		for _, item := range result {
			repoInfo := devops.NewRepoInfo(organization, "", item.Name, item.DefaultBranch)
			repoNames = append(repoNames, repoInfo)
		}

		// Get the next page URL
		link := resp.Header.Get("Link")
		logger.Debug("Link header: ", link)
		logger.Debug(strings.Contains(link, `rel="last"`))
		// If there is no next page, stop the loop
		if link == "" || !strings.Contains(link, `rel="next"`) {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}
