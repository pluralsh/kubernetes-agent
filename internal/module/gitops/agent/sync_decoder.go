package agent

import (
	"bytes"
	"errors"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/rpc"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/cli-utils/pkg/manifestreader"
)

type syncDecoder struct {
	restClientGetter resource.RESTClientGetter
	defaultNamespace string
}

func (d *syncDecoder) Decode(sources []rpc.ObjectSource) ([]*unstructured.Unstructured, error) {
	if len(sources) == 0 {
		return nil, nil
	}
	// 1. Parse in local mode to retrieve objects.
	builder := resource.NewBuilder(d.restClientGetter).
		ContinueOnError().
		Flatten().
		Unstructured().
		Local()
	for _, source := range sources {
		builder.Stream(bytes.NewReader(source.Data), source.Name)
	}
	result := builder.Do()
	var objs []*unstructured.Unstructured
	err := result.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		objs = append(objs, info.Object.(*unstructured.Unstructured))
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 2. Process parsed objects - set namespace to the default one if missing
	restMapper, err := d.restClientGetter.ToRESTMapper()
	if err != nil {
		return nil, err
	}
	triedToResetMapper := false
	for {
		err = manifestreader.SetNamespaces(restMapper, objs, d.defaultNamespace, false)
		if err != nil {
			if triedToResetMapper || !isUnknownTypesError(err) || !isResettableRESTMapper(restMapper) {
				return nil, err
			}
			// Haven't tried to reset the mapper AND it's an UnknownTypesError AND the mapper is resettable.
			// Reset the mapper and try one more time.
			triedToResetMapper = true
			meta.MaybeResetRESTMapper(restMapper)
			continue
		}
		return objs, nil
	}
}

func isUnknownTypesError(err error) bool {
	var typeErr *manifestreader.UnknownTypesError
	return errors.As(err, &typeErr)
}

func isResettableRESTMapper(mapper meta.RESTMapper) bool {
	_, ok := mapper.(meta.ResettableRESTMapper)
	return ok
}
