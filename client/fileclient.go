package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/media-informatics/aufgabe04b/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	fname  = flag.String("dir", ".", "directory to list")
	server = flag.String("server", service.Addr, "server address with port")
)

func main() {
	flag.Parse()
	// replace deprecated grpc.WithInsecure():
	cred := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(*server, cred)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := service.NewFileContentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := c.GetContent(ctx, &service.FileName{Name: *fname})
	if err != nil {
		log.Fatalf("did not receive stream: %v", err)
	}
	fmt.Printf("content of %s:\n\n", *fname)
	var builder strings.Builder
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			log.Printf("file content received")
			break
		}
		if err != nil {
			log.Fatalf("receive aborted %v", err)
		}
		builder.WriteString(line.GetLine())
		builder.WriteRune('\n')
	}
	fmt.Println(builder.String())
}
