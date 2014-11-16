package sabnzbd

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type Sabnzbd struct {
	mu    sync.RWMutex
	https bool
	addr  string
	path  string
	auth  authenticator
}

func New(options ...Option) (s *Sabnzbd, err error) {
	s = &Sabnzbd{
		addr: "localhost:8080",
		path: "api",
		auth: &noneAuth{},
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Sabnzbd) SetOptions(options ...Option) (err error) {
	for _, option := range options {
		if err := option(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sabnzbd) useHttps() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.https = true
}

func (s *Sabnzbd) useHttp() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.https = false
}

func (s *Sabnzbd) setAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.addr = addr
	return nil
}

func (s *Sabnzbd) setPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.path = path
	return nil
}

func (s *Sabnzbd) setAuth(a authenticator) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auth = a
	return nil
}

type sabnzbdURL struct {
	*url.URL
	v    url.Values
	auth authenticator
}

func (s *Sabnzbd) url() *sabnzbdURL {
	s.mu.RLock()
	defer s.mu.RUnlock()
	su := &sabnzbdURL{
		URL: &url.URL{
			Scheme: "http",
			Host:   s.addr,
			Path:   s.path,
		},
		auth: s.auth,
	}
	if s.https {
		su.Scheme = "https"
	}
	su.v = su.URL.Query()
	return su
}

func (su *sabnzbdURL) SetJsonOutput() {
	su.v.Set("output", "json")
}

func (su *sabnzbdURL) SetMode(mode string) {
	su.v.Set("mode", mode)
}

func (su *sabnzbdURL) Authenticate() {
	su.auth.Authenticate(su.v)
}

func (su *sabnzbdURL) String() string {
	su.RawQuery = su.v.Encode()
	return su.URL.String()
}

func (su *sabnzbdURL) CallJSON(r interface{}) error {
	resp, err := http.Get(su.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(r); err != nil {
		return err
	}
	if err, ok := r.(error); ok {
		return apiStringError(err.Error())
	}

	return nil
}

func (su *sabnzbdURL) CallJSONMultipart(reader io.Reader, contentType string, r interface{}) error {
	resp, err := http.Post(su.String(), contentType, reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(r); err != nil {
		return err
	}
	if err, ok := r.(error); ok {
		return apiStringError(err.Error())
	}

	return nil
}
