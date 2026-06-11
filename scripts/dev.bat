@echo off
REM ENIAC Workspace Dev Script

if "%1"=="" goto dev

if "%1"=="build" goto build
if "%1"=="test" goto test
goto usage

:dev
echo Starting backend in dev mode...
cd /d "%~dp0..\backend"
go run . --dev
goto end

:build
echo Building backend...
cd /d "%~dp0.."
go build -o bin\eniac-workspace.exe .\backend\
echo Build complete: bin\eniac-workspace.exe
goto end

:test
echo Running tests...
cd /d "%~dp0..\backend"
go test ./... -v
goto end

:usage
echo Usage: %~nx0 {dev^|build^|test}
exit /b 1

:end
