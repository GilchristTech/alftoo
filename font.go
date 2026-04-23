package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/adrg/sysfont"
)

func FontFindPath(name string) (string, error) {
	// First try fc-match, if its installed
	cmd := exec.Command("fc-match", "--format", "%{file}", name)
	var buf strings.Builder
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		// return "", fmt.Errorf("FontFindPath: fc-match \"%s\"\n - %w", name, err)
		goto SYSFONT_FALLBACK
	} else {
		return buf.String(), nil
	}

	return "", nil

SYSFONT_FALLBACK:
	if font := sysfont.NewFinder(nil).Match(name); font != nil {
		return font.Filename, nil
	} else {
		goto FAIL
	}

FAIL:
	return "", fmt.Errorf("FontFindPath: could not find font with name \"%s\"", name)
}
