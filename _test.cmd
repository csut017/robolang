@ECHO off
SETLOCAL ENABLEEXTENSIONS
SET me=%~n0

ECHO %me%: Running tests
go test

IF %ERRORLEVEL% NEQ 0 (
    ECHO %me%: One or more unit tests failed!
    EXIT /B 0
)

ECHO %me%: Generating coverage
go test -covermode=count -coverprofile=coverage.out

go tool cover -html=coverage.out