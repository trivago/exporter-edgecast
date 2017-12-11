package edgecast

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestClient(t *testing.T) {
	accountID := "ABCD"
	mediaType := MediaTypeSmall
	token := "1234-5678"

	bodyBW, _ := ioutil.ReadFile("fixtures/bandwidth.json")
	bodyConn, _ := ioutil.ReadFile("fixtures/connections.json")
	bodyCachestatus, _ := ioutil.ReadFile("fixtures/cachestatus.json")
	bodyStatuscodes, _ := ioutil.ReadFile("fixtures/statuscode.json")
	apiEndpoint := "https://api.edgecast.com/v2/realtimestats/customers/%s/media/%d/%s"

	bandwidthURL := fmt.Sprintf(apiEndpoint, accountID, mediaType, "bandwidth")
	connectionsURL := fmt.Sprintf(apiEndpoint, accountID, mediaType, "connections")
	cachestatusURL := fmt.Sprintf(apiEndpoint, accountID, mediaType, "cachestatus")
	statuscodesURL := fmt.Sprintf(apiEndpoint, accountID, mediaType, "statuscode")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", bandwidthURL, httpmock.NewStringResponder(200, string(bodyBW)))
	httpmock.RegisterResponder("GET", connectionsURL, httpmock.NewStringResponder(200, string(bodyConn)))
	httpmock.RegisterResponder("GET", cachestatusURL, httpmock.NewStringResponder(200, string(bodyCachestatus)))
	httpmock.RegisterResponder("GET", statuscodesURL, httpmock.NewStringResponder(200, string(bodyStatuscodes)))

	client := NewEdgecastClient(accountID, token)

	// Test bandwidth request
	dataBW, err := client.Bandwidth(mediaType)

	if err != nil {
		t.Errorf("Error while calling API %s", err)
	}

	expectedBW := BandwidthData{
		Bps:      42.5,
		Platform: mediaType,
	}

	if !reflect.DeepEqual(dataBW, &expectedBW) {
		t.Errorf("Expected %s (type %#v), got %s (type %#v)", expectedBW, expectedBW, dataBW, dataBW)
	}

	// Test connection request
	dataConn, err := client.Connections(mediaType)

	if err != nil {
		t.Errorf("Error while calling API %s", err)
	}

	expectedConn := ConnectionData{
		Connections: 7270.72352,
		Platform:    mediaType,
	}

	if !reflect.DeepEqual(dataConn, &expectedConn) {
		t.Errorf("Expected %s (type %#v), got %s (type %#v)", expectedConn, expectedConn, dataConn, dataConn)
	}

	// Test cachestatus request
	dataCachestatus, err := client.CacheStatus(mediaType)

	if err != nil {
		t.Errorf("Error while calling API %s", err)
	}

	expectedCachestatus := CacheStatusData{
		{
			CacheStatus: "TCP_HIT",
			Connections: 6740,
		},
		{
			CacheStatus: "TCP_EXPIRED_HIT",
			Connections: 1,
		},
		{
			CacheStatus: "TCP_MISS",
			Connections: 235,
		},
		{
			CacheStatus: "TCP_EXPIRED_MISS",
			Connections: 0,
		},
		{
			CacheStatus: "TCP_CLIENT_REFRESH_MISS",
			Connections: 123,
		},
		{
			CacheStatus: "NONE",
			Connections: 42,
		},
		{
			CacheStatus: "CONFIG_NOCACHE",
			Connections: 12,
		},
		{
			CacheStatus: "UNCACHEABLE",
			Connections: 4,
		},
	}

	if !reflect.DeepEqual(dataCachestatus, &expectedCachestatus) {
		t.Errorf("Expected\n%#v\nGot:\n%#v)", expectedCachestatus, dataCachestatus)
	}

	// Test statuscode request
	dataStatuscode, err := client.StatusCodes(mediaType)

	if err != nil {
		t.Errorf("Error while calling API %s", err)
	}

	expectedStatuscode := StatusCodeData{
		{
			Connections: 6671,
			StatusCode:  "2xx",
		},
		{
			Connections: 131,
			StatusCode:  "304",
		},
		{
			Connections: 0,
			StatusCode:  "3xx",
		},
		{
			Connections: 88,
			StatusCode:  "403",
		},
		{
			Connections: 2,
			StatusCode:  "404",
		},
		{
			Connections: 9,
			StatusCode:  "4xx",
		},
		{
			Connections: 0,
			StatusCode:  "503",
		},
		{
			Connections: 0,
			StatusCode:  "504",
		},
		{
			Connections: 0,
			StatusCode:  "5xx",
		},
		{
			Connections: 0,
			StatusCode:  "other",
		},
	}

	if !reflect.DeepEqual(dataStatuscode, &expectedStatuscode) {
		t.Errorf("Expected\n%#v\nGot:\n%#v)", expectedStatuscode, dataStatuscode)
	}
}

func TestFluentInterface(t *testing.T) {
	accountID := "ABCD"
	token := "1234-5678"

	client := NewEdgecastClient(accountID, token).SetRetries(3).SetTimeout(100)

	if client.Timeout != 100 {
		t.Errorf("Setting timeout did not work")
	}
	if client.Retries != 3 {
		t.Errorf("Setting retries did not work")
	}
}
