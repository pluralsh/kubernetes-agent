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
		_, err := ParseReferenceDiscovery(bytes.NewReader(data))
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

	d, err := ParseReferenceDiscovery(buf)
	require.NoError(t, err)
	require.Equal(t, []string{"capability"}, d.Caps)
	require.Equal(t, []Reference{{Oid: oid1, Name: "HEAD"}}, d.Refs)
}

func TestMultipleRefsAndCapsParse(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00first second")
	gittest.WritePktlineString(t, buf, oid2+" refs/heads/master")
	gittest.WritePktlineFlush(t, buf)

	d, err := ParseReferenceDiscovery(buf)
	require.NoError(t, err)
	require.Equal(t, []string{"first", "second"}, d.Caps)
	require.Equal(t, []Reference{{Oid: oid1, Name: "HEAD"}, {Oid: oid2, Name: "refs/heads/master"}}, d.Refs)
}

func TestInvalidHeaderFails(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=invalid\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")
	gittest.WritePktlineFlush(t, buf)

	_, err := ParseReferenceDiscovery(buf)
	require.Error(t, err)
}

func TestMissingRefsFail(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineFlush(t, buf)

	_, err := ParseReferenceDiscovery(buf)
	require.Error(t, err)
}

func TestInvalidRefFail(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")
	gittest.WritePktlineString(t, buf, oid2)
	gittest.WritePktlineFlush(t, buf)

	_, err := ParseReferenceDiscovery(buf)
	require.Error(t, err)
}

func TestMissingTrailingFlushFails(t *testing.T) {
	buf := &bytes.Buffer{}
	gittest.WritePktlineString(t, buf, "# service=git-upload-pack\n")
	gittest.WritePktlineFlush(t, buf)
	gittest.WritePktlineString(t, buf, oid1+" HEAD\x00caps")

	d := ReferenceDiscovery{}
	require.Error(t, d.Parse(buf))
}
