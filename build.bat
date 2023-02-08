set appName=rcd
set CPATH=%~dp0winfsp\inc\fuse

rem Compile as Console app
go build -ldflags "-s -w" -tags cmount -o %appName%.exe src/main.go
