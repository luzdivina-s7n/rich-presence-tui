APP_NAME  := rich-presence-tui
OUT_DIR   := bin
UNAME_S   := $(shell uname -s)
ifeq ($(OS),Windows_NT)
	EXT := .exe
else
	EXT :=
endif
.PHONY: all build icon run vet fmt clean
all: build
icon:
	go run ./cmd/genicon .
	@if not exist "$(USERPROFILE)\go\bin\rsrc.exe" go install github.com/akavel/rsrc@latest
	"$(USERPROFILE)\go\bin\rsrc.exe" -ico icon.ico -o rsrc.syso -arch amd64
build: icon
	go mod tidy
	mkdir -p $(OUT_DIR)
	go build -trimpath -ldflags "-s -w" -o $(OUT_DIR)/$(APP_NAME)$(EXT) .
run: build
	$(OUT_DIR)/$(APP_NAME)$(EXT)
vet:
	go vet ./...
fmt:
	go fmt ./...
clean:
	rm -rf $(OUT_DIR)