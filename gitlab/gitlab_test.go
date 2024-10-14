package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_gitlab_CreateCloneURLGitLab(t *testing.T) {
	organization := "organization"
	repository := "respository"
	accessToken := "accesstoken"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateCloneURLGitLab(accessToken, organization, repository, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://oauth2:accesstoken@gitlab.com/organization/respository.git", cloneUrl)
}
func Test_gitlab_CreateDiscoverURLGitLab(t *testing.T) {
	organization := "organization"
	pageNum := 1
	pageSize := 100
	accessToken := "accesstoken"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateDiscoverURLGitLab(accessToken, organization, pageNum, pageSize, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://accesstoken@gitlab.com/api/v4/groups/organization/projects?per_page=100&page=1", cloneUrl)
}
func Test_gitlab_CreateZipURLGitLab(t *testing.T) {
	organization := "organization"
	repository := "repository"
	defaultBranch := "main"
	devOpsBaseUrlOverride := ""
	useHttps := true

	cloneUrl := CreateZipURLGitLab(organization, repository, defaultBranch, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://gitlab.com/organization/repository/-/archive/main/repository-main.zip", cloneUrl)
}
