@ECHO OFF
SETLOCAL
SET PLATFORMS=windows;linux;freebsd;darwin;dragonfly
SET CGO_ENABLED=0
CALL make.bat docs
FOR %%P IN (%PLATFORMS%) DO (
    ECHO Creating distfile for %%P...
    SET GOOS=%%P
    CALL make.bat && CALL make.bat dist
)
ECHO All builds done...
ECHO Calculating hashes...
powershell -NoLogo "Get-ChildItem -Filter dist\docproc-*.zip | %%{ $_.Name+' (MD5): '+(Get-FileHash $_.Fullname -Algorithm MD5 | Select-Object -ExpandProperty Hash)}"
ECHO done
ENDLOCAL
@ECHO ON
