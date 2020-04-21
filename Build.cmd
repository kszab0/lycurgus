@echo off
set out=lycurgus.exe

@if -%1-==-- goto all
@if [%1]==[all] goto all
@if [%1]==[dependencies] goto dependencies
@if [%1]==[clean] goto clean
@echo Build target not found
@exit /b

:all
@echo [Lycurgus] Building Lycurgus...
@setx GOARCH "amd64" >nul
@setx GOOS "windows" >nul
@go build -ldflags="-s -w" -o %out% 
@echo OK
@exit /b

:dependencies
@echo [Lycurgus] Installing dependencies...
@go get -u gopkg.in/elazarl/goproxy 
@echo OK
@exit /b

:clean
@echo [Lycurgus] Cleaning up...
@del %out%
@echo OK
@exit /b