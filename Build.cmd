@echo off
set out=lycurgus.exe

@if -%1-==-- goto all
@if [%1]==[all] goto all
@if [%1]==[dev] goto dev
@if [%1]==[dependencies] goto dependencies
@if [%1]==[clean] goto clean
@echo Build target not found
@exit /b

:all
@echo [Lycurgus] Building Lycurgus...
@setx GOARCH "amd64" >nul
@setx GOOS "windows" >nul
@rsrc -arch amd64 -manifest lycurgus.manifest -ico=icon.ico -o rsrc.syso
@go build -ldflags="-s -w -H=windowsgui" -o %out%
@del rsrc.syso
@echo OK
@exit /b

:dev
@echo [Lycurgus] Building Lycurgus...
@setx GOARCH "amd64" >nul
@setx GOOS "windows" >nul
@rsrc -arch amd64 -manifest lycurgus.manifest -ico=icon.ico -o rsrc.syso
@go build -ldflags="-s -w" -o %out%
@del rsrc.syso
@echo OK
@exit /b

:dependencies
@echo [Lycurgus] Installing dependencies...
@go get gopkg.in/elazarl/goproxy.v1
@go get github.com/getlantern/systray
@go get github.com/akavel/rsrc
@echo OK
@exit /b

:clean
@echo [Lycurgus] Cleaning up...
@del %out%
@echo OK
@exit /b