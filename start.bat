@echo off
echo Starting Miniplay...

for %%P in (9863 5173) do (
    for /f "tokens=5" %%A in ('netstat -ano ^| findstr /r /c:"127.0.0.1:%%P " /c:"0.0.0.0:%%P "') do (
        taskkill /PID %%A /F >nul 2>&1
    )
)

start "Miniplay Backend" cmd /k "cd /d %~dp0backend && python -m uvicorn app.main:app --host 127.0.0.1 --port 9863"

if not exist "%~dp0frontend\node_modules" (
    start "Miniplay Frontend" cmd /k "cd /d %~dp0frontend && npm install && npx vite"
) else (
    start "Miniplay Frontend" cmd /k "cd /d %~dp0frontend && npx vite"
)

echo Both services starting...
