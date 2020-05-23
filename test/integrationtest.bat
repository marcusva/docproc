@ECHO OFF

SETLOCAL
SET FAILED=0
SET PNAME=docprocci
SET CIP=docprocci_docproc
SET DOCKER=docker
SET DOCKER_COMPOSE=docker-compose

ECHO Building docker environment...
%DOCKER% build -t docproc/base .
%DOCKER_COMPOSE% -p %PNAME% build

ECHO Starting docker environment...
%DOCKER_COMPOSE% -p %PNAME% up -d

ECHO Creating queues manually to speed up testing...
%DOCKER% exec -d %CIP%.fileinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
%DOCKER% exec -d %CIP%.webinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
%DOCKER% exec -d %CIP%.preproc_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=preprocessed
%DOCKER% exec -d %CIP%.renderer_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=rendered

REM Wait for all services to be started properly
TIMEOUT /T 5

ECHO Starting tests...

%DOCKER% cp examples/data/testrecords.csv %CIP%.fileinput_1:/app/data
%DOCKER% cp examples/data/raw.json %CIP%.webinput_1:/raw.json
%DOCKER% exec -d %CIP%.webinput_1 curl -X POST -H "Content-Type: application/json" ^
    --data @/raw.json http:/localhost/receive

REM TODO: Add fileupload test

TIMEOUT /T 20

%DOCKER% exec %CIP%.output_1 ls -al /app/output

REM DO NOT USE: the following line is to sync proper results with the test result dir
REM %DOCKER% cp %CIP%.output_1:/app/output/. ./test/results

%DOCKER% cp test/test-results.tar.gz %CIP%.output_1:/app
%DOCKER% exec %CIP%.output_1 tar -C /app -xzf test-results.tar.gz
%DOCKER% exec -it %CIP%.output_1 diff -Nur /app/output /app/test-results

IF %ERRORLEVEL% NEQ 0 (
    SET FAILED=1
    FOR %%A IN (%CIP%.fileinput_1 %CIP%.webinput_1 %CIP%.preproc_1 %CIP%.renderer_1 %CIP%.output_1) DO (
        %DOCKER% logs %%A 1> test\%%A.log 2>&1
    )
)

%DOCKER_COMPOSE% -p %PNAME% kill
%DOCKER_COMPOSE% -p %PNAME% rm -f

IF %FAILED% == 0 (
    ECHO Tests successful
) ELSE (
    ECHO Tests failed!
)
@ECHO ON
@EXIT /B %FAILED%
