package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	discordrichpresence "github.com/jacksonthemaster/discordrichpresence"
)

const (
	applicationID = "1560274580315177031"
	gameExe       = "NFS13_Game.exe"
	largeImageKey = "wnfsmw2012"
	largeImageTxt = "Need for Speed: Most Wanted 2012"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	messageBoxW     = user32.NewProc("MessageBoxW")
)

func showError(title, text string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	messageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		0x00000010,
	)
}

func main() {
	logFile, err := os.OpenFile("listener.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	exePath, err := os.Executable()
	if err != nil {
		showError("Error", "Error: Could not determine executable path.")
		return
	}
	gamePath := filepath.Join(filepath.Dir(exePath), gameExe)

	if _, err := os.Stat(gamePath); os.IsNotExist(err) {
		showError("Error", "Error: Game exe not found\n\nshould be named NFS13_Game.exe")
		return
	}

	client := discordrichpresence.NewClient(applicationID)

	activity := discordrichpresence.NewActivity().
		Details("Need for Speed: Most Wanted 2012").
		StartTime(time.Now()).
		LargeImage(largeImageKey, largeImageTxt).
		Type(0).
		Build()

	if err := client.StartWithActivity(activity, 30*time.Second); err != nil {
		showError("Error", "Error: Discord App not running")
		return
	}
	defer client.Close()

	cmd := exec.Command(gamePath)
	cmd.Dir = filepath.Dir(gamePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		showError("Error", "Error: Could not start game executable.")
		return
	}

	log.Printf("Game started (PID: %d)", cmd.Process.Pid)

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Exited Game with error: %v", err)
		} else {
			log.Println("Exited Game.")
		}
		client.ClearActivity()
		client.Close()
		os.Exit(0)
	}()

	select {}
}