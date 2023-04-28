package chartops

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/imdario/mergo"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/yaml"
)

const (
	noSourcePos = -1
)

type valueFromASource struct {
	sourcePos int
	values    ChartValues
}

type fetchedFromSource struct {
	// log and files are set if this is a chart.

	log   *zap.Logger
	chart *chart.Chart

	// values may be set or not.

	values []valueFromASource
}

type valuesHolder struct {
	values    ChartValues
	hasValues bool
}

type projectIdRef struct {
	projectId string

	// refs

	branch string
	tag    string
	commit string
}

func projectIdRefFromRef(projectId string, ref *agentcfg.GitRefCF) projectIdRef {
	return projectIdRef{
		projectId: projectId,
		branch:    ref.GetBranch(),
		tag:       ref.GetTag(),
		commit:    ref.GetCommit(),
	}
}

type watchInfo struct {
	req             *rpc.ObjectsToSynchronizeRequest
	valuesFileNames map[string]int // filename -> pos in valuesSources list

	// Shouldn't start with a slash because server returns relative paths. Root path is empty string.
	// Only set if chart is coming from a project.
	chartPathWithSlashAtEnd string
	isChart                 bool
}

func (w *worker) fetch(ctx context.Context, jobs chan<- job) {
	// 0. Boilerplate.
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	project2watch := make(map[projectIdRef]*watchInfo) // project id+ref -> watchInfo
	valuesFromSources := make([]valuesHolder, len(w.chartCfg.Values))
	fetched := make(chan fetchedFromSource)

	// 1. Setup chart fetching.
	switch src := w.chartCfg.Source.Source.(type) {
	case *agentcfg.ChartSourceCF_Project:
		pathWithSlash := strings.TrimSuffix(src.Project.Path, "/") + "/"
		project2watch[projectIdRefFromRef(src.Project.Id, src.Project.Ref)] = &watchInfo{
			req: &rpc.ObjectsToSynchronizeRequest{
				ProjectId: src.Project.Id,
				Ref:       rpc.NewRpcRef(src.Project.Ref),
				Paths: []*rpc.PathCF{
					{
						Path: &rpc.PathCF_Glob{
							Glob: pathWithSlash + "**", // TODO validate it's a path, not a glob?
						},
					},
				},
			},
			valuesFileNames:         map[string]int{},
			chartPathWithSlashAtEnd: strings.TrimLeft(pathWithSlash, "/"),
			isChart:                 true,
		}
	default:
		// Should never happen.
		panic(fmt.Errorf("unknown source type: %T", w.chartCfg.Source.Source))
	}

	// 2. Setup values fetching.
	for pos, valsFrom := range w.chartCfg.Values {
		switch src := valsFrom.From.(type) {
		case *agentcfg.ChartValuesCF_Inline:
			valuesFromSources[pos] = valuesHolder{values: src.Inline.AsMap(), hasValues: true}
		case *agentcfg.ChartValuesCF_File:
			key := projectIdRefFromRef(*src.File.ProjectId, src.File.Ref)
			wi := project2watch[key]
			if wi == nil {
				wi = &watchInfo{
					req: &rpc.ObjectsToSynchronizeRequest{
						ProjectId: *src.File.ProjectId,
						Ref:       rpc.NewRpcRef(src.File.Ref),
					},
					valuesFileNames: map[string]int{},
				}
				project2watch[key] = wi
			}
			wi.req.Paths = append(wi.req.Paths, &rpc.PathCF{
				Path: &rpc.PathCF_File{
					File: src.File.File,
				},
			})
			wi.valuesFileNames[src.File.File] = pos
		case *agentcfg.ChartValuesCF_Url:
			wg.StartWithContext(ctx, w.watchValuesUrlSource(src.Url, pos, fetched))
		default:
			panic(fmt.Errorf("unexpected values type: %T", valsFrom.From))
		}
	}

	// 3. Start GitLab repo watches.
	for _, wi := range project2watch {
		wg.StartWithContext(ctx, w.watchProjectSource(wi, fetched))
	}

	// 4. Handle incoming data.
	var (
		newJob      job
		nilableJobs chan<- job
		jobCancel   context.CancelFunc
	)
	defer func() {
		if jobCancel != nil {
			jobCancel()
		}
	}()
	done := ctx.Done()
handlingFetchedData:
	for {
		select {
		case <-done:
			return // nolint: govet
		case f := <-fetched:
			jobChanged := false
			for _, v := range f.values {
				oldVals := valuesFromSources[v.sourcePos]
				jobChanged = jobChanged || !oldVals.hasValues || !reflect.DeepEqual(oldVals.values, v.values)
				valuesFromSources[v.sourcePos] = valuesHolder{values: v.values, hasValues: true}
			}
			if f.chart != nil {
				newJob.log = f.log
				newJob.chart = f.chart
				jobChanged = true
			}
			if newJob.chart == nil {
				continue // haven't fetched the chart yet.
			}
			if !jobChanged {
				continue // nothing to do as neither values nor chart have changed.
			}
			var mergedValues ChartValues
			for pos, vals := range valuesFromSources {
				if !vals.hasValues {
					newJob.log.Debug("Haven't fetched this values source yet, waiting", w.logForSource(pos)...)
					continue handlingFetchedData
				}
				err := mergo.Merge(&mergedValues, vals.values, mergo.WithOverride)
				if err != nil {
					newJob.log.Error("Error merging chart values", append(w.logForSource(pos), logz.Error(err))...)
					continue handlingFetchedData
				}
			}
			if jobCancel != nil {
				jobCancel() // Cancel running/pending job ASAP
			}
			newJob.ctx, jobCancel = context.WithCancel(context.Background()) // nolint: govet
			newJob.values = mergedValues
			nilableJobs = jobs // enable select case
		case nilableJobs <- newJob:
			nilableJobs = nil // disable this case
		}
	}
}

