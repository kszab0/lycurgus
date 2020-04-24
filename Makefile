all:
	@/bin/echo -n "[Lycurgus] Building Lycurgus... "
	@go build -ldflags="-s -w" -o lycurgus 
	@/bin/echo "OK"

dependencies:
	@/bin/echo -n "[Lycurgus] Installing dependencies..."
	@go get gopkg.in/elazarl/goproxy.v1 
	@go get github.com/getlantern/systray
	@/bin/echo "OK"

clean:
	@/bin/echo -n "[Lycurgus] Cleaning up... "
	@rm -f ./lycurgus
	@/bin/echo "OK"

.PHONY: dependencies clean all