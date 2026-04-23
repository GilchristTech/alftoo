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

	if err := app.SetDefaults(); err != nil {
		fmt.Fprintf(os.Stderr, "Alftoo Error: could not set defaults\n - %s\n", err)
		os.Exit(1)
	}

	if err = app.Run(); err != nil {
		if exit_err, ok := err.(*exec.ExitError); ok {
			exit_code = exit_err.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "Alftoo Error:\n - %s\n", err)
			exit_code = 1
		}
	}

	os.Exit(exit_code)
}
