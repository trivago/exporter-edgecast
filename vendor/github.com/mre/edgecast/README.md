# edgecast

A golang client for the [Edgecast CDN API](https://www.programmableweb.com/api/edgecast-cdn).  

## Installation

```go
// Go vendor experiment (recommended)
gvt fetch "github.com/mre/edgecast"

// or globally (not recommended)
go get github.com/mre/edgecast
```

## Usage:

```go
import github.com/mre/edgecast

client := edgecast.NewEdgecastClient("AccountID", "Token")
data, err := client.Bandwidth()

// Example output:
// BandwidthData {
//    Bps:      42.5,
//    Platform: 2,
// }

```

## Methods


```go
client.Bandwidth()    // Return current bandwidth usage
client.Connections()  // Return number of active CDN connections
client.CacheStatus()  // Return cache hits and misses
client.StatusCodes()  // Return sum of HTTP status codes by category (404, 5xx,...)
```

## Fluent interface

You can also set additional parameters using a fluent interface:

```go
client := edgecast.NewEdgecastClient(config.AccountID, config.Token).
    SetRetries(3). // Setup HTTP request retries (e.g. for flaky connections)
    SetTimeout(5)  // Set request timeout per HTTP request (in seconds
```
