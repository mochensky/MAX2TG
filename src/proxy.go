package src

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func BuildSOCKS5Dialer(cfg *ProxyConfig) (proxy.Dialer, error) {
	if cfg == nil {
		return proxy.Direct, nil
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var auth *proxy.Auth
	if cfg.Username != "" || cfg.Password != "" {
		auth = &proxy.Auth{
			User:     cfg.Username,
			Password: cfg.Password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer for %s: %w", addr, err)
	}

	Logf("SOCKS5 proxy configured: %s", addr)
	return dialer, nil
}

func buildRemoteDNSDialer(proxyCfg *ProxyConfig) (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	addr := fmt.Sprintf("%s:%d", proxyCfg.Host, proxyCfg.Port)

	var auth *proxy.Auth
	if proxyCfg.Username != "" || proxyCfg.Password != "" {
		auth = &proxy.Auth{
			User:     proxyCfg.Username,
			Password: proxyCfg.Password,
		}
	}

	contextDialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer for %s: %w", addr, err)
	}

	cd, ok := contextDialer.(proxy.ContextDialer)
	if !ok {
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			return contextDialer.Dial(network, addr)
		}, nil
	}

	return cd.DialContext, nil
}

func BuildHTTPClientWithProxy(proxyCfg *ProxyConfig, timeout time.Duration) (*http.Client, error) {
	if proxyCfg == nil {
		return &http.Client{Timeout: timeout}, nil
	}

	dialContext, err := buildRemoteDNSDialer(proxyCfg)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: dialContext,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}, nil
}

func BuildHTTPTransportWithProxy(cfg *ProxyConfig) (*http.Transport, error) {
	if cfg == nil {
		return &http.Transport{}, nil
	}

	dialContext, err := buildRemoteDNSDialer(cfg)
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		DialContext: dialContext,
	}, nil
}

func buildProxyURL(cfg *ProxyConfig) *url.URL {
	if cfg == nil {
		return nil
	}
	u := &url.URL{
		Scheme: "socks5",
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	if cfg.Username != "" {
		u.User = url.UserPassword(cfg.Username, cfg.Password)
	}
	return u
}

func GetMaxProxy(cfg *Config) *ProxyConfig {
	for i := range cfg.Proxies {
		if cfg.Proxies[i].ForMax {
			return &cfg.Proxies[i]
		}
	}
	return nil
}

func GetTelegramProxy(cfg *Config) *ProxyConfig {
	for i := range cfg.Proxies {
		if cfg.Proxies[i].ForTelegram {
			return &cfg.Proxies[i]
		}
	}
	return nil
}
