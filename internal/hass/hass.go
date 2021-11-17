package hass

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const apiPrefix = "/api/services"

// Doer represents an http client that can "Do" a request
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Session is the access and credentials for the API
type Session struct {
	host   string
	token  string
	client Doer
}

// NewSession returns a new *Session to be used to interface with the
// Home Assistant system.
func NewSession(host, token string) *Session {
	return &Session{
		host:  host,
		token: token,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *Session) TurnOn(body LightState) error {
	return s.post(apiPrefix+"/light/turn_on", body)
}

func (a *Session) get(path string, v interface{}) error {
	req, err := http.NewRequest("GET", a.host+path, nil)

	if err != nil {
		return err
	}

	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}

	success := false
	for i := 0; i < 3; i++ {
		func() {
			var resp *http.Response
			resp, err = a.client.Do(req)
			if err != nil {
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				err = errors.New("hass: status not OK: " + resp.Status)
				return
			}

			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(v)
			success = true
		}()

		if success {
			break
		}
	}

	return err
}

func (a *Session) post(path string, v interface{}) error {
	var req *http.Request

	if v != nil {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		req, err = http.NewRequest("POST", a.host+path, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		var err error
		req, err = http.NewRequest("POST", a.host+path, nil)
		if err != nil {
			return err
		}
	}

	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}

	var err error
	success := false
	for i := 0; i < 3; i++ {
		func() {
			var resp *http.Response
			resp, err = a.client.Do(req)
			if err != nil {
				return
			}

			defer resp.Body.Close()

			success = true
		}()

		if success {
			break
		}
	}

	return err
}
