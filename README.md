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

This tool is intended for automated certificate renewal in scripts, cron jobs, or systemd timers where short-lived certificates from step-ca are in use.

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
