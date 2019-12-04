moov-io/watchman
===

[![GoDoc](https://godoc.org/github.com/moov-io/watchman?status.svg)](https://godoc.org/github.com/moov-io/watchman)
[![Build Status](https://travis-ci.com/moov-io/watchman.svg?branch=master)](https://travis-ci.com/moov-io/watchman)
[![Coverage Status](https://codecov.io/gh/moov-io/watchman/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/watchman)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/watchman)](https://goreportcard.com/report/github.com/moov-io/watchman)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/watchman/master/LICENSE)

Moov Watchman is an HTTP API and Go library to download, parse and offer search functions over numerous trade sanction lists from the United States, European Union governments, agencies, and non profits for complying with regional laws. Also included is a web UI and async webhook notification service to initiate processes on remote systems.

Lists included in search are:

- [Office of Foreign Asset Control](https://www.treasury.gov/about/organizational-structure/offices/Pages/Office-of-Foreign-Assets-Control.aspx) (OFAC)
   - The US Treasury department offers a [search page for OFAC records](https://sanctionssearch.ofac.treas.gov/)
- [BIS Denied Person's List](https://bis.data.commerce.gov/dataset/Denied-Persons-List-with-Denied-US-Export-Privileg/xwtd-wd7a/data) (DPL)

All United States or European Union companies are required to comply with various regulations and sanction lists (such as the US Patriot Act requiring compliance with the BIS Denied Person's List). Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

To get started using watchman download [the latest release](https://github.com/moov-io/watchman/releases/latest) or our [Docker image](https://hub.docker.com/r/moov/watchman/tags). We also have a [demo instance](https://moov.io/watchman/) as part of Moov's demo environment.

```
# Run as a binary
$ wget https://github.com/moov-io/watchman/releases/download/v0.12.0/watchman-darwin-amd64
$ chmod +x watchman-darwin-amd64
$ ./watchman-darwin-amd64
ts=2019-02-05T00:03:31.9583844Z caller=main.go:42 startup="Starting watchman server version v0.12.0"
...

# Run as a Docker image
$ docker run -p 8084:8084 -p 9094:9094 -it moov/watchman:latest
ts=2019-02-05T00:03:31.9583844Z caller=main.go:42 startup="Starting watchman server version v0.12.0"
...

# Perform a basic search
$ curl -s localhost:8084/search?name=...
{
  "SDNs": [
    {
      "entityID": "...",
      "sdnName": "...",
      "sdnType": "...",
      "program": "...",
      "title": "...",
      "callSign": "...",
      "vesselType": "...",
      "tonnage": "...",
      "grossRegisteredTonnage": "...",
      "vesselFlag": "...",
      "vesselOwner": "...",
      "remarks": "..."
    }
  ],
  "altNames": null,
  "addresses": null,
  "deniedPersons": null
}
```

We offer [hosted api docs as part of Moov's tools](https://api.moov.io/#tag/Sanctions) and an [OpenAPI specification](https://github.com/moov-io/watchman/blob/master/openapi.yaml) for use with generated clients.

Docs: [docs.moov.io](https://docs.moov.io/watchman/) | [api docs](https://api.moov.io/apps/watchman/)

### Web UI

Moov Sanction Search ships with a web interface for easier access searching the records. Our Docker image hosts the UI by default, but you can build and run it locally as well.

![](docs/images/webui.png)

```
$ make
...
CGO_ENABLED=1 go build -o ./bin/server github.com/moov-io/watchman/cmd/server
...
npm run build
...
Success!

$ go run ./cmd/server/ # Load http://localhost:8084 in a web browser
```

### Configuration

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `DATA_REFRESH_INTERVAL` | Interval for data redownload and reparse. `off` disables this refreshing. | 12h |
| `INITIAL_DATA_DIRECTORY` | Directory filepath with initial files to use instead of downloading. Periodic downloads will replace the initial files. | Empty |
| `WEBHOOK_BATCH_SIZE` | How many watches to read from database per batch of async searches. | 100 |
| `LOG_FORMAT` | Format for logging lines to be written as. | Options: `json`, `plain` - Default: `plain` |
| `BASE_PATH` | HTTP path to serve API and web UI from. | `/` |
| `HTTP_BIND_ADDRESS` | Address to bind HTTP server on. This overrides the command-line flag `-http.addr`. | Default: `:8084` |
| `HTTP_ADMIN_BIND_ADDRESS` | Address to bind admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9094` |
| `HTTPS_CERT_FILE` | Filepath containing a certificate (or intermediate chain) to be served by the HTTP server. Requires all traffic be over secure HTTP. | Empty |
| `HTTPS_KEY_FILE`  | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`. | Empty |
| `DATABASE_TYPE` | Which database option to use (Options: `sqlite`, `mysql`) | Default: `sqlite` |
| `WEB_ROOT` | Directory to serve web UI from | Default: `webui/` |

### List Configurations

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `OFAC_DOWNLOAD_TEMPLATE` | HTTP address for downloading raw OFAC files. | (OFAC website) |
| `DPL_DOWNLOAD_TEMPLATE` | HTTP address for downloading the DPL | (BIS website) |

#### Storage

Based on `DATABASE_TYPE` the following environment variables will be read to configure connections for a specific database.

##### MySQL

- `MYSQL_ADDRESS`: TCP address for connecting to the mysql server. (example: `tcp(hostname:3306)`)
- `MYSQL_DATABASE`: Name of database to connect into.
- `MYSQL_PASSWORD`: Password of user account for authentication.
- `MYSQL_USER`: Username used for authentication,

Refer to the mysql driver documentation for [connection parameters](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

- `MYSQL_TIMEOUT`: Timeout parameter specified on (DSN) data source name. (Default: `30s`)

##### SQLite

- `SQLITE_DB_PATH`: Local filepath location for the paygate SQLite database. (Default: `watchman.db`)

Refer to the sqlite driver documentation for [connection parameters](https://github.com/mattn/go-sqlite3#connection-string).

### Features

- Download OFAC, BIS Denied Persons List (DPL), and various other data sources on startup
  - Admin endpoint to [manually refresh OFAC and DPL data](docs/runbook.md#force-data-refresh)
- Index data for searches
- Async searches and notifications (webhooks)
- Manual overrides to mark a `Company` or `Customer` as `unsafe` (blocked) or `exception` (never blocked).
- Library for OFAC and BIS DPL data to download and parse their custom files

#### Webhook Notifications

When SancionSearch sends a [webhook](https://en.wikipedia.org/wiki/Webhook) to your application the body will contain a JSON representation of the [Company](https://godoc.org/github.com/moov-io/watchman/client#Company) or [Customer](https://godoc.org/github.com/moov-io/watchman/client#Customer) model as the body to a POST request. You can see an [example in Go](examples/webhook/webhook.go).

An `Authorization` header will also be sent with the `authToken` provided when setting up the watch. Clients should verify this token to ensure authenticated communicated.

Webhook notifications are ran after the OFAC data is successfully refreshed, which is determined by the `DATA_REFRESH_INTERVAL` environmental variable.

##### Watching a specific Customer or Company by ID

Moov Sanction Search supports sending a webhook periodically when a specific [Company](https://api.moov.io/#operation/addCompanyWatch) or [Customer](https://api.moov.io/#operation/addCustomerWatch) is to be watched. This is designed to update another system about an OFAC entry's sanction status.

##### Watching a customer or company name

Moov Sanction Search supports sending a webhook periodically with a free-form name of a [Company](https://api.moov.io/#operation/addCompanyNameWatch) or [Customer](https://api.moov.io/#operation/addCustomerNameWatch). This allows external applications to be notified when an entity matching that name is added to the OFAC list. The match percentage will be included in the JSON payload.

##### Prometheus Metrics

- `http_response_duration_seconds`: A Histogram of HTTP response timings
- `last_data_refresh_success`: Unix timestamp of when data was last refreshed successfully
- `match_percentages` A Histogram which holds the match percentages with a label (`type`) of searches
   - `type`: Can be address, q, remarksID, name, altName
- `mysql_connections`: How many MySQL connections and what status they're in.
- `sqlite_connections`: How many sqlite connections and what status they're in.

## Generating a Client

We use [openapi-generator](https://github.com/OpenAPITools/openapi-generator) from the [OpenAPI team](https://swagger.io/specification/) to generate API clients for popular programming languages from the API specification. To generate the Go client run `make client` from Watchman's root directory.

## Getting Help

We maintain a [runbook for common issues](docs/runbook.md) and configuration options. Also, if you've encountered a security issue please contact us at [`security@moov.io`](mailto:security@moov.io).

 channel | info
 ------- | -------
 [Project Documentation](https://docs.moov.io/) | Our project documentation available online.
 Google Group [moov-users](https://groups.google.com/forum/#!forum/moov-users)| The Moov users Google group is for contributors other people contributing to the Moov project. You can join them without a google account by sending an email to [moov-users+subscribe@googlegroups.com](mailto:moov-users+subscribe@googlegroups.com). After receiving the join-request message, you can simply reply to that to confirm the subscription.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io) | If you are able to reproduce an problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](http://moov-io.slack.com/) | Join our slack channel to have an interactive discussion about the development of the project. [Request an invite to the slack channel](https://join.slack.com/t/moov-io/shared_invite/enQtNDE5NzIwNTYxODEwLTRkYTcyZDI5ZTlkZWRjMzlhMWVhMGZlOTZiOTk4MmM3MmRhZDY4OTJiMDVjOTE2MGEyNWYzYzY1MGMyMThiZjg)

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](https://github.com/moov-io/ach/blob/master/CODE_OF_CONDUCT.md) to get started!

Note: This project uses Go Modules, which requires Go 1.11 or higher, but we ship the vendor directory in our repository.

## Links

- [Sanctions Search Page](https://sanctionssearch.ofac.treas.gov/)
- [Subscribe for OFAC updates](https://service.govdelivery.com/accounts/USTREAS/subscriber/new)
- [When should I call the OFAC Hotline?](https://www.treasury.gov/resource-center/faqs/Sanctions/Pages/directions.aspx)
- [BIS Denied Persons List with Denied US Export Privileges](https://bis.data.commerce.gov/dataset/Denied-Persons-List-with-Denied-US-Export-Privileg/xwtd-wd7a/data)

## License

Apache License 2.0 See [LICENSE](LICENSE) for details.
