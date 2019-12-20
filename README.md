# Money service
The experimental service for playing with go-kit

## Requirements
* git (for getting the source code)
* make
* docker
* docker-compose

## Build

The application and tests are built and run inside docker containers.

Clone sources
```bash
git clone https://github.com/denperov/money-service.git
```

Build
```bash
make build
```

Run
```bash
make run
```

Test API
```bash
make test
```

Direct access to database
```
docker run -it --rm --network=host postgres:12.1-alpine sh -c "PGPASSWORD=accounts psql -U accounts -h `hostname` -d accounts"
```

## API
[API description](docs/accounts/api.md)

## Packages

Package "github.com/caarlos0/env" is used to get parameters from environment variables, since this is a more common way for services running inside docker containers. It can be replaced with any other convenient package.

## Graceful shutdown

Service supports graceful shutdown with sequential shutdown of components

## Remark

To work with money in golang code, a slightly unusual data type is selected, but in the current example they are only passed to the API, and this is better than a float.
