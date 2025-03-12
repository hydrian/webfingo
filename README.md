# Webfingo

Webfingo is a simple HTTP server that serves WebFinger requests for Keycloak. It
is intended for integrating Keycloak with Tailscale.

## Basic Functionality

1. Acting as a webserver, receive Webfinger request from Tailscale, e.g. at
   `https://webfingo.example.com/.well-known/webfinger?resource=acct:john@example.com`.
2. Connect to Keycloak's Postgres database, and look up the user and their realm
3. Respond to the HTTP request with a normal WebFinger response

## Usage

```
make build
./bin/webfingo --config [your config file]
```

See `./config/config-example.json`. Note that the Keycloak DB config and the
hostname of your Keycloak endpoint will be needed.

## Requirements

- A Keycloak instance with a Postgres database
- Go 1.24 or later
