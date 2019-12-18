package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func CreateTlsConfiguration(certFile, keyFile, caFile *string, verifySsl bool) (t *tls.Config) {
	if StringWithNoSpace(*certFile) && StringWithNoSpace(*keyFile) && StringWithNoSpace(*caFile) {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			Panic("failed to LoadX509KeyPair(cert, key). %v", err)
		}
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			Panic("failed to read ca file. %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
