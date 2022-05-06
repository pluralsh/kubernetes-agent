package agent

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
	"sigs.k8s.io/yaml"
)

const (
	defaultNs = "ns-default"

	yamlNamespace = `apiVersion: v1
kind: Namespace
metadata:
  name: gitlab-agent
`
	yamlNamespaceNs = `apiVersion: v1
kind: Namespace
metadata:
  name: gitlab-agent
  namespace: ns-ns
`
	yamlConfigMapNoNs = `apiVersion: v1
kind: ConfigMap
metadata:
  name: map-no-ns
`
	yamlConfigMapNs = `apiVersion: v1
kind: ConfigMap
metadata:
  name: map-ns
  namespace: ns-map
`
	yamlCRD = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: issuers.cert-manager.io
spec:
  group: cert-manager.io
  names:
    kind: Issuer
    listKind: IssuerList
    plural: issuers
    singular: issuer
  scope: Namespaced
  versions:
    - name: v1
      subresources:
        status: {}
`
	yamlCRv1 = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
spec:
  ca:
    secretName: ca-key-pair
`

	yamlCRv2 = `apiVersion: cert-manager.io/v2
kind: Issuer
metadata:
  name: ca-issuer
spec:
  ca:
    secretName: ca-key-pair
`
)

func TestSyncDecoder_HappyPath(t *testing.T) {
	tests := []struct {
		name            string
		sources         []rpc.ObjectSource
		expectedErr     string
		expectedObjects []*unstructured.Unstructured
	}{
		{
			name: "namespaced with namespace",
			sources: []rpc.ObjectSource{
				{
					Name: "config-map-ns",
					Data: []byte(yamlConfigMapNs),
				},
			},
			expectedObjects: []*unstructured.Unstructured{
				yaml2unstructured(t, yamlConfigMapNs, ""),
			},
		},
		{
			name: "namespaced without namespace",
			sources: []rpc.ObjectSource{
				{
					Name: "config-map-no-ns",
					Data: []byte(yamlConfigMapNoNs),
				},
			},
			expectedObjects: []*unstructured.Unstructured{
				yaml2unstructured(t, yamlConfigMapNoNs, defaultNs),
			},
		},
		{
			name: "cluster-scoped",
			sources: []rpc.ObjectSource{
				{
					Name: "ns",
					Data: []byte(yamlNamespace),
				},
			},
			expectedObjects: []*unstructured.Unstructured{
				yaml2unstructured(t, yamlNamespace, ""),
			},
		},
		{
			name: "cluster-scoped with namespace",
			sources: []rpc.ObjectSource{
				{
					Name: "ns",
					Data: []byte(yamlNamespaceNs),
				},
			},
			expectedErr: `resource is cluster-scoped but has a non-empty namespace "ns-ns"`,
		},
		{
			name: "CRD",
			sources: []rpc.ObjectSource{
				{
					Name: "crd",
					Data: []byte(yamlCRD),
				},
			},
			expectedObjects: []*unstructured.Unstructured{
				yaml2unstructured(t, yamlCRD, ""),
			},
		},
		{
			name: "CRD and CRv1",
			sources: []rpc.ObjectSource{
				{
					Name: "crv1",
					Data: []byte(yamlCRv1),
				},
				{
					Name: "crd",
					Data: []byte(yamlCRD),
				},
			},
			expectedObjects: []*unstructured.Unstructured{
				yaml2unstructured(t, yamlCRv1, defaultNs),
				yaml2unstructured(t, yamlCRD, ""),
			},
		},
		{
			name: "CRv1",
			sources: []rpc.ObjectSource{
				{
					Name: "crv1",
					Data: []byte(yamlCRv1),
				},
			},
			expectedErr: "unknown resource types: cert-manager.io/v1/Issuer",
		},
		{
			name: "CRD and CRv2",
			sources: []rpc.ObjectSource{
				{
					Name: "crv2",
					Data: []byte(yamlCRv2),
				},
				{
					Name: "crd",
					Data: []byte(yamlCRD),
				},
			},
			expectedErr: "unknown resource types: cert-manager.io/v2/Issuer",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			factory := cmdtesting.NewTestFactory()
			mapper, err := factory.ToRESTMapper()
			require.NoError(t, err)
			crdGV := schema.GroupVersion{Group: "apiextensions.k8s.io", Version: "v1"}
			crdMapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{crdGV})
			crdMapper.AddSpecific(crdGV.WithKind("CustomResourceDefinition"),
				crdGV.WithResource("customresourcedefinitions"),
				crdGV.WithResource("customresourcedefinition"), meta.RESTScopeRoot)
			mapper = meta.MultiRESTMapper([]meta.RESTMapper{mapper, crdMapper})

			d := &syncDecoder{
				restMapper:       mapper,
				restClientGetter: factory,
				defaultNamespace: defaultNs,
			}
			res, err := d.Decode(tc.sources)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Empty(t, cmp.Diff(res, tc.expectedObjects))
			}
		})
	}
}

func yaml2unstructured(t *testing.T, yml, setNamespace string) *unstructured.Unstructured {
	var o unstructured.Unstructured
	err := yaml.Unmarshal([]byte(yml), &o.Object)
	require.NoError(t, err)
	if setNamespace != "" {
		err = unstructured.SetNestedField(o.Object, setNamespace, "metadata", "namespace")
		require.NoError(t, err)
	}
	return &o
}
