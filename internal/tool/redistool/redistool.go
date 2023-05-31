package redistool

import (
	"crypto/tls"
	"net"

	"github.com/redis/rueidis"
)

func MultiFirstError(resp []rueidis.RedisResult) error {
	for _, r := range resp {
		if err := r.Error(); err != nil {
			return err
		}
	}
	return nil
}

// UnixDialer can be used as DialFn in rueidis.ClientOption.
func UnixDialer(addr string, dialer *net.Dialer, tlsConfig *tls.Config) (net.Conn, error) {
	return net.DialUnix("unix", nil, &net.UnixAddr{
		Name: addr,
		Net:  "unix",
	})
}
