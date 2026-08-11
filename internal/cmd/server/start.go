package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/chrisarmitage/grpc-showcase/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"github.com/chrisarmitage/grpc-showcase/internal/grpcserver"
)

const (
	startFuncName = "start"
	startCmdDesc  = "Start the gRPC server"
	modeInsecure  = "insecure"
	modeTls       = "tls"
	modeMtls      = "mtls"
)

var tlsMode string

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   startFuncName,
		Short: startCmdDesc,
		Run:   runServerStart,
	}

	cmd.Flags().StringVar(&tlsMode, "tls", modeInsecure, "TLS mode (insecure|tls|mtls)")

	return cmd
}

func runServerStart(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var serverOptions []grpc.ServerOption

	switch tlsMode {
	case modeInsecure:
		log.Printf("Starting gRPC server in insecure mode")
	case modeTls:
		creds, err := credentials.NewServerTLSFromFile("pki/server.crt", "pki/server.key")
		if err != nil {
			log.Fatalf("Failed to load TLS credentials from pki/server.crt/pki/server.key: %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(creds))
		log.Printf("Starting gRPC server in TLS mode")
	case modeMtls:
		serverCert, err := tls.LoadX509KeyPair("pki/server.crt", "pki/server.key")
		if err != nil {
			log.Fatalf("Failed to load server certificate from pki/server.crt/pki/server.key: %v", err)
		}
		caCert, err := os.ReadFile("pki/ca.crt")
		if err != nil {
			log.Fatalf("Failed to read CA certificate from pki/ca.crt: %v", err)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("Failed to parse CA certificate from pki/ca.crt")
		}
		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caPool,
		}
		serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(tlsCfg)))
		serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(clientCertUnaryInterceptor))
		serverOptions = append(serverOptions, grpc.ChainStreamInterceptor(clientCertStreamInterceptor))
		log.Printf("Starting gRPC server in mTLS mode")
	default:
		log.Fatalf("Invalid TLS mode %q. Accepted values are: %s|%s|%s", tlsMode, modeInsecure, modeTls, modeMtls)
	}

	s := grpc.NewServer(serverOptions...)
	pb.RegisterUniServiceServer(s, &grpcserver.GrpcServer{})

	log.Printf("gRPC server listening on port 5050")
	log.Printf("Use Ctrl+C to stop the server")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func clientCertUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logClientCert(ctx)
	return handler(ctx, req)
}

func clientCertStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logClientCert(ss.Context())
	return handler(srv, ss)
}

func logClientCert(ctx context.Context) {
	p, ok := peer.FromContext(ctx)
	if !ok || p.AuthInfo == nil {
		return
	}
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok || len(tlsInfo.State.PeerCertificates) == 0 {
		return
	}
	cert := tlsInfo.State.PeerCertificates[0]
	log.Printf("Client TLS certificate:")
	log.Printf("  Subject:      %s", cert.Subject)
	log.Printf("  Organisation: %s", strings.Join(cert.Subject.Organization, ", "))
	log.Printf("  Issuer:       %s", cert.Issuer)
	log.Printf("  Valid from:   %s", cert.NotBefore.Format(time.RFC3339))
	log.Printf("  Valid until:  %s", cert.NotAfter.Format(time.RFC3339))
}
