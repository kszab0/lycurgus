all:
	@/bin/echo -n "[Lycurgus] Building Lycurgus... "
	@go build -ldflags="-s -w" -o lycurgus 
	@/bin/echo "OK"

clean:
	@/bin/echo -n "[Lycurgus] Cleaning up... "
	@rm -f ./lycurgus
	@/bin/echo "OK"

.PHONY: clean all