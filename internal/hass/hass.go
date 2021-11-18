package hass

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mishamyrt/ambihass/internal/log"
)

const apiPrefix = "/api"
const apiServices = apiPrefix + "/services"

// Session is the access and credentials for the API
type Session struct {
	host   string
	token  string
	client *http.Client
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
	return s.post(apiServices+"/light/turn_on", body)
}

func (s *Session) CheckAPI() error {
	return s.get(apiServices)
}

func (s *Session) get(path string) (err error) {
	req, err := s.createRequest("GET", path, nil)
	s.execute(req)
	return nil
}

func (s *Session) post(path string, v interface{}) (err error) {
	var data []byte
	if v != nil {
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	req, err := s.createRequest("POST", path, data)
	s.execute(req)
	return nil
}

func (s *Session) createRequest(method string, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, s.host+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)
	return req, nil
}

func (s *Session) execute(req *http.Request) (success bool) {
	for i := 0; i < 3; i++ {
		func() {
			resp, err := s.client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = errors.New("hass: status not OK: " + resp.Status)
				log.Debug("Request error: ", err)
				return
			}
			success = true
		}()

		if success {
			break
		}
	}
	return
}
