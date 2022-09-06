package serve

import (
	"context"

	pb "api/pkg/grpc/proto"
)

type UserService struct {
}

func (t *UserService) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	// 执行数据库操作
	// ...
	resp := pb.AddResponse{
		ErrorCode:    0,
		ErrorMessage: "",
		UserId:       10001,
	}
	return &resp, nil
}

func (t *UserService) List(ctx context.Context, in *pb.Empty) (*pb.AddResponse, error) {
	// 执行查询数据库操作
	// ...
	list := make(map[string]string)
	list["demo"] = "demo"
	list["dem1"] = "demo1"

	resp := pb.AddResponse{
		ErrorCode:    0,
		ErrorMessage: "",
		List:         list,
	}
	return &resp, nil
}
