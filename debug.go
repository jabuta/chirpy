package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func debugMode(dbPath string) string {
	dbg := flag.Bool("debug", false, "Enable Debug Mode, DB will be deleted")
	flag.Parse()

	if !*dbg {
		return dbPath
	}
	debugDbPath := "debug-" + dbPath

	os.Remove(debugDbPath)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		os.Remove(debugDbPath)
		log.Println("Debug mode enabled, database deleted before exit")
		os.Exit(0)
	}()
	return debugDbPath
}
