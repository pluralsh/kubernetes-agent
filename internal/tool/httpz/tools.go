package httpz

import (
	"mime"
	"net/http"
	"net/textproto"
	"strings"
)

// These headers must be in their canonical form. Only add headers used in production code, don't bother with tests.
const (
	ConnectionHeader         = "Connection" // https://datatracker.ietf.org/doc/html/rfc9110#section-7.6.1
	ProxyConnectionHeader    = "Proxy-Connection"
	KeepAliveHeader          = "Keep-Alive"
	HostHeader               = "Host"
	ProxyAuthenticateHeader  = "Proxy-Authenticate"
	ProxyAuthorizationHeader = "Proxy-Authorization"
	TeHeader                 = "Te"      // canonicalized version of "TE"
	TrailerHeader            = "Trailer" // not Trailers as per rfc2616; See errata https://www.rfc-editor.org/errata_search.php?eid=4522
	TransferEncodingHeader   = "Transfer-Encoding"
	UpgradeHeader            = "Upgrade" // https://datatracker.ietf.org/doc/html/rfc9110#section-7.8
	UserAgentHeader          = "User-Agent"
	AuthorizationHeader      = "Authorization" // https://datatracker.ietf.org/doc/html/rfc9110#section-11.6.2
	ContentTypeHeader        = "Content-Type"  // https://datatracker.ietf.org/doc/html/rfc9110#section-8.3
	AcceptHeader             = "Accept"        // https://datatracker.ietf.org/doc/html/rfc9110#section-12.5.1
	ServerHeader             = "Server"        // https://datatracker.ietf.org/doc/html/rfc9110#section-10.2.4
	ViaHeader                = "Via"           // https://datatracker.ietf.org/doc/html/rfc9110#section-7.6.3
)

// RemoveConnectionHeaders removes hop-by-hop headers listed in the "Connection" header of h.
// See https://datatracker.ietf.org/doc/html/rfc7230#section-6.1
func RemoveConnectionHeaders(h http.Header) {
	for _, f := range h[ConnectionHeader] {
		for _, sf := range strings.Split(f, ",") {
			if sf = textproto.TrimString(sf); sf != "" {
				// must use .Del() because connection options are case-insensitive and are likely in lower case, not in canonical case
				h.Del(sf)
			}
		}
	}
}

func IsContentType(actual string, expected ...string) bool {
	parsed, _, err := mime.ParseMediaType(actual)
	if err != nil {
		return false
	}
	for _, e := range expected {
		if e == parsed {
			return true
		}
	}
	return false
}
