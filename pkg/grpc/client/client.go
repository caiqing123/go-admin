package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	pb "api/pkg/grpc/proto"
)

type GrpcClientCommand struct {
}

func (t *GrpcClientCommand) Main() {
	addr := "127.0.0.1:3006"
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	cli := pb.NewUserClient(conn)

	//添加
	req := pb.AddRequest{
		Name: "xiaoliu",
	}
	resp, err := cli.Add(ctx, &req)

	//列表
	//req := pb.Empty{}
	//resp, err := cli.List(ctx, &req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.List)
	fmt.Println(fmt.Sprintf("Add User: %d", resp.UserId))
}
