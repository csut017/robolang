@ECHO off
SETLOCAL ENABLEEXTENSIONS
SET me=%~n0

ECHO %me%: Compiling
go build 
