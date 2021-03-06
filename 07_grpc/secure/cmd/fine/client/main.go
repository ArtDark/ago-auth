package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	fineV1Pb "lectiongrpc/pkg/fine/v1"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9999"
const defaultHost = "netology.local"
const defaultCertificatePath = "./tls/certificate.pem"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	certificatePath, ok := os.LookupEnv("APP_CERTIFICATE_PATH")
	if !ok {
		certificatePath = defaultCertificatePath
	}

	if err := execute(net.JoinHostPort(host, port), certificatePath); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string, certificatePath string) (err error) {
	creds, err := credentials.NewClientTLSFromFile(certificatePath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := fineV1Pb.NewFineServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	response, err := client.FindByUserId(ctx, &fineV1Pb.FinesRequest{UserId: 1})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}

	log.Print(response)
	return nil
}
