package main

import (
	"com.sentry.dev/app/sqlight"
	"os"
)

func main() {
	dbPath := os.Args[1]
	command := os.Args[2]

	executor := sqlight.GetInstance(dbPath)
	defer executor.Close()
	executor.Execute(command)
}
