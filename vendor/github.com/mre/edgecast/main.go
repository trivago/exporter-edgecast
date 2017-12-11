package edgecast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// APIEndpoint holds the Edgecast API endpoint URL
	// Meaning of the wildcards:
	// AccountID
	// Platform (see below)
	// Method endpoint (see Metric* below)
	APIEndpoint = "https://api.edgecast.com/v2/realtimestats/customers/%s/media/%d/%s"

	// DefaultRequestRetries defines the default number of http requests before giving up
	DefaultRequestRetries = 1

	// DefaultRequestTimeout defines the default request timeout in seconds
	DefaultRequestTimeout = 5

	// MethodBandwidth is the endpoint for Edgecast bandwidth
	MethodBandwidth = "bandwidth"
	// MethodConnections is the endpoint for Edgecast connections
	MethodConnections = "connections"
	// MethodCachestatus is the endpoint for the Edgecast cache status
	MethodCachestatus = "cachestatus"
	// MethodStatuscodes is the endpoint for the Edgecast status codes
	MethodStatuscodes = "statuscode"
)

// Media types (also known as "platform")
// Specify, which Edgecast platforms to monitor.
// They are identified by an integer value and passed to the API url.
// The following types are available:
// flash, http_large, http_small, adn
//
// Unfortunately the stats aren't more fine grained than this. If you have
// more than one 'service' using the platform(s), you'll get them added together.
const (
	MediaTypeFlash = 2
	MediaTypeLarge = 3
	MediaTypeSmall = 8
	MediaTypeADN   = 14
)

// Edgecast client for Go
type Edgecast struct {
	AccountID string
	BaseURL   string
	Token     string
	Retries   int
	Timeout   int
}

// NewEdgecastClient creates a new Edgecast client
func NewEdgecastClient(accountID, token string) *Edgecast {
	return &Edgecast{
		AccountID: accountID,
		BaseURL:   APIEndpoint,
		Token:     token,
		Retries:   DefaultRequestRetries,
		Timeout:   DefaultRequestTimeout,
	}
}

// SetRetries sets the number of consecutive query replies until giving up
func (e *Edgecast) SetRetries(retries int) *Edgecast {
	e.Retries = retries
	return e
}

// SetTimeout sets the request timeout
func (e *Edgecast) SetTimeout(timeout int) *Edgecast {
	e.Timeout = timeout
	return e
}

// addHeaders adds some special Edgecast headers that are required for querying the API
func (e *Edgecast) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", "TOK:"+e.Token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

// Bandwidth returns the current bandwidth usage
func (e Edgecast) Bandwidth(platform int) (*BandwidthData, error) {
	body, err := e.request(platform, MethodBandwidth)
	if err != nil {
		return nil, err
	}
	var data RawEdgecastResult
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &BandwidthData{Bps: data.Result, Platform: platform}, nil
}

// Connections returns the current bandwidth usage
func (e Edgecast) Connections(platform int) (*ConnectionData, error) {
	body, err := e.request(platform, MethodConnections)
	if err != nil {
		return nil, err
	}
	var data RawEdgecastResult
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &ConnectionData{Connections: data.Result, Platform: platform}, nil
}

// CacheStatus returns the current cache status usage
func (e Edgecast) CacheStatus(platform int) (*CacheStatusData, error) {
	body, err := e.request(platform, MethodCachestatus)
	if err != nil {
		return nil, err
	}
	var data CacheStatusData
	err = json.Unmarshal(body, &data)
	return &data, err
}

// StatusCodes returns the current HTTP status codes
func (e Edgecast) StatusCodes(platform int) (*StatusCodeData, error) {
	body, err := e.request(platform, MethodStatuscodes)
	if err != nil {
		return nil, err
	}
	var data StatusCodeData
	err = json.Unmarshal(body, &data)
	return &data, err
}

// fullURL creates a queryable URL for the current method
func (e Edgecast) fullURL(platform int, method string) string {
	return fmt.Sprintf(e.BaseURL, e.AccountID, platform, method)
}

// request runs an API request using the given parameters and returns the raw request body or an error
func (e Edgecast) request(platform int, method string) ([]byte, error) {
	url := e.fullURL(platform, method)

	client := http.Client{
		Timeout: time.Duration(time.Duration(e.Timeout) * time.Second),
	}
	var err error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	e.addHeaders(req)

	for i := 0; i < e.Retries; i++ {

		resp, err := client.Do(req)

		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		return body, nil
	}
	return nil, err
}
