package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	alf_window   *sdl.Window
	alf_renderer *sdl.Renderer

	alf_base_w int32 = 800
	alf_base_h int32 = 72
	alf_margin int32 = 16
)

func init() {
	var cmd CommandWindowGeometry
	AddColonCommand(":base-w", &cmd)
	AddColonCommand(":base-h", &cmd)
	AddColonCommand(":margin", &cmd)
}

func Window() *sdl.Window {
	return alf_window
}

func Renderer() *sdl.Renderer {
	return alf_renderer
}

func BaseW() int {
	return int(alf_base_w)
}

func BaseH() int {
	return int(alf_base_h)
}

func Margin() int {
	return int(alf_margin)
}

type CommandWindowGeometry int

func (bw *CommandWindowGeometry) Run(name, args string) error {
	var (
		err_h string = "alftoo.CommandWindowGeometry.Run" + name
		value int32
	)

	if args == "" {
		return fmt.Errorf(`%s: no value specified`)

	} else if value_uint64, err := strconv.ParseUint(args, 10, 32); err != nil {
		return fmt.Errorf(`%s: argument is not a positive integer`)

	} else {
		value = int32(value_uint64)
	}

	switch name {
	case ":base-w":
		alf_base_w = value
	case ":base-h":
		alf_base_h = value
	case ":margin":
		alf_margin = value
	default:
		return fmt.Errorf(`%s: command not recognized: "%s"`, err_h, name)
	}

	Draw()

	return nil
}
