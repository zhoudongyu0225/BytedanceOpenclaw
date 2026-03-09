@echo off

echo ing....

rem 编译 go 程序
go run  ./toolProtoToPb/generateProto.go

echo ok
pause>nul