func (w *worker) watchValuesUrlSource(src *agentcfg.ChartValuesUrlCF, pos int, fetched chan<- fetchedFromSource) func(context.Context) {
	return func(ctx context.Context) {
		_ = wait.PollImmediateInfiniteWithContext(ctx, src.PollPeriod.AsDuration(), func(ctx context.Context) (bool /* done */, error) {
			values, err := w.pollValuesUrlSource(ctx, src)
			if err != nil {
				w.log.Error("Failed to fetch values from URL", append(w.logForSource(pos), logz.Error(err))...)
				return false, nil
			}
			f := fetchedFromSource{
				values: []valueFromASource{
					{
						sourcePos: pos,
						values:    values,
					},
				},
			}
			select {
			case <-ctx.Done():
			case fetched <- f:
			}
			return false, nil
		})
	}
}

func (w *worker) pollValuesUrlSource(ctx context.Context, src *agentcfg.ChartValuesUrlCF) (retVal ChartValues, retErr error) {
	// 1. Setup request.
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, src.Url, nil)
	if err != nil {
		return nil, err
	}
	r.Header[httpz.AcceptHeader] = []string{"application/yaml, application/json"}

	// 2. Make request.
	resp, err := w.httpClient.RoundTrip(r) // nolint: bodyclose
	if err != nil {
		return nil, err
	}
	defer errz.SafeClose(resp.Body, &retErr)
	// 3. Check HTTP status code of the response.
	switch resp.StatusCode {
	case http.StatusOK:
		// handled below
	case http.StatusNoContent:
		return nil, nil // don't event look at response body
	default:
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	// 4. Check response content type.
	contentType := resp.Header.Get(httpz.ContentTypeHeader)
	if !httpz.IsContentType(contentType, "application/json", "application/yaml") {
		return nil, fmt.Errorf("unexpected %s in response: %q", httpz.ContentTypeHeader, contentType)
	}
	// 5. Read response body.
	maxSize := *src.MaxFileSize
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxSize)))
	if err != nil {
		return nil, err
	}
	if uint32(len(body)) == maxSize { // don't want to parse truncated file.
		return nil, fmt.Errorf("max file size reached: %d bytes", maxSize)
	}
	// 6. Parse response body.
	var values ChartValues
	err = yaml.Unmarshal(body, &values)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}
	return values, nil
}

