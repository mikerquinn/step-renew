# step-renew(1)

## NAME

**step-renew** — renew a certificate from step-ca using mTLS

## SYNOPSIS

**step-renew**
**-ca** *cafile*
**-cert** *certfile*
**-key** *keyfile*
**-server** *server*

## DESCRIPTION

The **step-renew** utility renews an X.509 certificate issued by a step-ca server.
It authenticates to the CA using mutual TLS (mTLS) with the existing client certificate and key, posts a renewal request to the `/1.0/renew` endpoint, and atomically overwrites the certificate file with the full chain returned by the server.

This tool is intended for automated certificate renewal in scripts, cron jobs, or systemd timers where short-lived certificates from step-ca are in use.  It was created as a minimal replacement for the `step` cli tool in constrained environments where a large executable like that is just too big (OpenWRT routers with limited flash storage and such).

The client certificate file is both the authentication credential and the target for the renewed chain (leaf certificate followed by any intermediates or CA).

## OPTIONS

The following options are mandatory:

**-ca** *path*

	Path to the PEM file containing the CA certificate (or certificate bundle) used to verify the step-ca server's TLS certificate. Required.

**-cert** *path*

	Path to the PEM-encoded client certificate. This file is read for mTLS authentication and then atomically replaced with the renewed certificate chain on success. Required.

**-key** *path*

	Path to the PEM-encoded private key corresponding to the client certificate. Required.

**-server** *host* | *url*

	Target step-ca server.

	If the argument does not begin with `https://`, the program builds the request URL as:

	    https:// + host + (":443" if no port specified) + "/1.0/renew"

	If a complete `https://...` URL is supplied, it is used verbatim (no path is appended).

	Required.

## EXAMPLES

Renew using a hostname (default port 443):

```sh
step-renew -ca /etc/ssl/ca.crt \
           -cert /etc/ssl/cert.pem \
           -key /etc/ssl/key.pem \
           -server step-ca.example.com
```
## EXIT STATUS

The utility exits with status 0 on successful renewal.
Any error (missing required flags, I/O failure, TLS handshake failure, empty response from server, etc.) results in a non-zero exit status.

## SEE ALSO

step-ca(8), step(1)

## AUTHORS

Written by mikerquinn.

## BUILDING

step-renew is a single-file Go program with no external dependencies and no cgo, making cross-compilation straightforward.

### Native build (development machine)

```sh
go build -o step-renew main.go
```

Install system-wide:

```sh
go build -o /usr/local/bin/step-renew main.go
```

Run directly (for testing):

```sh
go run main.go -ca ... -cert ... -key ... -server ...
```

### Cross-compiling for embedded / non-x86 systems

Most deployments are on embedded devices that lack a Go toolchain. Build on your development machine (x86_64 or arm64) and copy the resulting static binary to the target.

**Recommended flags** (fully static binary, stripped):

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=<arch> go build -ldflags="-s -w" -o step-renew main.go
```

Common architectures:

- **arm64** (Raspberry Pi 4/5, most modern embedded boards):

  ```sh
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o step-renew main.go
  ```

- **arm** (32-bit ARM, older Raspberry Pi, many IoT/embedded devices):

  ```sh
  CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o step-renew main.go
  ```

- Other frequent embedded targets (replace `GOARCH`):

  ```sh
  CGO_ENABLED=0 GOOS=linux GOARCH=mips     go build ...
  CGO_ENABLED=0 GOOS=linux GOARCH=riscv64  go build ...
  CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le  go build ...
  CGO_ENABLED=0 GOOS=linux GOARCH=s390x    go build ...
  ```

Transfer the binary to the device (scp, etc.) and make executable:

```sh
chmod +x step-renew
```

The binary has no runtime dependencies and will run on any most any system that matches the target architecture.
