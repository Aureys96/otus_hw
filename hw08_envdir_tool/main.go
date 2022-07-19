package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Should be at least two args")
	}
	args := os.Args

	args = args[1:]

	envDir := args[0]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	cmdArgs := args[1:]
	os.Exit(RunCmd(cmdArgs, env))
}
