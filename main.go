package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type renewResponse struct {
	Crt       string   `json:"crt"`
	CA        string   `json:"ca"`
	CertChain []string `json:"certChain"`
}

func main() {
	caPath := flag.String("ca", "", "CA certificate path (required)")
	certPath := flag.String("cert", "", "Client certificate path (required)")
	keyPath := flag.String("key", "", "Client key path (required)")
	serverHost := flag.String("server", "", "CA hostname or host:port (required)")

	flag.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", name)
		fmt.Fprintln(os.Stderr, "Renew a certificate from step-ca using mTLS. Writes the full chain returned by the server.")
		fmt.Fprintln(os.Stderr, "\nRequired flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintf(os.Stderr, "  %s -ca root_ca.crt -cert ssl-cert.pem -key ssl-cert.key -server step.mydomain.com\n", name)
	}

	flag.Parse()

	if *caPath == "" || *certPath == "" || *keyPath == "" || *serverHost == "" {
		flag.Usage()
		os.Exit(1)
	}

	serverURL := buildServerURL(*serverHost)

	rootPEM, _ := os.ReadFile(*caPath)
	certPEM, _ := os.ReadFile(*certPath)
	keyPEM, _ := os.ReadFile(*keyPath)

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(rootPEM)

	cert, _ := tls.X509KeyPair(certPEM, keyPEM)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      roots,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Post(serverURL, "application/json", strings.NewReader("{}"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result renewResponse
	json.Unmarshal(body, &result)

	fullChain := buildFullChain(result)
	if fullChain == "" {
		fmt.Println("No certificate returned")
		os.Exit(1)
	}

	writeFileAtomic(*certPath, fullChain)
	fmt.Println("Renewed successfully (full chain written)")
}

func buildServerURL(host string) string {
	if strings.HasPrefix(host, "https://") {
		return host
	}
	if !strings.Contains(host, ":") {
		host += ":443"
	}
	return "https://" + host + "/1.0/renew"
}

func buildFullChain(r renewResponse) string {
	if len(r.CertChain) > 0 {
		return strings.Join(r.CertChain, "\n")
	}
	if r.Crt != "" && r.CA != "" {
		return r.Crt + "\n" + r.CA
	}
	return r.Crt
}

func writeFileAtomic(path, content string) {
	tmp := path + ".new"
	os.WriteFile(tmp, []byte(content), 0644)
	os.Rename(tmp, path)
}
