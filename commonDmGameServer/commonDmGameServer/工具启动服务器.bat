@echo off

call :buildServer
call :oneGroup

:buildServer
    echo "build server"
    go build -o ./Server.exe  ./main.go
GOTO:EOF

:oneGroup
    echo "start one group"
    start "Logic" ./Server.exe -f conf/logic_1.json
GOTO:EOF

