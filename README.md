# RCD (wrapped Rclone) component
Rcd is a wrapped version of [Rclone](https://github.com/rclone/rclone "Rclone"). The wrapping is needed to be sure that server will not work if front-end dies (because front-end is responsible for verifyng user access).

To compile server (Rclone) with mount capablities sources of [WinFsp](https://github.com/winfsp/winfsp "WinFsp") are required and MiniGW-64 must be installed.

## How to setup
Update conf.go file with login, password and url password that are match to the ones you put in conf.go of Omega Drive client