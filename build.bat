@echo off
echo [1/3] Building core...
go build -o cmd\loader\core.exe ./cmd/core

if errorlevel 1 (
    echo Failed to build core.
    exit /b 1
)

echo [2/3] Building loader with embedded core.exe...
go build -o insyra-insights.exe ./cmd/loader

if errorlevel 1 (
    echo Failed to build loader.
    exit /b 1
)

echo [3/3] Done! Output: insyra-insights.exe
