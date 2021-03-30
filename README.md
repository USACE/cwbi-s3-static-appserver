# CWBI S3 Static Appserver

A Virtual-Host style webserver written in Go using the [Echo Web Framework](https://echo.labstack.com/).

A `sidecar` container syncs files from an S3 bucket to a container volume `data` mounted at `/data`.  The `data` volume is shared between the `sidecar` and `appserver` containers, providing the files to be served by `appserver`

The entire stack can be brought-up locally using `docker-compose`.

[MINIO](https://min.io/) - which has a AWS S3 compatible API - is used locally in place of AWS S3.  Files and folders placed in the `./apps` directory of this repository will be loaded into MINIO on `docker-compose up`.

Used to host the following applications in `development`:

| Application                                          | Hosted Locally At ...    |
| ---------------------------------------------------- | ------------------------ |
| Landing Page                                         | `localhost:8000`         |
| [Cumulus](https://github.com/USACE/cumulus-ui)       | `cumulus.localhost:8000` |
| [MIDAS](https://github.com/USACE/instrumentation-ui) | `midas.localhost:8000`   |
| [Access to Water](https://github.com/USACE/water-ui) | `water.localhost:8000`   |
