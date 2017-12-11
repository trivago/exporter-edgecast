package edgecast

// BandwidthData holds that data of a request
// to the edgecast bandwidth API
type BandwidthData struct {
	Bps      float64
	Platform int
}

// ConnectionData holds that data of a request
// to the edgecast connections API
type ConnectionData struct {
	Connections float64
	Platform    int
}

// CacheStatusData represends all fields returned from
// a request to the /cachestatus endpoint
type CacheStatusData []struct {
	CacheStatus string `json:"CacheStatus"`
	Connections int64  `json:"Connections"`
}

// StatusCodeData represends all fields returned from
// a request to the /statuscode endpoint
type StatusCodeData []struct {
	Connections int64  `json:"Connections"`
	StatusCode  string `json:"StatusCode"`
}

// RawEdgecastResult represends a raw JSON response object from the API
type RawEdgecastResult struct {
	Result float64
}
