package stats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitaly/copied/gittest"
)

const (
	oid1 = "78fb81a02b03f0013360292ec5106763af32c287"
	oid2 = "0f6394307cd7d4909be96a0c818d8094a4cb0e5b"
)

func BenchmarkMultipleRefsAndCapsParse(b *testing.B) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(b, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(b, buf)
	gittest.WritePktlineString(b, buf, oid1+" HEAD\x00first second")
	gittest.WritePktlineString(b, buf, oid2+" refs/heads/master")
	gittest.WritePktlineFlush(b, buf)
	data := buf.Bytes()
	r := bytes.NewReader(data)
	b.ReportAllocs()
	b.ResetTimer() // don't take the stuff above into account
	for i := 0; i < b.N; i++ {
		err := ParseReferenceDiscovery(r, func(ref Reference) bool { return false })
		if err != nil {
			b.Fatal(err)
		}
		r.Reset(data)
	}
}

func TestSingleRefParses(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00capability")
	gittest.WritePktlineFlush(t, buf)

	var refs []Reference
	err := ParseReferenceDiscovery(buf, accumulateRefs(&refs))
	require.NoError(t, err)
	require.Equal(t, []Reference{{Oid: []byte(oid1), Name: []byte("HEAD")}}, refs)
}

func TestMultipleRefsAndCapsParse(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00first second")
	gittest.WritePktlineString(t, buf, oid2+" refs/heads/master")
	gittest.WritePktlineFlush(t, buf)

	var refs []Reference
	err := ParseReferenceDiscovery(buf, accumulateRefs(&refs))
	require.NoError(t, err)
	require.Equal(t, []Reference{{Oid: []byte(oid1), Name: []byte("HEAD")}, {Oid: []byte(oid2), Name: []byte("refs/heads/master")}}, refs)
}

func TestInvalidHeaderFails(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=invalid\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")
	gittest.WritePktlineFlush(t, buf)

	err := ParseReferenceDiscovery(buf, func(ref Reference) bool { return false })
	require.Error(t, err)
}

func TestMissingRefsFail(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineFlush(t, buf)

	err := ParseReferenceDiscovery(buf, func(ref Reference) bool { return false })
	require.Error(t, err)
}

func TestInvalidRefFail(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")
	gittest.WritePktlineString(t, buf, oid2)
	gittest.WritePktlineFlush(t, buf)

	err := ParseReferenceDiscovery(buf, func(ref Reference) bool { return false })
	require.Error(t, err)
}

func TestMissingTrailingFlushFails(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")

	err := ParseReferenceDiscovery(buf, func(ref Reference) bool { return false })
	require.Error(t, err)
}

func cloneSlice(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func accumulateRefs(refs *[]Reference) ReferenceCb {
	return func(ref Reference) bool {
		*refs = append(*refs, Reference{
			Oid:  cloneSlice(ref.Oid),
			Name: cloneSlice(ref.Name),
		})
		return false
	}
}
