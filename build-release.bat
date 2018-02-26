@SETLOCAL ENABLEDELAYEDEXPANSION

@SET PLATFORMS=windows;linux;freebsd;darwin;dragonfly
@SET APPS=docproc.fileinput;docproc.proc
@SET FILES=LICENSE;README.md
@SET FOLDERS=examples
@FOR /F "tokens=*" %%A IN ('go env GOARCH') DO @SET ARCH=%%A
@SET /P VERSION=<VERSION

@ECHO Creating release packages for version %VERSION%...

@RMDIR /S /Q dist
@RMDIR /S /Q doc\_build

@MKDIR dist
@ECHO Creating documentation...
@CD doc
@CALL make html
@CD ..
@FOR %%P IN (%PLATFORMS%) DO (
    SET SUFFIX=""
    if "%%P" == "windows" (
        SET SUFFIX=.exe
    )
    SET DISTNAME=docproc-%VERSION%-%%P-%ARCH%
    SET DESTDIR=dist\!DISTNAME!
    ECHO Building release for %%P in !DESTDIR!...
    XCOPY /Q /E /I doc\_build\html !DESTDIR!\doc
    ECHO Building application...
    FOR %%A IN (%APPS%) DO (
        go build -tags "beanstalk nats nsq" -o !DESTDIR!\%%A!SUFFIX! ./%%A
    )
    ECHO Copying dist files...
    FOR %%A IN (%FOLDERS%) DO XCOPY /Q /E /I %%A !DESTDIR!\%%A
    FOR %%A IN (%FILES%) DO XCOPY /Q %%A !DESTDIR!

    ECHO Creating package...
    powershell -NoLogo Compress-Archive -Path !DESTDIR! -CompressionLevel Optimal -DestinationPath dist\!DISTNAME!.zip
    RMDIR /S /Q !DESTDIR!
)

@ECHO Calculating hashes...
@powershell -NoLogo "Get-ChildItem -Filter dist\docproc-*.zip | %%{ $_.Name+' (MD5): '+(Get-FileHash $_.Fullname -Algorithm MD5 | Select-Object -ExpandProperty Hash)}"
@ECHO done
@ENDLOCAL