package main

import (
	"os"
	"os/signal"
	"runtime"
	app "shortlink2/internal"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1)
	// debug.SetGCPercent(100)
	// debug.SetMemoryLimit(2 831 155 200)

	myApp := app.NewApp()
	myAppStop := myApp.Start()

	defer func() {
		if err := recover(); err != nil {
			myAppStop(err.(error))
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for {
		s := <-sig
		switch s {
		case syscall.SIGHUP: // kill -SIGHUP <pid> // restarting all for sake of config reload
			myAppStop(nil)
			myApp = app.NewApp()
			myAppStop = myApp.Start()
		default:
			myAppStop(nil)
			os.Exit(0)
		}
	}
	/* 	sig := make(chan os.Signal, 1)
	   	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	   	<-sig */
}
