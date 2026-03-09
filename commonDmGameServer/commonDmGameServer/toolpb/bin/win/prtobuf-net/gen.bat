@echo off
set CurrentDir=%~dp0
set ProtoDir=%CurrentDir%
set JavaOutDir=%CurrentDir%

cd /d "%CurrentDir%"
setlocal enabledelayedexpansion

for %%i in (*.proto) do (
protogen -ns:pb  -i:%%i -o:.\output\%%~ni.cs)

pause