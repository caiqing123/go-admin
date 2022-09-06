package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	client1 "api/pkg/grpc/client"
	pd "api/pkg/grpc/proto"
	"api/pkg/grpc/serve"

	"google.golang.org/grpc"
)

func TestServe(t *testing.T) {

	// listen
	listener, err := net.Listen("tcp", "127.0.0.1:3006")
	if err != nil {
		panic(err)
	}

	// signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		if err := listener.Close(); err != nil {
			panic(err)
		}
	}()

	s := grpc.NewServer()
	pd.RegisterUserServer(s, &serve.UserService{})
	fmt.Println("serve 127.0.0.1:3006")

	if err := s.Serve(listener); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}

func TestClient(t *testing.T) {
	run := client1.GrpcClientCommand{}
	run.Main()
}
