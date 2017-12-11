# Exporter-Edgecast

### What is this repository for?

This is a Prometheus Exporter/Collector for Edgecast CDN.

Exporter-Edgecast uses the [edgecast-client](https://github.com/mre/edgecast) created by [Matthias Endler](https://github.com/mre) to fetch metrics from the EdgeCast CDN API and then transforms and exposes them to be scraped and displayed by [Prometheus](https://prometheus.io/).

### Package Management
* This project uses **dep** as package manager
* versions are tracked in `Gopkg.lock`
* glide settings are included in `Gopkg.toml`
* get dep here: [https://github.com/golang/dep](https://github.com/golang/dep)

### Static Analysis
- ```make lint``` (uses gometalinter, downloads and installs it in case of absence)

### Build
- ```make build``` (builds for Windows or Unix, after checking ```$(OS),Windows_NT```)

### Configure
- You need to set two environment-variables to configure your Edgecast-Account:
    + EDGECAST_ACCOUNT_ID
    + EDGECAST_TOKEN
- e.g. on Linux: `export EDGECAST_TOKEN=B12AC`

### Run
- `./bin/main` (Unix) or `.\bin\main.exe` (Windows)
- via Docker:
    + build Docker image: `make docker`
    + run Docker image: `(sudo) docker run -p=<some_free_port>:80 trivago/monitoring:edgecast-v1 -e "EDGECAST_TOKEN=<your_token>" -e "EDGECAST_ACCOUNTID=<your_id>"`
        * NOTE: <some_free_port> must be the same as specified in the job-description in prometheus.yml


### View Exposed Metrics:
- via Browser on the same machine: visit [http://localhost:80/metrics](http://localhost:80/metrics)
    + via Browser on different machine: change "localhost" to endpoint address
- via existing Prometheus server installation: 
    + start new server locally using the provided configuration file:
        * `prometheus -config-file=prometheus.yml`
        * view results here: [http://localhost:9090](http://localhost:9090)
    + copy & paste job from provided `prometheus.yml` to running server's configuration to scrape the service metrics

### Exposed Metrics

See information to all the possible metrics offered by the API in the [official documentation](./docs/[Documentation]EdgeCast_Web_Services_REST_API.pdf).

#### EdgeCast Metrics
- `Edgecast_metrics_bandwidth_bps`
    + HELP:     Current amount of bandwidth usage per platform (bits per second)
    + TYPE:     GaugeValue
    + Labels:
        * platform = [http_small|http_large|adn|flash]
- `Edgecast_metrics_cachestatus`
    + HELP:     Breakdown of the cache statuses currently being returned for requests to CDN account.
    + TYPE:     GaugeValue
    + Labels:
        * platform = [http_small|http_large|adn|flash]
        * CacheStatus = [TCP_HIT|TCP_MISS|...]
- `Edgecast_metrics_connections`
    + HELP:     Total active connections per second per platform.
    + TYPE:     GaugeValue
    + Labels:
        * platform = [http_small|http_large|adn|flash]
- `Edgecast_metrics_statuscodes`
    + HELP:     Breakdown of the HTTP status codes currently being returned for requests to CDN account.
    + TYPE:     GaugeValue
    + Labels:
        * platform = [http_small|http_large|adn|flash]
        * StatusCode = [2xx|3xx|404|...]

#### Service Metrics
- `Edgecast_service_metrics_request_count`
    + HELP:     Number of requests received.
    + TYPE:     CounterValue
    + Labels:
        * method
        * error
- `Edgecast_service_metrics_request_latency_seconds`
    + HELP:     Duration of request in seconds.
    + TYPE:     GaugeValue
    + Labels:
        * method
        * error
- `Edgecast_service_metrics_request_latency_distribution_seconds`
    + HELP:     Total duration of requests in seconds.
    + TYPE:     Summary
    + Labels:
        * method
        * error

### Queried Platforms:
| MediaTypeId | Platform                     | Naming     |
|-------------|------------------------------|------------|
| 2           | Flash Media Streaming        | flash      |
| 3           | HTTP Large                   | http_large |
| 8           | HTTP Small                   | http_small |
| 14          | Application Delivery Network | adn        |

**Note**: MediaTypeId 7, 9, 15 refer to SSL-Traffic only for the platforms 3, 8, 14 respectively (Docs page 467) and are not queried yet.
