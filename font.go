package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/adrg/sysfont"
	"github.com/veandco/go-sdl2/sdl"
)

type FontCommand string

var command_font FontCommand

func init() {
	AddColonCommand(":font", &command_font)
}

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

func (fc *FontCommand) Run(a *AlftooApp, name, args string) error {
	args = strings.TrimSpace(args)

	if fp, err := FontFindPath(args); err != nil {
		return err

	} else {
		*fc = FontCommand(fp)

		if a.window != nil {
			fmt.Println("Set font:", fp)
		}

		if err := a.SetFontPath(fp, a.font_size); err != nil {
			return fmt.Errorf("FontCommand.Run: error setting font\n - %w", err)
		}

		if a.window != nil {
			a.Draw()
			sdl.Delay(1000)
		}

		return nil
	}
}
