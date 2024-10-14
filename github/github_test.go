package github

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_github_CreateCloneURLGithub(t *testing.T) {
	organization := "organization"
	repository := "repository"
	accessToken := "accesstoken"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateCloneURLGithub(accessToken, organization, repository, devOpsBaseUrlOverride, useHttps)

	fmt.Println(cloneUrl)
	// Assert
	assert.Equal(t, "https://oauth2:accesstoken@github.com/organization/repository.git", cloneUrl)
}

func Test_github_CreateZipURLGithub(t *testing.T) {

	organization := "organization"
	repository := "repository"
	defaultBranch := "main"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateZipURLGithub(organization, repository, defaultBranch, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://github.com/organization/repository/archive/refs/heads/main.zip", cloneUrl)
}

func Test_github_CreateGetDefaultBranchURLGitHub(t *testing.T) {
	organization := "organization"
	repository := "repository"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateGetDefaultBranchURLGitHub(organization, repository, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://api.github.com/repos/organization/repository", cloneUrl)
}

func Test_github_CreateDiscoverURLGitHub(t *testing.T) {
	organization := "organization"
	pageNum := 1
	pageSize := 100
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateDiscoverReposURLGitHub(organization, pageNum, pageSize, devOpsBaseUrlOverride, useHttps)

	// Assert
	assert.Equal(t, "https://api.github.com/orgs/organization/repos?per_page=100&page=1", cloneUrl)
}
