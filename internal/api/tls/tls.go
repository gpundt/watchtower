package tls

import (
	"fmt"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	Config "watchtower/internal/config"
	Endpoints "watchtower/internal/api/endpoints"

	"github.com/rs/zerolog/log"
)

func InitializeServermTLS() {
	// 1. Load server's cert and key
	cert, loadErr := tls.LoadX509KeyPair(
		Config.ServerPaths.CertFilepath,
		Config.ServerPaths.KeyFilepath,
	)
	if loadErr != nil {
		log.Fatal().Err(loadErr).Str("func", "InitializeServermTLS").Msg("failed to load x509 keypair")
	}

	// 2. Load the CA cert that signed these certificates
	caCert, readErr := ioutil.ReadFile(Config.ServerPaths.CACertFilepath)
	if readErr != nil {
		log.Fatal().Err(readErr).Str("func", "InitializeServermTLS").Msg("failed to read CA cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs: caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", Config.ServerConfig.Server.Port),
		TLSConfig: tlsConfig,
	}

	server.ListenAndServeTLS("", "")
}

func InitializeAgentmTLS() {
	// 1. Load the client's cert and key
	cert, loadErr := tls.LoadX509KeyPair(
		Config.AgentPaths.CertFilepath,
		Config.AgentPaths.KeyFilepath,
	)
	if loadErr != nil {
		log.Fatal().Err(loadErr).Str("func", "InitializeAgentmTLS").Msg("failed to load x509 keypair")
	}

	// 2. Load the CA cet that signed the server's certificate
	caCert, readErr := ioutil.ReadFile(Config.AgentPaths.CACertFilepath)
	if readErr != nil {
		log.Fatal().Err(readErr).Str("func", "InitializeAgentmTLS").Msg("failed to read CA cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. Setup HTTPs transport
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(fmt.Sprintf("%s%s", Config.AgentConfig.Agent.ServerURL, Endpoints.HealthCheckEndpoint))
	if err != nil {
		log.Fatal().Err(err).Str("func", "InitializeAgentmTLS").Msg("Failed to contact api endpoint")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", body)
}