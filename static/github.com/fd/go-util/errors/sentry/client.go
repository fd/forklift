package sentry

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	c_AUTH_HEADER       = "X-Sentry-Auth"
	c_ERROR_HEADER      = "X-Sentry-Error"
	c_USER_AGENT_HEADER = "User-Agent"
	c_USER_AGENT        = "fd-go-raven/1.0"
	c_AUTH_HEADER_FMT   = "Sentry sentry_version=4, sentry_client=%s, sentry_timestamp=%d, sentry_key=%s, sentry_secret=%s"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			ResponseHeaderTimeout: 3 * time.Second,
		},
	}
}

func set_auth_header(req *http.Request) {
	req.Header.Set(c_USER_AGENT_HEADER, c_USER_AGENT)
	req.Header.Set(c_AUTH_HEADER, fmt.Sprintf(
		c_AUTH_HEADER_FMT,
		c_USER_AGENT,
		time.Now().Unix(),
		public_key,
		secret_key,
	))
}

func call(packet *Packet) error {
	data, err := packet.to_json()
	if err != nil {
		return err
	}

	// fmt.Printf("packet=%s\n", data)

	inner := func() error {
		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(data))
		if err != nil {
			return err
		}

		set_auth_header(req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		if resp.StatusCode != 200 {
			return &apiError{resp}
		}

		return nil
	}

	for i := 1; i <= 5; i++ {
		err = inner()
		if err == nil {
			return nil
		}

		time.Sleep(time.Duration(i) * 50 * time.Millisecond)
	}

	return err
}

type apiError struct {
	resp *http.Response
}

func (err *apiError) Error() string {
	return fmt.Sprintf("sentry: (status=%d) %s", err.resp.StatusCode, err.resp.Header.Get(c_ERROR_HEADER))
}
