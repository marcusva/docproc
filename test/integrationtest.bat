@ECHO OFF

SETLOCAL
SET FAILED=0
SET PNAME=docprocci
SET CIP=docprocci_docproc
SET DOCKER=docker
SET DOCKER_COMPOSE=docker-compose

ECHO Building docker environment...
%DOCKER% build -t docproc/base .
%DOCKER_COMPOSE% -p %PNAME% build --no-cache

ECHO Starting docker environment...
%DOCKER_COMPOSE% -p %PNAME% up -d

ECHO Creating queues manually to speed up testing...
%DOCKER% exec -d %CIP%.fileinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
%DOCKER% exec -d %CIP%.preproc_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=preprocessed
%DOCKER% exec -d %CIP%.renderer_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=rendered

ECHO Starting tests...
%DOCKER% cp examples/data/testrecords.csv %CIP%.fileinput_1:/app/data

TIMEOUT /T 10

REM DO NOT USE: the following lines are to sync proper results with the test result dir
REM %DOCKER% exec %CIP%.output_1 ls -al /app/output
REM %DOCKER% cp %CIP%.output_1:/app/output/. ./test/results

%DOCKER% cp test/test-results.tar.gz %CIP%.output_1:/app
%DOCKER% exec %CIP%.output_1 tar -C /app -xzf test-results.tar.gz
%DOCKER% exec -it %CIP%.output_1 diff -Nur /app/output /app/test-results

IF %ERRORLEVEL% NEQ 0 (
    SET FAILED=1
)

%DOCKER_COMPOSE% -p %PNAME% kill
%DOCKER_COMPOSE% -p %PNAME% rm -f

IF %FAILED% == 0 (
    ECHO Tests successful
) ELSE (
    ECHO Tests failed!
)
EXIT /B %FAILED%
