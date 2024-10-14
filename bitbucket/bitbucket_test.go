package bitbucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bitbucket_CreateDiscoverURLBitbucket(t *testing.T) {

	organization := "organization"
	pageNum := 1
	pageSize := 100
	devOpsBaseUrlOverride := ""
	useHttps := true
	apiURL := CreateDiscoverRepositoriesURLBitbucket(organization, pageNum, pageSize, devOpsBaseUrlOverride, useHttps)
	expected := "https://api.bitbucket.org/2.0/repositories/organization?pagelen=100&page=1"
	// Assert
	assert.Equal(t, expected, apiURL)
}
func Test_bitbucket_CreateCloneURLBitbucket(t *testing.T) {

	organization := "organization"
	repository := "repository"
	accessToken := "accessToken"
	devOpsBaseUrlOverride := ""
	useHttps := true
	cloneUrl := CreateCloneURLBitbucket(accessToken, organization, repository, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://x-token-auth:accessToken@bitbucket.org/organization/repository.git", cloneUrl)
}

func Test_bitbucket_CreateZipURLBitbucket(t *testing.T) {
	organization := "organization"
	repository := "repository"
	accessToken := "accessToken"
	defaultBranch := "main"
	devOpsBaseUrlOverride := ""
	useHttps := true
	cloneUrl := CreateZipURLBitbucket(accessToken, organization, repository, defaultBranch, devOpsBaseUrlOverride, useHttps)
	// Assert
	assert.Equal(t, "https://bitbucket.org/organization/repository/get/HEAD.zip", cloneUrl)
}
