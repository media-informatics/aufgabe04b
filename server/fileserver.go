package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/media-informatics/aufgabe04b/service"
	"google.golang.org/grpc"
)

type FileServer struct {
	service.UnimplementedFileContentServer
}

var (
	server = flag.String("server", service.Addr, "server address with port")
)

func (ds *FileServer) GetContent(in *service.FileName, srv service.FileContent_GetContentServer) error {
	fname := in.GetName()
	log.Printf("received: %s", fname)
	f, err := os.ReadFile(fname)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("could not read %s: %w", fname, err)
	}

	ctx := srv.Context()
	deadline, ok := ctx.Deadline()
	if ok {
		log.Printf("context aborts at %v", deadline)
	} else {
		log.Printf("context has not deadline set")
	}
	lines := strings.Split(string(f), "\n")
	log.Printf("read %d lines from %s", len(lines), fname)
	for i, l := range lines {
		select {
		case <-ctx.Done():
			err := fmt.Errorf("serving file content aborted after %d lines: %w", i, ctx.Err())
			log.Print(err)
			return err

		default:
			if err = srv.Send(&service.Line{Line: l}); err != nil {
				err = fmt.Errorf("send error after %d lines: %w", i, err)
				return err
			}
		}

	}
	log.Printf("finshed after %d lines", len(lines))
	return nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *server)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	service.RegisterFileContentServer(s, &FileServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
