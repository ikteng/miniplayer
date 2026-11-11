Write-Host "Starting Miniplay..." -ForegroundColor Green
Write-Host ""

$root = $PSScriptRoot

foreach ($port in 9863, 5173) {
    $pids = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique
    foreach ($procId in $pids) {
        Stop-Process -Id $procId -Force -ErrorAction SilentlyContinue
    }
}

if (-not (Test-Path "$root\backend\server.exe")) {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\backend'; go build -o server.exe ./cmd/server; ./server.exe"
} else {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\backend'; ./server.exe"
}

Start-Sleep -Seconds 3

Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$root\frontend'; npm install; npx vite"

Write-Host ""
Write-Host "Backend: http://localhost:9863" -ForegroundColor Yellow
Write-Host "Frontend: http://localhost:5173" -ForegroundColor Yellow
