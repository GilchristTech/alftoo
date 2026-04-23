package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var (
		err       error
		app       AlftooApp
		exit_code int = 0
	)

	app.SetDefaults()

	if err = app.Run(); err != nil {
		if exit_err, ok := err.(*exec.ExitError); ok {
			exit_code = exit_err.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			exit_code = 1
		}
	}

	os.Exit(exit_code)
}
