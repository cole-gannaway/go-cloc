package azuredevops

import (
	"encoding/json"
	"go-cloc/devops"
	"go-cloc/logger"
	"go-cloc/utilities"
	"io"
	"log"
	"net/http"
)

// Define the nested struct types
type item struct {
	Name string `json:"name"`
}

type response struct {
	Value []item `json:"value"`
}

const defaultBaseUrl = "dev.azure.com"

// Example URL: https://oauth2:accessToken@dev.azure.com/organization/projectName/_git/repoName
func CreateCloneURLAzureDevOps(accessToken string, organization string, projectName string, repoName string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + accessToken + "@" + baseUrl + "/" + organization + "/" + projectName + "/_git/" + repoName
}

// Example URL: https://dev.azure.com/organization/projectName/_apis/git/repositories/repoName/items/items?path=/&versionDescriptor[versionOptions]=0&versionDescriptor[versionType]=0&versionDescriptor[version]=defaultBranch&resolveLfs=true&$format=zip&api-version=5.0&download=true
func CreateZipURLAzureDevOps(organization string, projectName string, repoName string, defaultBranch string, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/" + projectName + "/_apis/git/repositories/" + repoName + "/items/items?path=/&versionDescriptor[versionOptions]=0&versionDescriptor[versionType]=0&versionDescriptor[version]=" + defaultBranch + "&resolveLfs=true&$format=zip&api-version=5.0&download=true"
}

// Example URL: https://dev.azure.com/organization/_apis/projects?api-version=7.0
func CreateDiscoverProjectsURLAzureDevOps(organization string, pageNum int, pageSize int, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/_apis/projects?api-version=7.0"
}

// Example URL: https://dev.azure.com/organization/projectName/_apis/git/repositories?api-version=7.0
func CreateDiscoverRepositoriesURLAzureDevOps(organization string, projectName string, pageNum int, pageSize int, devopsBaseUrlOverride string, useHttps bool) string {
	httpProtocolSetting := utilities.GetHttpProtocolSetting(useHttps)
	baseUrl := defaultBaseUrl
	if devopsBaseUrlOverride != "" {
		baseUrl = devopsBaseUrlOverride
	}
	return httpProtocolSetting + "://" + baseUrl + "/" + organization + "/" + projectName + "/_apis/git/repositories?api-version=7.0"
}

func DiscoverReposAzureDevOps(organization string, accessToken string, devopsBaseUrlOverride string, useHttps bool) []devops.RepoInfo {
	discoverProjectsUrl := CreateDiscoverProjectsURLAzureDevOps(organization, 1, 100, devopsBaseUrlOverride, useHttps)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", discoverProjectsUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set basic auth
	req.SetBasicAuth("", accessToken)

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

	repoNames := []devops.RepoInfo{}
	// Access the nested Name field
	for _, item := range r.Value {
		projectName := item.Name
		logger.Debug("Project Name:", projectName)

		discoverRepositoriesUrl := CreateDiscoverRepositoriesURLAzureDevOps(organization, projectName, 1, 100, devopsBaseUrlOverride, useHttps)
		// Create a new HTTP request
		req, err := http.NewRequest("GET", discoverRepositoriesUrl, nil)
		if err != nil {
			log.Fatalf("Failed to create HTTP request: %v", err)
		}

		// Set basic auth
		req.SetBasicAuth("", accessToken)

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

		r := response{}
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}
		for _, item := range r.Value {
			repoName := item.Name
			repoInfo := devops.NewRepoInfo(organization, projectName, repoName, "")
			repoNames = append(repoNames, repoInfo)
		}
	}

	return repoNames
}
