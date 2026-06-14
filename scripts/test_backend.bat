@echo off
cd /d %~dp0\..

:: Backend (Go)
echo Testing backend with dynamic port discovery...
start /b cmd /c "backend\eniac.exe serve --debug"

timeout /t 2

:: Check backend
curl -s http://localhost/api/health || (
    echo [ERROR] Backend offline — 'eniac serve' failed.
    echo Check API on http://localhost/ with:
    netstat -ano | findstr LISTEN | findstr :0
    goto :error
)

echo [OK] Backend online!
echo -------------------------------------------------
goto :eof

:error
echo [ERROR] Backend offline. Run 'eniac serve' before desktop.
pause