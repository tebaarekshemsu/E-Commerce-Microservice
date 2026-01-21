# Health checks for E-Commerce stack (via Docker and HTTP)
$ErrorActionPreference = 'Stop'

Write-Host "== Broker =="
curl.exe -s http://localhost/ping | Out-String | Write-Host
curl.exe -s http://localhost/grpc/health | Out-String | Write-Host

Write-Host "== Product =="
curl.exe -s http://localhost:8082/product-service/ping | Out-String | Write-Host

Write-Host "== Order =="
curl.exe -s http://localhost:8083/order-service/health | Out-String | Write-Host

Write-Host "== User (HTTP) =="
# Expect 405 on GET register (indicates HTTP is up)
try { curl.exe -s -o NUL -w "%{http_code}" http://localhost:8001/api/auth/register/ | Write-Host } catch {}

Write-Host "== RabbitMQ (HTTP API) =="
$rbUser = $env:RABBIT_USER; if (-not $rbUser) { $rbUser = 'guest' }
$rbPass = $env:RABBIT_PASS; if (-not $rbPass) { $rbPass = 'guest' }
$pair = "$rbUser`:$rbPass"
$basic = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes($pair))
$headers = @{ Authorization = "Basic $basic" }
try { Invoke-RestMethod -Uri http://localhost:15672/api/overview -Headers $headers -TimeoutSec 5 | ConvertTo-Json -Compress | Write-Host } catch { Write-Host $_ }

Write-Host "== Postgres (exec SELECT 1) =="
try { docker compose exec -T postgres psql -U postgres -tAc "SELECT 1;" | Write-Host } catch { Write-Host $_ }

Write-Host "== Mongo (ping) =="
try { docker compose exec -T mongo mongosh --quiet --eval "db.adminCommand('ping').ok" | Write-Host } catch { Write-Host $_ }
