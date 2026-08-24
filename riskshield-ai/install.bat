@echo off
title RiskShield AI - Dependency Installer
color 0B

echo ===================================================
echo     RiskShield AI Dependency Installer
echo ===================================================
echo.

:: Check for Administrator privileges
net session >nul 2>&1
if %errorLevel% == 0 (
    echo Administrator privileges confirmed.
) else (
    echo ERROR: This installer requires Administrator privileges.
    echo Please right-click the .exe and select "Run as administrator".
    pause
    exit /b 1
)

echo.
echo [1/4] Installing PostgreSQL 16...
winget install -e --id PostgreSQL.PostgreSQL.16 --accept-package-agreements --accept-source-agreements --silent

echo.
echo [2/4] Installing Node.js (for frontend)...
winget install -e --id OpenJS.NodeJS --accept-package-agreements --accept-source-agreements --silent

echo.
echo [3/4] Configuring Database...
:: Set up postgres password, user, and database
set PGPASSWORD=postgres
"C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "CREATE USER riskshield WITH PASSWORD 'riskshield';"
"C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "CREATE DATABASE riskshield OWNER riskshield;"

echo.
echo [4/4] Installing Frontend Dependencies...
cd /d "%~dp0frontend"
call npm install
call npm run build

echo.
echo ===================================================
echo Installation Complete! 
echo You can now use the RiskShield AI Launcher.
echo ===================================================
pause
