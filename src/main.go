package main

import (
	"app/src/conf"
	"app/src/settings"
	"app/src/state"

	"app/src/utils"
	"encoding/base64"

	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/s3"
	"github.com/rclone/rclone/cmd"
	_ "github.com/rclone/rclone/cmd/all"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/process"

	log "github.com/sirupsen/logrus"
)

type RcloneInput struct {
	Name string `json:"name,omitempty"`
}

type RcloneResponse struct {
	Error  string      `json:"error,omitempty"`
	Input  RcloneInput `json:"input,omitempty"`
	Path   string      `json:"path,omitempty"`
	Status int         `json:"status,omitempty"`
}

func main() {
	// Create log
	initLog()
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("App has started, PID is " + strconv.Itoa(os.Getpid()))

	// Load settings
	settings.LoadSettingsValues()
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Settings loaded:")
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info(state.SettingsValues)

	// Setup CRON on schedule
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))
	c.AddFunc("@every "+conf.CheckParentInterval, func() {
		verifyPID()
	})
	go c.Run()

	verifyPID()

	// Cache
	var cachePath []string
	if !state.SettingsValues.Cache.Disabled {
		if state.SettingsValues.Cache.OverridePath != "" {
			cachePath = append(cachePath, "--cache-dir")
			cachePath = append(cachePath, state.SettingsValues.Cache.OverridePath)
		} else {
			cachePath = append(cachePath, "--cache-dir")
			cachePath = append(cachePath, state.SettingsValues.Cache.DefaultPath)
		}
	}

	path := []string{"rclone", "rcd", "--rc-user", conf.RcUsername, "--rc-pass", conf.RcPassword, "--rc-addr", strings.Replace(strings.Replace(conf.RcHost, "http://", "", -1), "/", "", -1)}

	if len(cachePath) > 0 {
		path = append(path, cachePath...)
	}

	// Exec rclone with arguments
	os.Args = path
	cmd.Main()

}

// Init log
func initLog() {
	_ = os.Remove(conf.LogFilename)

	f, err := os.OpenFile(conf.LogFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Error opening file:" + err.Error())
	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}

// Verify PID
func verifyPID() {
	validResp := true

	// Verify parent PID with reference PID taken fron env var
	rpid, err := base64.StdEncoding.DecodeString(os.Getenv("MD_PID"))
	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Can't decode reference PID")
		validResp = false
	}

	ppid := os.Getppid()

	isValid, err := process.PidExists(int32(ppid))

	if err != nil {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Parent PID does not exist. Details:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning(err)
		validResp = false
	}

	if !isValid {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Parent PID is not valid")
		validResp = false
	}

	if string(rpid) != strconv.Itoa(ppid) {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Parent PID does not match with reference")
		validResp = false
	}

	if !validResp {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Error("Quit due to error")
		Quit()
	}
}

// Quit Rclone Server
func Quit() {
	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Quiting Rclone has started...")

	var resp RcloneResponse
	path := "mount/unmountall"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Error unmounting all letters via 'mount/unmountall'. Details:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning(resp)
	}

	path = "fscache/clear"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Error clearing cache via 'fscache/clear'. Details:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning(resp)
	}

	path = "core/quit"
	utils.Request(state.RcAuthEncoded, http.MethodPost, conf.RcHost+path, &resp, nil)
	if len(resp.Error) > 0 {
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning("Error quitting Rclone via 'core/quit'. Details:")
		log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Warning(resp)
	}

	log.WithFields(log.Fields{"func": utils.CallerFuncLoc(), "loc": utils.CallerFileLoc()}).Info("Quiting Rclone has finished")
}