func (w *worker) watchProjectSource(wi *watchInfo, fetched chan<- fetchedFromSource) func(context.Context) {
	return func(ctx context.Context) {
		var chartHash []byte
		w.objWatcher.Watch(ctx, wi.req, func(ctx context.Context, data rpc.ObjectsToSynchronizeData) {
			var f fetchedFromSource
			var err error
			f.log = w.log.With(logz.CommitId(data.CommitId))
			if wi.isChart {
				files, newHash := loadChartFromSources(wi.chartPathWithSlashAtEnd, data.Sources)
				if !bytes.Equal(newHash, chartHash) {
					// Data has changed, reload chart.
					f.chart, err = loader.LoadFiles(files)
					if err != nil {
						f.log.Error("Failed to load chart", logz.Error(err))
						return
					}
					chartHash = newHash
				}
			}
			var errPos int
			f.values, errPos, err = loadValuesFromSources(wi.valuesFileNames, data.Sources)
			if err != nil {
				f.log.Error("Failed to load values for chart", append(w.logForSource(errPos), logz.Error(err))...)
				return
			}
			select {
			case <-ctx.Done():
			case fetched <- f:
			}
		})
	}
}

func (w *worker) logForSource(pos int) []zap.Field {
	if pos == noSourcePos {
		return nil
	}
	from := w.chartCfg.Values[pos].GetFrom()
	switch v := from.(type) {
	case *agentcfg.ChartValuesCF_Inline:
		return []zap.Field{logz.Filename("inline")}
	case *agentcfg.ChartValuesCF_File:
		ref := v.File.Ref
		return []zap.Field{
			logz.ProjectId(*v.File.ProjectId),
			logz.GitRef(resolvedRef(ref.GetBranch(), ref.GetTag(), ref.GetCommit())),
			logz.Filename(v.File.File),
		}
	case *agentcfg.ChartValuesCF_Url:
		return []zap.Field{logz.Url(v.Url.Url)}
	default:
		// Shouldn't happen
		panic(fmt.Errorf("unknown type: %T", from))
	}
}

func loadChartFromSources(chartPathWithSlashAtEnd string, sources []rpc.ObjectSource) ([]*loader.BufferedFile, []byte) {
	files := make([]*loader.BufferedFile, 0, len(sources))
	for _, source := range sources {
		if !strings.HasPrefix(source.Name, chartPathWithSlashAtEnd) {
			// Not part of the chart, but a source for values.
			continue
		}
		files = append(files, &loader.BufferedFile{
			Name: source.Name[len(chartPathWithSlashAtEnd):],
			Data: source.Data,
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
	h := fnv.New128()
	for _, file := range files {
		_, _ = h.Write([]byte(file.Name))
		_, _ = h.Write([]byte{11}) // delimiter
		_, _ = h.Write(file.Data)
		_, _ = h.Write([]byte{11}) // delimiter
	}
	return files, h.Sum(nil)
}

func loadValuesFromSources(valuesFileNames map[string]int, sources []rpc.ObjectSource) ([]valueFromASource, int /* error pos */, error) {
	missingFiles := make(map[string]struct{}, len(valuesFileNames))
	for file := range valuesFileNames {
		missingFiles[file] = struct{}{}
	}
	values := make([]valueFromASource, 0, len(valuesFileNames))
	for _, source := range sources {
		if pos, ok := valuesFileNames[source.Name]; ok {
			if _, firstEncounter := missingFiles[source.Name]; !firstEncounter {
				// This should never happen, but if it does, we want to report the problem and stop the sync
				// to avoid any further surprises.
				return nil, pos, errors.New("server sent file with chart values more than once")
			}
			var vals ChartValues
			err := yaml.Unmarshal(source.Data, &vals)
			if err != nil {
				return nil, pos, fmt.Errorf("unmarshaling chart values: %w", err)
			}
			values = append(values, valueFromASource{
				sourcePos: pos,
				values:    vals,
			})
			delete(missingFiles, source.Name)
		}
	}
	if len(missingFiles) > 0 {
		return nil, noSourcePos, fmt.Errorf("server didn't send expected files for chart values: %s", missingFiles)
	}
	return values, 0, nil
}

func resolvedRef(branch, tag, commit string) string {
	switch {
	case tag != "":
		return "refs/tags/" + tag
	case branch != "":
		return "refs/heads/" + branch
	case commit != "":
		return commit
	default:
		return gitaly.DefaultBranch
	}
}
