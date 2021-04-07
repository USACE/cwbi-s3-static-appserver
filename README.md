# CWBI S3 Static Appserver

## Description

A Virtual-Host style webserver written in Go using the [Echo Web Framework](https://echo.labstack.com/).

A `sidecar` container syncs files from an S3 bucket to a container volume `data` mounted at `/data`. The `data` volume is shared between the `sidecar` and `appserver` containers, providing the files to be served by `appserver`

The entire stack can be brought-up locally using `docker-compose`.

[MINIO](https://min.io/) - which has a AWS S3 compatible API - is used locally in place of AWS S3. Files and folders placed in the `./apps` directory of this repository will be loaded into MINIO on `docker-compose up`.

## Used to host the following applications in `development`

| Application                                                        | Hosted Locally At ...    |
| ------------------------------------------------------------------ | ------------------------ |
| [Home](https://github.com/USACE/cwbi-application-development-docs) | `home.localhost:8080`    |
| [Cumulus](https://github.com/USACE/cumulus-ui)                     | `cumulus.localhost:8080` |
| [MIDAS](https://github.com/USACE/instrumentation-ui)               | `midas.localhost:8080`   |
| [Access to Water](https://github.com/USACE/water-ui)               | `water.localhost:8080`   |

## Environment Variables for Container `appserver`

- `APPSSERVER_DOMAIN`: Top level domain where all apps will be served. Individual apps are served as a sub-domain to this top level domain. Default value for development
- `APPSERVER_SUBDOMAIN_PREFIX`: Prepend all hostnames with this prefix. This is most easily explained using a practical example:
- `APPSERVER_RELOAD_INTERVAL`: How often to reload the appserver.  A reload causes the appserver to serve any new content made available by `appserver-sidecar` since `appserver` was last started. Example Values: 5m, 10s, 1h (values supported by bash `sleep` command)

  For application `water`, setting `APPSERVER_DOMAIN=rsgis.dev` and `APPSERVER_SUBDOMAIN_PREFIX=develop-` serves the application `water` at `develop-water.rsgis.dev`. Setting `APPSERVER_DOMAIN=rsgis.dev` and omitting `APPSERVER_SUBDOMAIN_PREFIX` (i.e. `""`) serves the application `water` at `water.rsgis.dev`. Practically, this is used to serve `development` and `stable` versions of apps using different s3 buckets, but the same container image.

## Environment Variables for Container `appserver-sidecar`

- `AWS_REGION`: Credentials to serve files from S3 Bucket (may or may not be explicitly required depending on your deployment details)
- `AWS_ACCESS_KEY_ID`: Credentials to serve files from S3 Bucket (may or may not be explicitly required depending on your deployment details)
- `AWS_SECRET_ACCESS_KEY`: Credentials to serve files from S3 Bucket (may or may not be explicitly required depending on your deployment details)
- `S3_BUCKET`: Name of S3 Bucket where files to be served exist
- `S3_ENDPOINT_URL`: Endpoint for S3 API Access - used for local development and testing against `minio` (set to http://minio:9000 for development)
