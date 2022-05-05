package httpz

import (
	"net/http"
	"net/textproto"
	"strings"
)

// These headers must be in their canonical form. Only add headers used in production code, don't bother with tests.
const (
	ConnectionHeader         = "Connection" // https://datatracker.ietf.org/doc/html/rfc7230#section-6.1
	ProxyConnectionHeader    = "Proxy-Connection"
	KeepAliveHeader          = "Keep-Alive"
	HostHeader               = "Host"
	ProxyAuthenticateHeader  = "Proxy-Authenticate"
	ProxyAuthorizationHeader = "Proxy-Authorization"
	TeHeader                 = "Te"      // canonicalized version of "TE"
	TrailerHeader            = "Trailer" // not Trailers as per rfc2616; See errata https://www.rfc-editor.org/errata_search.php?eid=4522
	TransferEncodingHeader   = "Transfer-Encoding"
	UpgradeHeader            = "Upgrade"
	UserAgentHeader          = "User-Agent"
	AuthorizationHeader      = "Authorization"
	ContentTypeHeader        = "Content-Type"
	AcceptHeader             = "Accept"
	ServerHeader             = "Server"
	ViaHeader                = "Via" // https://datatracker.ietf.org/doc/html/rfc7230#section-5.7.1
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
