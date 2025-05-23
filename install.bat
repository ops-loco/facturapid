@echo off
SETLOCAL

REM === Configuration ===
SET ServiceName=FacturapidSynchronizer
SET DisplayName="Facturapid Synchronizer Service"
SET ExecutableName=FacturapidSynchronizer.exe
REM Define the source path of your compiled executable
SET SourcePath=%~dp0%ExecutableName%
REM Define the installation directory
SET InstallDir=%ProgramFiles%\%ServiceName%

echo Installing %DisplayName%...

REM Check if running as administrator
REM (A more robust check might be needed for complex scenarios)
openfiles > NUL 2>&1
if %errorlevel% NEQ 0 (
    echo Requesting administrative privileges...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

REM Create installation directory
if not exist "%InstallDir%" (
    echo Creating directory: %InstallDir%
    mkdir "%InstallDir%"
    if errorlevel 1 (
        echo Failed to create directory. Aborting.
        goto :eof
    )
)

REM Copy executable
echo Copying %ExecutableName% to %InstallDir%...
copy /Y "%SourcePath%" "%InstallDir%\"
if errorlevel 1 (
    echo Failed to copy executable. Aborting.
    goto :eof
)

REM Install the service
echo Installing Windows Service...
pushd "%InstallDir%"
%ExecutableName% install
if errorlevel 1 (
    echo Failed to install service. Check logs or run %ExecutableName% install manually from %InstallDir%.
    popd
    goto :eof
)
popd

echo Service %DisplayName% installed successfully.
echo You may need to configure the service user account or startup type via services.msc.
echo To start the service, run: sc start %ServiceName%

:eof
pause
ENDLOCAL
