NAME := Lycurgus
APP := lycurgus
VERSION := `git rev-parse --short HEAD`
BUILD_DATE := `date +%FT%T%z`

all:
	@/bin/echo -n "[$(NAME)] Building... "
	@go build -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}" -o $(APP) 
	@/bin/echo "OK"

clean:
	@/bin/echo -n "[$(NAME)] Cleaning up... "
	@rm -f ./$(APP)
	@/bin/echo "OK"

.PHONY: clean all