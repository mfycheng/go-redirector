# go-redirector

go-redirector is a basic HTTP/HTTPS web server that allows multiple vanity go import paths.

### Installing

```
go get github.com/mfycheng/go-retriever
```

### Example Setup

Currently, I have two domains, `ephyra.io`, and `mfycheng.com` that point to a single
Digital Ocean server. Using go-redirector, both domains can be used for vanity import
paths such as:

```
import "mfycheng.com/go-redirector"
import "ephyra.io/pkg"
```

The configuration for this is simply:

```
{
    "tls": [
        {
            // Pre-LetsEncrypt era certs.
            "cert": "/etc/go-redirector/ephyra.io/server.crt",
            "key": "/etc/go-redirector/ephyra.io/server.key"
        },
        {
            // LetsEncrypt generated certs.
            "cert": "/etc/go-redirector/mfycheng.com/fullchain.pem",
            "key": "/etc/go-redirector/mfycheng.com/privkey.pem"
        }
    ],
    "redirections": [
        {
            "url": "ephyra.io",
            "redirect": "https://github.com/mfycheng/ephyra"
        },
        {
            "url": "mfycheng.com",
            "redirect": "https://github.com/mfycheng"
        }
    ]
}
```

### Sample Systemd Service File
```
[Unit]
Description=Go Redirector

[Service]
ExecStart=/usr/local/bin/go-redirector -config /etc/go-redirector/config.json

[Install]
WantedBy=multi-user.target
```
