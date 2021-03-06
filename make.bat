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
@SET TAGS=beanstalk nsq
@IF "%CGO_ENABLED%" == "0" SET TAGS=%TAGS% netgo

@IF "%~1" == "" GOTO :all
@GOTO :%~1

:all
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

:docker-image
docker build -t docproc/base -f ./test/dockerfiles/Dockerfile .
@GOTO :eof

:nsq-test
docker-compose -f ./test/dockerfiles/docker-compose.yml build && ^
docker-compose -f ./test/dockerfiles/docker-compose.yml up --abort-on-container-exit && ^
docker-compose -f ./test/dockerfiles/docker-compose.yml down -v
@GOTO :eof

:beanstalk-test
docker-compose -f ./test/dockerfiles/docker-compose.beanstalk.yml build && ^
docker-compose -f ./test/dockerfiles/docker-compose.beanstalk.yml up --abort-on-container-exit && ^
docker-compose -f ./test/dockerfiles/docker-compose.beanstalk.yml down -v
@GOTO :eof

:integration-test
@CALL :docker-image
@CALL :nsq-test
@CALL :beanstalk-test
@GOTO :eof

@ENDLOCAL
