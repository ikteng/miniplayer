# Builds a standalone miniplayer.exe: runs the frontend production build,
# copies it into backend/internal/staticassets/dist for embedding, then
# compiles the Go server with the embedded frontend baked in.
#
# Output: dist/miniplayer.exe + dist/assets/AudioService/ (sibling folder,
# required at runtime since the AudioService subprocess can't be embedded).
# Run the exe from inside dist/ - settings/ is created next to it on first run.

$root = $PSScriptRoot
$embedDist = "$root\backend\internal\staticassets\dist"
$outDir = "$root\dist"

Write-Host "Building frontend..." -ForegroundColor Green
Push-Location "$root\frontend"
npm run build
Pop-Location
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Embedding frontend into backend..." -ForegroundColor Green
Remove-Item "$embedDist\*" -Recurse -Force
Copy-Item "$root\frontend\dist\*" $embedDist -Recurse -Force

Write-Host "Building server..." -ForegroundColor Green
New-Item -ItemType Directory -Force $outDir | Out-Null
Push-Location "$root\backend"
go build -ldflags "-H=windowsgui" -o "$outDir\miniplayer.exe" ./cmd/server
Pop-Location
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Copying AudioService..." -ForegroundColor Green
New-Item -ItemType Directory -Force "$outDir\assets" | Out-Null
Copy-Item "$root\backend\assets\AudioService" "$outDir\assets\AudioService" -Recurse -Force

Write-Host ""
Write-Host "Done: $outDir\miniplayer.exe" -ForegroundColor Yellow
Write-Host "Run it from inside $outDir - AudioService must stay next to it." -ForegroundColor Yellow
