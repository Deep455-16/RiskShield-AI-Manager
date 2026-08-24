@echo off
title RiskShield AI Launcher
color 0A

echo Starting RiskShield AI Core Services...

:: Clean up any existing orphaned processes from previous runs to prevent port conflicts
powershell -Command "Stop-Process -Name server -Force -ErrorAction SilentlyContinue"
powershell -Command "Stop-Process -Name node -Force -ErrorAction SilentlyContinue"

cd /d "%~dp0backend"
set DATABASE_URL=postgres://riskshield:riskshield@localhost:5432/riskshield?sslmode=disable
set DEMO_MODE=true
set PORT=8080
set JWT_SECRET=dev-secret-key-123
start "RiskShield Backend" /MIN server.exe

cd /d "%~dp0frontend"
:: CRITICAL: Reset PORT to 3000 so Next.js doesn't try to use the backend's port!
set PORT=3000
start "RiskShield Frontend" /MIN npm start

:: Use ping for a safe delay that works in all environments
ping 127.0.0.1 -n 6 > nul

:: Launch Microsoft Edge in App Mode
start msedge --app="http://localhost:3000" --profile-directory="Default" || start chrome --app="http://localhost:3000" || start http://localhost:3000
