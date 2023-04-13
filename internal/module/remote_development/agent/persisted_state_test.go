package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testVersion = &WorkspaceRailsInfo{
	Name:                               "workspace-123",
	Namespace:                          "namespace",
	PersistedDeploymentResourceVersion: "123",
	PersistedActualStateIsTerminated:   false,
	DesiredStateIsTerminated:           false,
	ConfigToApply:                      "",
}

func TestRecordAndRetrieve(t *testing.T) {
	tracker := newPersistedStateTracker()
	tracker.recordVersion(testVersion)

	assert.True(t, tracker.isPersisted("workspace-123", "123"))
}

func TestDeleteWorkspace(t *testing.T) {
	tracker := newPersistedStateTracker()
	tracker.recordVersion(testVersion)

	tracker.delete("workspace-123")

	assert.False(t, tracker.isPersisted("workspace-123", "123"))
}
