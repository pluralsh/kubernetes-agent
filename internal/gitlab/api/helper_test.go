package api

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/entity"
	"google.golang.org/protobuf/testing/protocmp"
)

func AssertGitalyRepository(t *testing.T, gitalyRepository *entity.GitalyRepository, apiGitalyRepository *gitalypb.Repository) {
	assert.Equal(t, gitalyRepository.StorageName, apiGitalyRepository.StorageName)
	assert.Equal(t, gitalyRepository.RelativePath, apiGitalyRepository.RelativePath)
	assert.Equal(t, gitalyRepository.GitObjectDirectory, apiGitalyRepository.GitObjectDirectory)
	assert.Equal(t, gitalyRepository.GitAlternateObjectDirectories, apiGitalyRepository.GitAlternateObjectDirectories)
	assert.Equal(t, gitalyRepository.GlRepository, apiGitalyRepository.GlRepository)
	assert.Equal(t, gitalyRepository.GlProjectPath, apiGitalyRepository.GlProjectPath)
}

func AssertGitalyInfo(t *testing.T, gitalyInfo, apiGitalyInfo *entity.GitalyInfo) {
	assert.Empty(t, cmp.Diff(gitalyInfo, apiGitalyInfo, protocmp.Transform()))
}
