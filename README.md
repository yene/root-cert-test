# Working with Root CA

## Notes
Do not use this document, these are just my learning notes.

Use a config file for creation, pass it in with `-config`.
If the certificate is in text format, then it is in PEM format.
Trusted root certificate is sometimes refered as anchor.
I recommend to use Linux with latest openssl for key creation.
Make sure older versions of TLS are not supported by your server, it should require 1.2
`curl --tlsv1.1 --cacert rootCa.crt https://www.mydomain.com`
`curl --tlsv1.0 --cacert rootCa.crt https://www.mydomain.com`

A lot of good stuff can be found in [Simple Golang HTTPS/TLS Examples](https://github.com/denji/golang-tls)

## Questions
* What is the right way to create multi domain certificates? Pass in the v3_req extension when generating csr and cert.
* Can you re-generate a new root cert? no
* How to use CAcreateserial? openssl creates a txt file and increments the serial inside for you.

## Create Root Key

**Attention:** this is the key used to sign the certificate requests, anyone holding this can sign certificates on your behalf. So keep it in a safe place!

```bash
openssl genrsa -out rootCA.key 4096
```

If you want a password protected key add the `-des3` option

## Create and self sign the Root Certificate

```bash
openssl req -new -x509 -sha256 -nodes -key rootCA.key -days 365000 -out rootCA.crt
```

Here we used our root key to create the root certificate that needs to be distributed in all the computers that have to trust us.

Include the purpose in the Common Name (root ca product x).


# Create a certificate for each server

This procedure needs to be followed for each server/appliance that needs a trusted certificate from our CA.

## Create the certificate key

```
openssl genrsa -out mydomain.com.key 2048
```

## Create the signing  (csr)

The certificate signing request is where you specify the details for the certificate you want to generate.
This request will be processed by the owner of the Root key (you in this case since you create it earlier) to generate the certificate.

**Important:** Please mind that while creating the signign request is important to specify the `Common Name` providing the IP address or domain name for the service, otherwise the certificate cannot be verified.

I will describe here two ways to generate them:

### Method A (Interactive)

If you generate the csr in this way, openssl will ask you questions about the certificate to generate like the organization details and the `Common Name` (CN) that is the web address you are creating the certificate for, e.g `mydomain.com`.

```
openssl req -new -key mydomain.com.key -out mydomain.com.csr
```

### Method B (config file)

Google around until you find the right combinations of paramter and config file and OpenSSL weirdness.

## Verify the csr's content

```
openssl req -in mydomain.com.csr -noout -text
```

## Generate the certificate using the `mydomain` csr and key along with the CA Root key

```
openssl x509 -req -in mydomain.com.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out mydomain.com.crt -days 36000 -sha256
```

## Verify the certificate's content

```
openssl x509 -in mydomain.com.crt -text -noout
```

## Verify issuer
```
openssl verify -untrusted rootCA.crt mydomain.com.crt
openssl verify -CAfile rootCA.crt mydomain.com.crt
```

## Multidomain SSL

> For everybody, who doesn´t like to edit the system-wide openssl.conf, there´s a native openssl CLI option for adding the SANs to the .crt from a .csr. All you have to use is openssl´s -extfile and -extensions CLI parameters.

> Extensions in certificates are not transferred to certificate requests and vice versa. Thats why we specify it with as parameter.

```
openssl req -new -key mydomain.com.key -out mydomain.com.csr -config mydomain.com.cnf
openssl x509 -req -in mydomain.com.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out mydomain.com.crt -days 36000 -sha256 -extfile mydomain.com.cnf -extensions v3_req
openssl x509 -in mydomain.com.crt -text -noout

```

A multi domain certificate should have a section with `X509v3 Subject Alternative Name`
https://stackoverflow.com/a/47779814 https://security.stackexchange.com/a/176084
