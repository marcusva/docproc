@ECHO OFF

SETLOCAL
SET PNAME=docprocperf
SET CIP=docprocperf_docproc
SET DOCKER=docker
SET DOCKER_COMPOSE=docker-compose

ECHO Building docker environment...
%DOCKER% build -t docproc/base .
%DOCKER_COMPOSE% -p %PNAME% build

ECHO Starting docker environment...
%DOCKER_COMPOSE% -p %PNAME% up -d

ECHO Creating queues manually to speed up testing...
%DOCKER% exec -d %CIP%.fileinput_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=input
%DOCKER% exec -d %CIP%.preproc_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=preprocessed
%DOCKER% exec -d %CIP%.renderer_1 curl -X POST http://127.0.0.1:4151/topic/create?topic=rendered

ECHO Starting performance tests...
%DOCKER% cp examples/data/performance.csv %CIP%.fileinput_1:/app/data

TIMEOUT /T 10

%DOCKER% exec -it %CIP%.output_1 cat /app/output/performance.txt
FOR %%A IN (%CIP%.fileinput_1 %CIP%.preproc_1 %CIP%.renderer_1 %CIP%.output_1) DO (
    %DOCKER% logs %%A 1> test\%%A.log 2>&1
)

%DOCKER_COMPOSE% -p %PNAME% kill
%DOCKER_COMPOSE% -p %PNAME% rm -f

@ECHO ON
