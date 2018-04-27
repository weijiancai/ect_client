@echo off
set GOPATH=C:\workspace\my\ect_client


REM go install github.com/acroca/go-symbols
REM go install github.com/cweill/gotests/gotests
REM go install github.com/fatih/gomodifytags
REM go install github.com/golang/lint/golint
REM go install github.com/kardianos/govendor
REM go install github.com/nsf/gocode
REM go install github.com/ramya-rao-a/go-outline
REM go install github.com/rogpeppe/godef
REM go install github.com/sqs/goreturns
REM go install github.com/tpng/gopkgs
REM go install github.com/sqs/goreturns
REM go get github.com/denisenkom/go-mssqldb
go install ect
start "" "./bin/ect.exe"
REM go test -v ect