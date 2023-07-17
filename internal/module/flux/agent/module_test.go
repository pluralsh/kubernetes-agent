package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_k8s"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/mock/gomock"
	apiextensionsv1api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	testSupportedCrd             = requiredFluxCrds[0]
	testSupportedCrdResourceName = testSupportedCrd.GroupResource().String()
	testSupportedCrdVersion      = testSupportedCrd.Version
)

var (
	_ modagent.Module = &module{}
)

func TestModule_DefaultAndValidateConfiguration_WithoutFluxConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, defaultServiceApiBaseUrl, cfg.Flux.WebhookReceiverUrl)
}

func TestModule_DefaultAndValidateConfiguration_WithoutWebhookReceiverUrlConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{
		Flux: &agentcfg.FluxCF{},
	}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, defaultServiceApiBaseUrl, cfg.Flux.WebhookReceiverUrl)
}

func TestModule_DefaultAndValidateConfiguration_WithWebhookReceiverUrlConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{
		Flux: &agentcfg.FluxCF{
			WebhookReceiverUrl: "https://example.com",
		},
	}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", cfg.Flux.WebhookReceiverUrl)
}

func TestModule_FailedToGetCRD(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(nil, errors.New("expected error during test"))

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	assert.ErrorContains(t, err, "unable to get CRD")
	assert.False(t, ok)
}

func TestModule_CRDNotFound(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(nil, kubeerrors.NewNotFound(testSupportedCrd.GroupResource(), testSupportedCrd.GroupResource().String()))

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_CRDMultipleVersionsNoSupport(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).
		Return(&v1.CustomResourceDefinition{
			Spec: v1.CustomResourceDefinitionSpec{
				Versions: []v1.CustomResourceDefinitionVersion{
					{
						Name:   "v1beta1",
						Served: true,
					},
					{
						Name:   "v1beta2",
						Served: true,
					},
				},
			},
			Status: v1.CustomResourceDefinitionStatus{
				Conditions: []v1.CustomResourceDefinitionCondition{
					{
						Type:   apiextensionsv1api.Established,
						Status: apiextensionsv1api.ConditionTrue,
					},
				},
			},
		}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_CRDMultipleVersionsSupportNotServed(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).
		Return(&v1.CustomResourceDefinition{
			Spec: v1.CustomResourceDefinitionSpec{
				Versions: []v1.CustomResourceDefinitionVersion{
					{
						Name:   "v1beta1",
						Served: true,
					},
					{
						Name:   "v1beta2",
						Served: true,
					},
					{
						Name:   testSupportedCrdVersion,
						Served: false,
					},
				},
			},
			Status: v1.CustomResourceDefinitionStatus{
				Conditions: []v1.CustomResourceDefinitionCondition{
					{
						Type:   apiextensionsv1api.Established,
						Status: apiextensionsv1api.ConditionTrue,
					},
				},
			},
		}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_CRDMultipleVersionsSupportedAndServedAndEstablished(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)
	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).
		Return(&v1.CustomResourceDefinition{
			Spec: v1.CustomResourceDefinitionSpec{
				Versions: []v1.CustomResourceDefinitionVersion{
					{
						Name:   "v1beta1",
						Served: true,
					},
					{
						Name:   "v1beta2",
						Served: true,
					},
					{
						Name:   testSupportedCrdVersion,
						Served: true,
					},
				},
			},
			Status: v1.CustomResourceDefinitionStatus{
				Conditions: []v1.CustomResourceDefinitionCondition{
					{
						Type:   apiextensionsv1api.Established,
						Status: apiextensionsv1api.ConditionTrue,
					},
				},
			},
		}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestModule_CRDNotEstablished(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(&v1.CustomResourceDefinition{
		Spec: v1.CustomResourceDefinitionSpec{
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:   testSupportedCrdVersion,
					Served: true,
				},
			},
		},
	}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_CRDNotEstablishedBecauseOfWrongCondition(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(&v1.CustomResourceDefinition{
		Spec: v1.CustomResourceDefinitionSpec{
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:   testSupportedCrdVersion,
					Served: true,
				},
			},
		},
		Status: v1.CustomResourceDefinitionStatus{
			Conditions: []v1.CustomResourceDefinitionCondition{
				{
					Type:   apiextensionsv1api.NamesAccepted,
					Status: apiextensionsv1api.ConditionTrue,
				},
			},
		},
	}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_CRDNotEstablishedBecauseOfWrongConditionStatus(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(&v1.CustomResourceDefinition{
		Spec: v1.CustomResourceDefinitionSpec{
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:   testSupportedCrdVersion,
					Served: true,
				},
			},
		},
		Status: v1.CustomResourceDefinitionStatus{
			Conditions: []v1.CustomResourceDefinitionCondition{
				{
					Type:   apiextensionsv1api.Established,
					Status: apiextensionsv1api.ConditionFalse,
				},
			},
		},
	}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestModule_SuccessfullyEstablishedCRD(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockApiExtClient := mock_k8s.NewMockApiextensionsV1Interface(ctrl)
	mockCRDInterface := mock_k8s.NewMockCustomResourceDefinitionInterface(ctrl)

	// setup mock expectations
	mockApiExtClient.EXPECT().CustomResourceDefinitions().Return(mockCRDInterface)
	mockCRDInterface.EXPECT().Get(gomock.Any(), testSupportedCrdResourceName, gomock.Any()).Return(&v1.CustomResourceDefinition{
		Spec: v1.CustomResourceDefinitionSpec{
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:   testSupportedCrdVersion,
					Served: true,
				},
			},
		},
		Status: v1.CustomResourceDefinitionStatus{
			Conditions: []v1.CustomResourceDefinitionCondition{
				{
					Type:   apiextensionsv1api.Established,
					Status: apiextensionsv1api.ConditionTrue,
				},
			},
		},
	}, nil)

	// WHEN
	ok, err := checkCRDExistsAndEstablished(context.Background(), mockApiExtClient, testSupportedCrd)

	// THEN
	require.NoError(t, err)
	assert.True(t, ok)
}
