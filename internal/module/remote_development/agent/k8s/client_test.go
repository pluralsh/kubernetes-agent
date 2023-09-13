package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_k8s"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
)

func setupK8sClient(t *testing.T) *K8sClient {
	ctrl := gomock.NewController(t)

	mockFactory := mock_k8s.NewMockFactory(ctrl)
	mockFactory.EXPECT().KubernetesClientSet().AnyTimes()
	mockFactory.EXPECT().DynamicClient().AnyTimes()
	mockFactory.EXPECT().ToRESTMapper().AnyTimes()
	mockFactory.EXPECT().ToDiscoveryClient().AnyTimes()
	mockFactory.EXPECT().ToRESTConfig().AnyTimes()
	logger := zaptest.NewLogger(t)
	k8sClient, err := New(logger, mockFactory)
	if err != nil {
		t.Fatalf("Error during setup: %v", err)
	}

	return k8sClient
}

func Test_groupObjectsByInventory(t *testing.T) {
	k8sClient := setupK8sClient(t)

	inv1Name := "inv1"
	inv1 := getMockInventory(t, inv1Name)
	inv2Name := "inv2"
	inv2 := getMockInventory(t, inv2Name)
	secret1 := getMockSecret(t, "secret1", inv1Name)
	secret2 := getMockSecret(t, "secret2", inv1Name)
	secret3 := getMockSecret(t, "secret3", inv2Name)
	secretNoInv := getMockSecret(t, "secretNoInv", "")

	tests := []struct {
		name           string
		input          []*unstructured.Unstructured
		expectedOutput map[string]*applierInfo
		expectedErr    error
	}{
		{
			name:  "one-inventory-with-one-resource",
			input: []*unstructured.Unstructured{inv1, secret1},
			expectedOutput: map[string]*applierInfo{
				inv1Name: {
					invInfo: inv1,
					objects: []*unstructured.Unstructured{secret1},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "one-inventory-with-one-resource-with-inventory-not-as-first-object",
			input: []*unstructured.Unstructured{secret1, inv1},
			expectedOutput: map[string]*applierInfo{
				inv1Name: {
					invInfo: inv1,
					objects: []*unstructured.Unstructured{secret1},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "one-inventory-with-multiple-resources",
			input: []*unstructured.Unstructured{inv1, secret1, secret2},
			expectedOutput: map[string]*applierInfo{
				inv1Name: {
					invInfo: inv1,
					objects: []*unstructured.Unstructured{secret1, secret2},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "multiple-inventories-with-multiple-resources",
			input: []*unstructured.Unstructured{inv1, secret1, secret2, inv2, secret3},
			expectedOutput: map[string]*applierInfo{
				inv1Name: {
					invInfo: inv1,
					objects: []*unstructured.Unstructured{secret1, secret2},
				},
				inv2Name: {
					invInfo: inv2,
					objects: []*unstructured.Unstructured{secret3},
				},
			},
			expectedErr: nil,
		},
		{
			name:           "no-inventory-provided",
			input:          []*unstructured.Unstructured{secret1, secret2},
			expectedOutput: nil,
			expectedErr:    noInventoryFoundErr,
		},
		{
			name:           "no-owning-inventory-found-for-objects",
			input:          []*unstructured.Unstructured{inv1, secretNoInv},
			expectedOutput: nil,
			expectedErr:    noOwningInventoryFoundErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, actualErr := k8sClient.groupObjectsByInventory(tc.input)
			assert.Equal(t, tc.expectedOutput, actualOutput)
			assert.ErrorIs(t, tc.expectedErr, actualErr)
		})
	}
}

func getMockInventory(t *testing.T, name string) *unstructured.Unstructured {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				common.InventoryLabel: name,
			},
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&cm)
	if err != nil {
		t.Fatalf("unable to type cast inventory configmap: %v", err)
	}
	return &unstructured.Unstructured{Object: obj}
}

func getMockSecret(t *testing.T, name, owningInventory string) *unstructured.Unstructured {
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "default",
			Annotations: map[string]string{},
		},
	}
	if owningInventory != "" {
		secret.Annotations[inventory.OwningInventoryKey] = owningInventory
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&secret)

	if err != nil {
		t.Fatalf("unable to type cast secret: %v", err)
	}
	return &unstructured.Unstructured{Object: obj}
}
