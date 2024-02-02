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

	cfgfilename := os.Args[0] + ".env"
	if _, err := os.Stat(cfgfilename); err != nil {
		cfgfilename = ""
	}
	myApp := app.NewApp(cfgfilename)
	myAppStop := myApp.Start()
	defer func() {
		if err := recover(); err != nil {
			myAppStop(err.(error))
			os.Exit(1)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig
	myAppStop(nil)
}
