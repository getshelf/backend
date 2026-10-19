package webpage

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"
)

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: nil,
		DialContext: dialPublic,
		TLSHandshakeTimeout: 5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		MaxResponseHeaderBytes: 64 << 10,
		IdleConnTimeout: 30 * time.Second,
		MaxIdleConns: 20,
		MaxConnsPerHost: 4,
	}

	return &http.Client{
		Transport: transport,
		Timeout: 12 * time.Second,
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			if len(via) > 5 {
				return fmt.Errorf("too many redirects")
			}

			return validateURL(req.URL)
		},
	}
}

func validateURL(value *url.URL) error {
	if value.Scheme != "http" && value.Scheme != "https" {
		return fmt.Errorf("page URL must use HTTP or HTTPS")
	}

	if value.Hostname() == "" {
		return fmt.Errorf("page URL must have a hostname")
	}

	if value.User != nil {
		return fmt.Errorf("page URL must not contain credentials")
	}

	if port := value.Port(); port != "" && port != "80" && port != "443" {
		return fmt.Errorf("page URL uses an unsupported port")
	}

	return nil
}

func dialPublic(
	ctx context.Context,
	network string,
	address string,
) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)

	if err != nil {
		return nil, fmt.Errorf("parse destination: %w", err)
	}

	addresses, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)

	if err != nil {
		return nil, fmt.Errorf("resolve destination: %w", err)
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("destination has no IP addresses")
	}

	for _, address := range addresses {
		if !allowedIP(address) {
			return nil, fmt.Errorf("destination resolves to a blocked address")
		}
	}

	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	var lastErr error
	for _, address := range addresses {
		conn, err := dialer.DialContext(
			ctx,
			network,
			net.JoinHostPort(address.Unmap().String(), port),
		)

		if err == nil {
			return conn, nil
		}

		lastErr = err
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("connect to destination: %w", lastErr)
}

var globalIPv6 = netip.MustParsePrefix("2000::/3")

var blockedRanges = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("3fff::/20"),
}

func allowedIP(address netip.Addr) bool {
	address = address.Unmap()

	if address.Zone() != "" ||
	!address.IsGlobalUnicast() ||
	address.IsPrivate() ||
	address.IsLoopback() ||
	address.IsLinkLocalUnicast() {
		return false
	}

	if address.Is6() && !globalIPv6.Contains(address) {
		return false
	}

	for _, prefix := range blockedRanges {
		if prefix.Contains(address) {
			return false
		}
	}

	return true
}
