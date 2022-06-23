package stats

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitaly/copied/pktline"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/memz"
)

// Reference as used by the reference discovery protocol.
type Reference struct {
	// Oid is the object ID the reference points to
	Oid []byte
	// Name of the reference. The name will be suffixed with ^{} in case
	// the reference is the peeled commit.
	Name []byte
}

// ReferenceCb is a callback that consumes parsed references.
// WARNING: It must not hold onto the byte slices as the backing array is reused! Make copies if needed.
// Returns true if reference parsing should stop.
type ReferenceCb func(Reference) bool

type referenceDiscoveryState int

const (
	referenceDiscoveryExpectService referenceDiscoveryState = iota
	referenceDiscoveryExpectFlush
	referenceDiscoveryExpectRefWithCaps
	referenceDiscoveryExpectRef
	referenceDiscoveryExpectEnd
)

// ParseReferenceDiscovery parses a client's reference discovery stream and
// calls cb with references. It returns an error in case
// it couldn't make sense of the client's request.
//
// Expected protocol:
// - "# service=git-upload-pack\n"
// - FLUSH
// - "<OID> <ref>\x00<capabilities>\n"
// - "<OID> <ref>\n"
// - ...
// - FLUSH
func ParseReferenceDiscovery(body io.Reader, cb ReferenceCb) error {
	state := referenceDiscoveryExpectService
	buf := memz.Get64k()
	defer memz.Put64k(buf)
	scanner := pktline.NewScanner(body, buf)

	for scanner.Scan() {
		pkt := scanner.Bytes()
		data := bytes.TrimSuffix(pktline.Data(pkt), []byte{'\n'})

		switch state {
		case referenceDiscoveryExpectService:
			if !bytes.Equal(data, []byte("# service=git-upload-pack")) {
				return fmt.Errorf("unexpected header %q", data)
			}

			state = referenceDiscoveryExpectFlush
		case referenceDiscoveryExpectFlush:
			if !pktline.IsFlush(pkt) {
				return errors.New("missing flush after service announcement")
			}

			state = referenceDiscoveryExpectRefWithCaps
		case referenceDiscoveryExpectRefWithCaps:
			if len(data) == 0 { // no refs in an empty repo
				state = referenceDiscoveryExpectEnd
				continue
			}
			split := bytes.SplitN(data, []byte{0}, 2)
			if len(split) != 2 {
				return errors.New("invalid first reference line")
			}
			ref := bytes.SplitN(split[0], []byte{' '}, 2)
			if len(ref) != 2 {
				return errors.New("invalid reference line")
			}
			if cb(Reference{Oid: ref[0], Name: ref[1]}) {
				return nil
			}

			state = referenceDiscoveryExpectRef
		case referenceDiscoveryExpectRef:
			if pktline.IsFlush(pkt) {
				state = referenceDiscoveryExpectEnd
				continue
			}

			split := bytes.SplitN(data, []byte{' '}, 2)
			if len(split) != 2 {
				return errors.New("invalid reference line")
			}
			if cb(Reference{Oid: split[0], Name: split[1]}) {
				return nil
			}
		case referenceDiscoveryExpectEnd:
			return errors.New("received packet after flush")
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if state != referenceDiscoveryExpectEnd {
		return errors.New("discovery ended prematurely")
	}

	return nil
}
