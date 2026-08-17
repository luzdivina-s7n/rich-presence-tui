@echo off
setlocal
cd /d "%~dp0"
where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH.
    exit /b 1
)
echo [1/5] go mod tidy...
call go mod tidy
if errorlevel 1 exit /b 1
echo [2/5] generating icon.ico...
call go run ./cmd/genicon .
if errorlevel 1 exit /b 1
echo [3/5] generating rsrc.syso...
set "RSRC=%USERPROFILE%\go\bin\rsrc.exe"
if not exist "%RSRC%" (
    echo rsrc not found, installing...
    call go install github.com/akavel/rsrc@latest
    if errorlevel 1 exit /b 1
)
"%RSRC%" -ico icon.ico -o rsrc.syso -arch amd64
if errorlevel 1 exit /b 1
if not exist bin mkdir bin
echo [4/5] building bin\rich-presence-tui.exe ...
call go build -trimpath -ldflags "-s -w" -o bin\rich-presence-tui.exe .
if errorlevel 1 exit /b 1
echo [5/5] done: bin\rich-presence-tui.exe
exit /b 0