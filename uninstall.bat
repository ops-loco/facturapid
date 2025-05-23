@echo off
SETLOCAL

REM === Configuration ===
SET ServiceName=FacturapidSynchronizer
SET DisplayName="Facturapid Synchronizer Service"
SET ExecutableName=FacturapidSynchronizer.exe
REM Define the installation directory (should match install.bat)
SET InstallDir=%ProgramFiles%\%ServiceName%

echo Uninstalling %DisplayName%...

REM Check if running as administrator
openfiles > NUL 2>&1
if %errorlevel% NEQ 0 (
    echo Requesting administrative privileges...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

REM Check if executable exists in the expected location
if not exist "%InstallDir%\%ExecutableName%" (
    echo Executable not found at %InstallDir%\%ExecutableName%.
    echo Please ensure the service was installed correctly or provide the correct path.
    echo You might need to manually run 'sc delete %ServiceName%' if the executable is missing.
    goto :eof
)

REM Uninstall the service
echo Uninstalling Windows Service...
pushd "%InstallDir%"
%ExecutableName% uninstall
if errorlevel 1 (
    echo Failed to uninstall service. Check logs or run %ExecutableName% uninstall manually.
    echo You might need to manually run 'sc stop %ServiceName%' and then 'sc delete %ServiceName%'.
    popd
    goto :eof
)
popd

REM Optional: Remove the installation directory
:promptDelete
echo Do you want to remove the installation directory (%InstallDir%)? (Y/N)
set /p "choice="
if /i "%choice%"=="Y" (
    echo Removing directory: %InstallDir%
    rmdir /S /Q "%InstallDir%"
    if errorlevel 1 (
        echo Failed to remove directory. You may need to remove it manually.
    ) else (
        echo Directory removed successfully.
    )
) else if /i "%choice%"=="N" (
    echo Installation directory not removed.
) else (
    echo Invalid choice. Please enter Y or N.
    goto :promptDelete
)

echo Service %DisplayName% uninstalled successfully.

:eof
pause
ENDLOCAL
