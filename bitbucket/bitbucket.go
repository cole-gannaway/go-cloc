package bitbucket

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"go-cloc/utilities"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Define the nested struct types
type item struct {
	Name    string  `json:"name"`
	Project project `json:"project"`
}
type project struct {
	Name string `json:"name"`
}

type response struct {
	Value []item `json:"values"`
	// Next is string or null
	Next *string `json:"next"`
}

const defaultBaseUrl = "bitbucket.org"

// Example URL: https://x-token-auth:accessToken@bitbucket.org/organization/repoName.git
func CreateCloneURLBitbucket(accessToken string, organization string, respository string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://x-token-auth:" + accessToken + "@" + baseUrl + "/" + organization + "/" + respository + ".git"
}

// Example URL: https://bitbucket.org/organization/repoName/get/HEAD.zip
func CreateZipURLBitbucket(accessToken string, organization string, repoName string, defaultBranch string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	logger.Warn("Assuming git repository since using HEAD.zip for the commit")
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/" + repoName + "/get/HEAD.zip"
}

// Example URL: https://api.bitbucket.org/2.0/repositories/organization?pagelen=100&page=1
func CreateDiscoverRepositoriesURLBitbucket(organization string, pageNum int, pageSize int, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://api." + baseUrl + "/2.0/repositories/" + organization + "?pagelen=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(pageNum)
}

func DiscoverReposBitbucket(organization string, accessToken string, devopsBaseUrlOverride string, useHttps bool) []devops.RepoInfo {
	pageSize := 100
	pageNum := 1
	repoNames := []devops.RepoInfo{}

	for pageNum != -1 {
		apiURL := CreateDiscoverRepositoriesURLBitbucket(organization, pageNum, pageSize, devopsBaseUrlOverride, useHttps)
		logger.Debug("Discovering repos using url: ", apiURL)

		// Create a new HTTP request
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf("Failed to create HTTP request: %v", err)
		}

		// Set basic auth
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
			logger.Error("Unexpected status code: ", resp.StatusCode, ", expected 200")
			logger.Error("Response: ", resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// Unmarshal the JSON data into the Response struct
		var r response
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		logger.Debug("Response: ", r)

		for _, item := range r.Value {
			repoInfo := devops.NewRepoInfo(organization, item.Project.Name, item.Name, "")
			repoNames = append(repoNames, repoInfo)
		}
		// If there is no next page, stop the loop
		if r.Next == nil {
			pageNum = -1
		} else {
			pageNum = pageNum + 1
		}
	}

	return repoNames
}
