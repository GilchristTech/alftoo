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
	AddColonCommand(":base-w", (*CommandWindowGeometry)(&alf_base_w))
	AddColonCommand(":base-h", (*CommandWindowGeometry)(&alf_base_h))
	AddColonCommand(":margin", (*CommandWindowGeometry)(&alf_margin))
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

func ResizeWindow(w, h int32) {
	var (
		original_w, original_h int32 = alf_window.GetSize()
	)

	if w == -1 {
		w = original_w
	}

	if h == -1 {
		h = original_h
	}

	if w == original_w && h == original_h {
		return
	} else {
		alf_window.SetSize(w, h)
		CenterWindow()
	}
}

func CenterWindow() error {
	if display_bounds, err := sdl.GetDisplayBounds(0); err != nil {
		return fmt.Errorf("alftoo.CenterWindow\n - %w", err)

	} else {
		var W, H int32 = alf_window.GetSize()
		alf_window.SetPosition(display_bounds.W/2-W/2, display_bounds.H/2-H/2)
		return nil
	}
}

type CommandWindowGeometry int32

func (bw *CommandWindowGeometry) Run(name, args string) error {
	var (
		err_h string = "alftoo.CommandWindowGeometry.Run" + name
		value int32
	)

	if args == "" {
		return fmt.Errorf(`%s: no value specified`, err_h)

	} else if value_uint64, err := strconv.ParseUint(args, 10, 32); err != nil {
		return fmt.Errorf(`%s: argument is not a positive integer`, err_h)

	} else {
		value = int32(value_uint64)
	}

	*bw = CommandWindowGeometry(value)

	Draw()

	return nil
}
