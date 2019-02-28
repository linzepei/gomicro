# Homeweb Service

This is the Homeweb service

Generated with

```
micro new go-1/homeweb --namespace=go.micro --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.web.homeweb
- Type: web
- Alias: homeweb

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./homeweb-web
```

Build a docker image
```
make docker
```