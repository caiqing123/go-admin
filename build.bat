@echo off
:loop
@echo off&amp;color 0A
cls
echo,
echo ��ѡ��Ҫ�����ϵͳ������
echo,
echo 1. Windows_amd64
echo 2. linux_amd64
echo 3. Mac_amd64
echo 4. All
echo 0. quit
echo,

set/p action=��ѡ��:
if %action% == 1 goto build_Windows_amd64
if %action% == 2 goto build_linux_amd64
if %action% == 3 goto build_Mac_amd64
if %action% == 4 goto all
if %action% == 0 goto end
cls &amp; goto :loop

:build_Windows_amd64
echo ����Windows�汾64λ...
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/windows/amd64/app.exe main.go
echo ѹ���ļ�...
upx -9 ./release/windows/amd64/app.exe
goto end

:build_linux_amd64
echo ����Linux�汾64λ...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/linux/amd64/app main.go
echo ѹ���ļ�...
upx -9 ./release/linux/amd64/app
goto end

:build_Mac_amd64
echo ����Mac�汾64λ...
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/mac/amd64/app main.go
echo ѹ���ļ�...
upx -9 ./release/mac/amd64/app
goto end

:all
echo ׼���������а汾�������ĵȴ�...
timeout /t 3 /nobreak
echo,

echo ����Windows�汾64λ...
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/windows/amd64/app.exe main.go
echo ѹ���ļ�...
upx -9 ./release/windows/amd64/app.exe

echo ===============���Ƿָ���=====================

echo ����Linux�汾64λ...
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/linux/amd64/app main.go
echo ѹ���ļ�...
upx -9 ./release/linux/amd64/app

echo ===============���Ƿָ���=====================

echo ����Mac�汾64λ...
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-s -w" -o ./release/mac/amd64/app main.go
echo ѹ���ļ�...
upx -9 ./release/mac/amd64/app

:end
@exit