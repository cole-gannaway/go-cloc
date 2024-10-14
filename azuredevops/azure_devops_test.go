package azuredevops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_azuredevops_CreateCloneURLAzureDevOps(t *testing.T) {
	accessToken := "abcdefg"
	organization := "organization"
	projectName := "project"
	repoName := "repo"
	devOpsBaseUrlOverride := ""
	useHttps := true
	azdoCloneURL := CreateCloneURLAzureDevOps(accessToken, organization, projectName, repoName, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://abcdefg@dev.azure.com/organization/project/_git/repo", azdoCloneURL)
}
func Test_azuredevops_CreateZipURLAzureDevOps(t *testing.T) {
	organization := "organization"
	projectName := "project"
	repository := "repository"
	defaultBranch := "main"
	devOpsBaseUrlOverride := ""
	useHttps := true
	azdoCloneURL := CreateZipURLAzureDevOps(organization, projectName, repository, defaultBranch, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://dev.azure.com/organization/project/_apis/git/repositories/repository/items/items?path=/&versionDescriptor[versionOptions]=0&versionDescriptor[versionType]=0&versionDescriptor[version]=main&resolveLfs=true&$format=zip&api-version=5.0&download=true", azdoCloneURL)
}
func Test_azuredevops_CreateDiscoverProjectsURLAzureDevOps(t *testing.T) {
	organization := "organization"
	pageSize := 100
	pageNum := 1
	devOpsBaseUrlOverride := ""
	useHttps := true
	azdoCloneURL := CreateDiscoverProjectsURLAzureDevOps(organization, pageNum, pageSize, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://dev.azure.com/organization/_apis/projects?api-version=7.0", azdoCloneURL)
}

func Test_azuredevops_CreateDiscoverRepositoriesURLAzureDevOps(t *testing.T) {
	organization := "organization"
	projectName := "project"
	pageSize := 100
	pageNum := 1
	devOpsBaseUrlOverride := ""
	useHttps := true
	azdoCloneURL := CreateDiscoverRepositoriesURLAzureDevOps(organization, projectName, pageNum, pageSize, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://dev.azure.com/organization/project/_apis/git/repositories?api-version=7.0", azdoCloneURL)
}

//
//
