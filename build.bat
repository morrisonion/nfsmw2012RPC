@echo off
setlocal

echo nfsmw2012RPC

go build -ldflags "-H=windowsgui -s -w" -o NFS13.exe main.go

if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Build failed.
    pause
    exit /b %errorlevel%
)

echo.
echo [OK] Build complete: NFS13.exe
echo.
pause
endlocal