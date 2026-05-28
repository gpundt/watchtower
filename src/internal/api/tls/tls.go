package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
)

var TLSServer *http.Server

func InitializeServermTLS() {
	// 1. Load server's cert and key
	cert, loadErr := tls.LoadX509KeyPair(
		Config.ServerPaths.CertFilepath,
		Config.ServerPaths.KeyFilepath,
	)
	if loadErr != nil {
		log.Fatal().Err(loadErr).
			Str("func", "InitializeServermTLS").
			Msg("failed to load x509 keypair")
	}

	// 2. Load the CA cert that signed these certificates
	caCert, readErr := ioutil.ReadFile(Config.ServerPaths.CACertFilepath)
	if readErr != nil {
		log.Fatal().Err(readErr).
			Str("func", "InitializeServermTLS").
			Msg("failed to read CA cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	TLSServer := &http.Server{
		Addr:      fmt.Sprintf(":%d", Config.ServerConfig.Server.Port),
		TLSConfig: tlsConfig,
	}

	log.Fatal().Err(TLSServer.ListenAndServeTLS("", ""))
}

var AgentTLSClient *http.Client

func InitializeAgentmTLS() {
	// 1. Load the client's cert and key
	cert, loadErr := tls.LoadX509KeyPair(
		Config.AgentPaths.CertFilepath,
		Config.AgentPaths.KeyFilepath,
	)
	if loadErr != nil {
		log.Fatal().Err(loadErr).
			Str("func", "InitializeAgentmTLS").
			Msg("failed to load x509 keypair")
	}

	// 2. Load the CA cet that signed the server's certificate
	caCert, readErr := ioutil.ReadFile(Config.AgentPaths.CACertFilepath)
	if readErr != nil {
		log.Fatal().Err(readErr).
			Str("func", "InitializeAgentmTLS").
			Msg("failed to read CA cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. Setup HTTPs transport
	AgentTLSClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
}
