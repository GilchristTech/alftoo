package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/adrg/sysfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type CommandFont string

var command_font CommandFont

func init() {
	AddColonCommand(":font", &command_font)
}

func SetFontPath(fpath string, size_px int) error {
	if alf_font != nil {
		alf_font.Close()
		alf_font = nil
	}

	if font, err := ttf.OpenFont(fpath, size_px*3/4); err != nil {
		return fmt.Errorf(
			"alftoo.SetFontPath: Failed to load font at path:\n     \"%s\"\n - %w",
			fpath,
			err,
		)
	} else {
		alf_font = font
		alf_font_fp = fpath
	}

	return nil
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

func (fc *CommandFont) Run(name, args string) error {
	var err_h string = "alftoo.CommandFont.Run"
	args = strings.TrimSpace(args)

	if fp, err := FontFindPath(args); err != nil {
		return fmt.Errorf("%s: error finding font path\n - %w", err_h, err)

	} else {
		*fc = CommandFont(fp)

		if alf_window != nil {
			fmt.Println("Set font:", fp)
		}

		if err := SetFontPath(fp, alf_font_size); err != nil {
			return fmt.Errorf("%s: error setting font\n - %w", err_h, err)
		}

		if alf_window != nil {
			Draw()
			sdl.Delay(1000)
		}

		return nil
	}
}

func RenderText(text string) (
	texture *sdl.Texture,
	surface *sdl.Surface,
	err error,
) {
	if surface, err = alf_font.RenderUTF8Blended(
		text,
		alf_font_color,
	); err != nil {
		return
	}

	if texture, err = alf_renderer.CreateTextureFromSurface(surface); err != nil {
		return
	}

	return
}

func RenderTextWrapped(text string, wrap_length_px int) (
	texture *sdl.Texture,
	surface *sdl.Surface,
	err error,
) {
	if surface, err = alf_font.RenderUTF8BlendedWrapped(
		text,
		alf_font_color,
		wrap_length_px,
	); err != nil {
		return
	}

	if texture, err = alf_renderer.CreateTextureFromSurface(surface); err != nil {
		return
	}

	return
}
