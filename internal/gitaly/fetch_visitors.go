package gitaly

import (
	"fmt"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DupBehavior byte

const (
	// DupError means "return error on duplicate file".
	DupError DupBehavior = 1
	// DupSkip means "skip duplicate files".
	DupSkip DupBehavior = 2
)

type ChunkingFetchVisitor struct {
	FetchVisitor
	maxChunkSize int
}

func NewChunkingFetchVisitor(delegate FetchVisitor, maxChunkSize int) *ChunkingFetchVisitor {
	return &ChunkingFetchVisitor{
		FetchVisitor: delegate,
		maxChunkSize: maxChunkSize,
	}
}

func (v ChunkingFetchVisitor) StreamChunk(path []byte, data []byte) (bool /* done? */, error) {
	for {
		bytesToSend := minInt(len(data), v.maxChunkSize)
		done, err := v.FetchVisitor.StreamChunk(path, data[:bytesToSend])
		if err != nil || done {
			return done, err
		}
		data = data[bytesToSend:]
		if len(data) == 0 {
			break
		}
	}
	return false, nil
}

type MaxNumberOfFilesError struct {
	MaxNumberOfFiles uint32
}

func (e *MaxNumberOfFilesError) Error() string {
	return fmt.Sprintf("maximum number of files limit reached: %d", e.MaxNumberOfFiles)
}

type EntryCountLimitingFetchVisitor struct {
	FetchVisitor
	maxNumberOfFiles uint32
	FilesVisited     uint32
	FilesSent        uint32
}

func NewEntryCountLimitingFetchVisitor(delegate FetchVisitor, maxNumberOfFiles uint32) *EntryCountLimitingFetchVisitor {
	return &EntryCountLimitingFetchVisitor{
		FetchVisitor:     delegate,
		maxNumberOfFiles: maxNumberOfFiles,
	}
}

func (v *EntryCountLimitingFetchVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	if v.FilesVisited == v.maxNumberOfFiles {
		return false, 0, &MaxNumberOfFilesError{
			MaxNumberOfFiles: v.maxNumberOfFiles,
		}
	}
	v.FilesVisited++
	return v.FetchVisitor.Entry(entry)
}

func (v *EntryCountLimitingFetchVisitor) EntryDone(entry *gitalypb.TreeEntry, err error) {
	v.FetchVisitor.EntryDone(entry, err)
	if err != nil {
		return
	}
	v.FilesSent++
}

type TotalSizeLimitingFetchVisitor struct {
	FetchVisitor
	RemainingTotalFileSize int64
}

func NewTotalSizeLimitingFetchVisitor(delegate FetchVisitor, maxTotalFileSize int64) *TotalSizeLimitingFetchVisitor {
	return &TotalSizeLimitingFetchVisitor{
		FetchVisitor:           delegate,
		RemainingTotalFileSize: maxTotalFileSize,
	}
}

func (v *TotalSizeLimitingFetchVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	shouldDownload, maxSize, err := v.FetchVisitor.Entry(entry)
	if err != nil || !shouldDownload {
		return false, 0, err
	}
	return true, minInt64(v.RemainingTotalFileSize, maxSize), nil
}

func (v *TotalSizeLimitingFetchVisitor) StreamChunk(path []byte, data []byte) (bool /* done? */, error) {
	v.RemainingTotalFileSize -= int64(len(data))
	if v.RemainingTotalFileSize < 0 {
		// This should never happen because we told Gitaly the maximum file size that we'd like to get.
		// i.e. we should have gotten an error from Gitaly if file is bigger than the limit.
		return false, status.Error(codes.Internal, "unexpected negative remaining total file size")
	}
	return v.FetchVisitor.StreamChunk(path, data)
}

type HiddenDirFilteringFetchVisitor struct {
	FetchVisitor
}

func NewHiddenDirFilteringFetchVisitor(delegate FetchVisitor) *HiddenDirFilteringFetchVisitor {
	return &HiddenDirFilteringFetchVisitor{
		FetchVisitor: delegate,
	}
}

func (v HiddenDirFilteringFetchVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	if isHiddenDir(string(entry.Path)) {
		return false, 0, nil
	}
	return v.FetchVisitor.Entry(entry)
}

type GlobMatchFailedError struct {
	Cause error
	Glob  string
}

func (e *GlobMatchFailedError) Error() string {
	return fmt.Sprintf("glob %s match failed: %v", e.Glob, e.Cause)
}

func (e *GlobMatchFailedError) Unwrap() error {
	return e.Cause
}

type GlobFilteringFetchVisitor struct {
	FetchVisitor
	Glob string
}

func NewGlobFilteringFetchVisitor(delegate FetchVisitor, glob string) *GlobFilteringFetchVisitor {
	return &GlobFilteringFetchVisitor{
		FetchVisitor: delegate,
		Glob:         glob,
	}
}

func (v GlobFilteringFetchVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	shouldDownload, err := doublestar.Match(v.Glob, string(entry.Path))
	if err != nil {
		return false, 0, &GlobMatchFailedError{
			Cause: err,
			Glob:  v.Glob,
		}
	}
	if !shouldDownload {
		return false, 0, nil
	}
	return v.FetchVisitor.Entry(entry)
}

type DuplicatePathFoundError struct {
	Path string
}

func (e *DuplicatePathFoundError) Error() string {
	return fmt.Sprintf("path visited more than once: %s", e.Path)
}

type DuplicatePathDetectingVisitor struct {
	FetchVisitor
	visited     map[string]struct{}
	DupBehavior DupBehavior
}

func NewDuplicateFileDetectingVisitor(delegate FetchVisitor, dupBehavior DupBehavior) DuplicatePathDetectingVisitor {
	return DuplicatePathDetectingVisitor{
		FetchVisitor: delegate,
		visited:      map[string]struct{}{},
		DupBehavior:  dupBehavior,
	}
}

func (v DuplicatePathDetectingVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	p := string(entry.Path)
	if _, visited := v.visited[p]; visited {
		switch v.DupBehavior {
		case DupError:
			return false, 0, &DuplicatePathFoundError{
				Path: p,
			}
		case DupSkip:
			return false, 0, nil
		default:
			panic(fmt.Errorf("unknown dup behavior: %d", v.DupBehavior))
		}
	}
	v.visited[p] = struct{}{}
	return v.FetchVisitor.Entry(entry)
}

// isHiddenDir checks if a file is in a directory, which name starts with a dot.
func isHiddenDir(filename string) bool {
	dir := path.Dir(filename)
	if dir == "." { // root directory special case
		return false
	}
	parts := strings.Split(dir, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
