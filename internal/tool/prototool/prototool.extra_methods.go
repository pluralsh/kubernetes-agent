package prototool

import (
	"net/http"
	"net/url"

	"github.com/pluralsh/kuberentes-agent/internal/tool/httpz"
)

func (x *HttpRequest) HttpHeader() http.Header {
	return ValuesMapToHttpHeader(x.GetHeader())
}

func (x *HttpRequest) UrlQuery() url.Values {
	return ValuesMapToUrlValues(x.GetQuery())
}

func (x *HttpRequest) IsUpgrade() bool {
	return x.GetHeader()[httpz.UpgradeHeader] != nil
}

func (x *HttpResponse) HttpHeader() http.Header {
	return ValuesMapToHttpHeader(x.GetHeader())
}

func ValuesMapToHttpHeader(from map[string]*Values) http.Header {
	res := make(http.Header, len(from))
	for key, val := range from {
		res[key] = val.GetValue()
	}
	return res
}

func HttpHeaderToValuesMap(from http.Header) map[string]*Values {
	res := make(map[string]*Values, len(from))
	for key, val := range from {
		res[key] = &Values{
			Value: val,
		}
	}
	return res
}

func UrlValuesToValuesMap(from url.Values) map[string]*Values {
	res := make(map[string]*Values, len(from))
	for key, val := range from {
		res[key] = &Values{
			Value: val,
		}
	}
	return res
}

func ValuesMapToUrlValues(from map[string]*Values) url.Values {
	query := make(url.Values, len(from))
	for key, val := range from {
		query[key] = val.GetValue()
	}
	return query
}
