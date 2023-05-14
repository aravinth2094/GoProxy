package actions

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/aravinth2094/GoProxy/auth"
	"github.com/oov/socks5"
	"github.com/urfave/cli/v2"
)

func TunnelServer(ctx *cli.Context) error {
	if err := GenerateKeyPair(); err != nil {
		return err
	}
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return nil
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	listener, err := tls.Listen("tcp", ctx.String("listen"), config)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Println("Listening on ", ctx.String("listen"))
	targetAddr := ctx.String("target")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleClient(conn, targetAddr)
	}
}

func ProxyServer(ctx *cli.Context) error {
	server := socks5.New()
	server.AuthUsernamePasswordCallback = func(conn *socks5.Conn, username, password []byte) error {
		conn.Data = string(username)
		return auth.Authenticate(string(username), string(password))
	}
	server.HandleConnectFunc(func(conn *socks5.Conn, host string) (newHost string, err error) {
		return host, auth.CheckAllowed(conn.Data.(string), host)
	})
	return server.ListenAndServe(ctx.String("listen"))
}

func TunnelClient(ctx *cli.Context) error {
	listenAddr := ctx.String("listen")
	targetAddr := ctx.String("target")
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Relay server started on %s\n", listenAddr)

	for {
		// Accept incoming client connections
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept client connection: %v", err)
			continue
		}

		go handleClientTls(clientConn, targetAddr)
	}
}

func handleClient(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Connect to the target server
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Printf("Failed to connect to the target server: %v", err)
		return
	}
	defer targetConn.Close()

	// Start relaying messages between the client and target server
	go relayMessages(clientConn, targetConn)
	relayMessages(targetConn, clientConn)
}

func relayMessages(src, dest net.Conn) {
	_, err := io.Copy(dest, src)
	if err != nil {
		fmt.Printf("Error while relaying messages: %v", err)
	}
}

func handleClientTls(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Connect to the target server
	targetConn, err := tls.Dial("tcp", targetAddr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Printf("Failed to connect to the target server: %v", err)
		return
	}
	defer targetConn.Close()

	// Start relaying messages between the client and target server
	go relayMessages(clientConn, targetConn)
	relayMessages(targetConn, clientConn)
}

func GenerateKeyPair() error {
	// Generate a new RSA key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a self-signed X.509 certificate
	template := &x509.Certificate{
		SerialNumber:          bigIntOne(),
		Subject:               pkixName(),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return err
	}

	// Write the certificate and private key to files
	certFile, err := os.Create("server.crt")
	if err != nil {
		return err
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return err
	}

	keyFile, err := os.OpenFile("server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err != nil {
		return err
	}

	fmt.Println("X.509 key pair generated successfully.")
	return nil
}

func bigIntOne() *big.Int {
	return big.NewInt(1)
}

func pkixName() pkix.Name {
	return pkix.Name{
		Country:            []string{"IN"},
		Organization:       []string{"Example Inc."},
		OrganizationalUnit: []string{"IT"},
		CommonName:         "localhost",
	}
}
