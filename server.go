package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"text/template"
)

var redirectTemplate = template.Must(template.New("redirect").Parse(`
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
        <meta name="go-import" content="{{.URL}} git {{.RedirectURL}}">
        <meta http-equiv="refresh" content="0; url={{.RedirectURL}}">
        <!--
        <meta name="go-import" content="ephyra.io git https://github.com/mfycheng/ephyra">
        <meta http-equiv="refresh" content="0; url=https://github.com/mfycheng/ephyra">
        -->
    </head>
    <body>
        <p>Nothing to see here...</p>
    </body>
</html>
`))

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	if v, ok := redirectMap[req.Host]; ok {
		redirectTemplate.Execute(w, Redirect{URL: req.Host, RedirectURL: v})
	} else {
		log.Println("Unknown host:", req.Host)
		http.Error(w, "Unknown host", http.StatusBadRequest)
	}
}

// createTLSConfig creates a tls.Config for multiple cert/key pairs.
func createTLSConfig(tlsConfigs []TLSConfig) *tls.Config {
	var err error
	tlsConfig := &tls.Config{
		Certificates: make([]tls.Certificate, len(tlsConfigs)),
	}

	for i, conf := range tlsConfigs {
		tlsConfig.Certificates[i], err = tls.LoadX509KeyPair(conf.CertFile, conf.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}

func listenAndServeTLS(config Config, wg *sync.WaitGroup) {
	defer wg.Done()

	tlsConfig := createTLSConfig(config.TLSConfigs)
	server := &http.Server{
		TLSConfig: tlsConfig,
	}

	lis, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	log.Fatal(server.Serve(lis))
}

func listenAndServe(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Fatal(http.ListenAndServe(":80", nil))
}
