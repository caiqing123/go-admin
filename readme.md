## 说明

## 所有路由



## 所有命令

```
$ go run main.go -h
Default will run "serve" command, you can use "-h" flag to see all subcommands

Usage:
   [command]

Available Commands:
  cache       Cache management
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  key         Generate App Key, will print the generated Key
  make        Generate file and code
  make        Generate file and code
  migrate     Run database migration
  init        Initialize the database.
  seed        Insert fake data to the database
  serve       Start web server

Flags:
  -e, --env string   load .env file, example: --env=testing will use .env.testing file
  -h, --help         help for this command

Use " [command] --help" for more information about a command.
```

make 命令：

```
$ go run main.go make -h
Generate file and code

Usage:
   make [command]

Available Commands:
  apicontroller Create api controller，exmaple: make apicontroller v1/user
  cmd           Create a command, should be snake_case, exmaple: make cmd buckup_database
  factory       Create model's factory file, exmaple: make factory user
  migration     Create a migration file, example: make migration add_users_table
  model         Crate model file, example: make model user
  request       Create request file, example make request user
  seeder        Create seeder file, example:  make seeder user

Flags:
  -h, --help   help for make

Global Flags:
  -e, --env string   load .env file, example: --env=testing will use .env.testing file

Use " make [command] --help" for more information about a command.
```

migrate 命令：

```
$ go run main.go migrate -h
Run database migration

Usage:
   migrate [command]

Available Commands:
  down        Reverse the up command
  fresh       Drop all tables and re-run all migrations
  refresh     Reset and re-run all migrations
  reset       Rollback all database migrations
  up          Run unmigrated migrations

Flags:
  -h, --help   help for migrate

Global Flags:
  -e, --env string   load .env file, example: --env=testing will use .env.testing file

Use " migrate [command] --help" for more information about a command.
```

打包：

```
go build main.go
go build -ldflags="-s -w" main.go  //-s：忽略符号表和调试信息。 -w：忽略DWARFv3调试信息，使用该选项后将无法使用gdb进行调试。
upx -9 main.exe //upx 压缩包大小

cd到main.go所在目录，执行命令：
SET CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build main.go
打包后会生成一个main程序，将此程序拷贝至linux服务器，两种方式启动：

1、在当前会话执行
./main
2、后台启动
setsid ./main
```

mysql 5.7

win运行
修改.env配置
./main.exe init
./main.exe
