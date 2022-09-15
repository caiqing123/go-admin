### 代码与服务器端都需要客户端使用.proto生成的使用，因为代码双方需要数据结构通讯

- .proto是gRPC通信的数据结构文件，采用protobuf协议

```protobuf
syntax = "proto3";

package go.micro.grpc.user;
option go_package = "./;protos";

service User {
    rpc Add(AddRequest) returns (AddResponse) {}
}

message AddRequest {
    string Name = 1;
}

message AddResponse {
    int32 error_code = 1;
    string error_message = 2;
    int64 user_id = 3;
}
```

- 然后我们需要安装 gRPC 相关的编译程序：
- <https://www.cnblogs.com/oolo/p/11840305.html#%E5%AE%89%E8%A3%85-grpc>
- 我们开始编译.proto文件：
- 编译成功后会在当前目录生成protos/user.pb.go文件
- protoc --go_out=plugins=grpc:. user.proto

## chat

### Server

```bash
go run .\main.go -s -p "password"
```

### Client

```bash
go run .\main.go -h "127.0.0.1:6262" -p "password" -n "username"
```