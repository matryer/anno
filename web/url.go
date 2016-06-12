package web

import (
	neturl "net/url"
	"strings"
)

// ShortURL gets a nice human readable URL string.
// Loses information from the URL.
func ShortURL(u string, maxlength int) string {

	if maxlength < 4 {
		panic("anno: maxlength must be >4")
	}

	// ensure it has a protocol
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}
	url, err := neturl.Parse(u)
	if err != nil {
		// invalid URL; ignore
		return u
	}
	host := url.Host

	// trim www
	host = strings.TrimPrefix(host, "/")
	if strings.HasPrefix(strings.ToLower(host), "www.") {
		host = host[4:]
	}

	path := strings.TrimPrefix(url.Path, "/")

	if len(path) > 0 {
		host += "/" + path
	}

	if len(host) > maxlength {
		host = host[0:maxlength-3] + "..."
	}

	return host
}
