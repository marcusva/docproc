@SETLOCAL ENABLEDELAYEDEXPANSION

@SET /P VERSION=<VERSION
@FOR /F "tokens=*" %%A IN ('go env GOARCH') DO @SET GOARCH=%%A
@FOR /F "tokens=*" %%A IN ('go env GOOS') DO @SET GOOS=%%A
@SET EXT=
@IF "%GOOS%" == "windows" (
    SET EXT=.exe
)
@SET ZIPPER="C:\Program Files\7-Zip\7z.exe"

@SET APPS=docproc.fileinput;docproc.proc;docproc.webinput
@SET DISTFILES=LICENSE;README.md
@SET DDIRS=examples
@SET DISTNAME=docproc-%VERSION%-%GOOS%-%GOARCH%
@SET DISTDIR=dist\%DISTNAME%

@SET LDFLAGS=-X main.version=%VERSION%
@SET TAGS=beanstalk nats nsq
@IF "%CGO_ENABLED%" == "0" SET TAGS=%TAGS% netgo

@IF "%~1" == "" GOTO :all
@GOTO :%~1

:all
@CALL :docs
@CALL :vendor
@CALL :build
@GOTO :eof

:clean
@ECHO Cleaning up...
@IF EXIST dist RMDIR /S /Q dist
@IF EXIST doc\_build RMDIR /S /Q doc\_build
@IF EXISt vendor RMDIR /S /Q vendor
@GOTO :eof

:docs
@ECHO Building docs...
@CD doc && CALL make html && CD ..
@GOTO :eof

:vendor
@ECHO Fetching dependencies...
%GOPATH%\bin\dep ensure -v
@GOTO :eof

:build
@ECHO Building apps for %GOOS% on arch %GOARCH%...
FOR %%A IN (%APPS%) DO (
    go build -tags "%TAGS%" -ldflags "%LDFLAGS%" -o %DISTDIR%\%%A%EXT% ./%%A
)
@GOTO :eof

:dist
@ECHO Creating distfile...
@XCOPY /Q /E /I /Y doc\_build\html %DISTDIR%\doc
@FOR %%A IN (%DISTFILES%) DO @COPY %%A %DISTDIR%\%%A
@FOR %%A IN (%DDIRS%) DO @XCOPY /Q /E /I /Y %%A %DISTDIR%\%%A
@powershell -NoLogo Compress-Archive -Path %DISTDIR% -Force -CompressionLevel Optimal -DestinationPath dist\%DISTNAME%.zip
@RMDIR /S /Q %DISTDIR%
@GOTO :eof

:test
go test -tags "%TAGS%" -ldflags "%LDFLAGS%" ./...
@GOTO :eof

@ENDLOCAL
