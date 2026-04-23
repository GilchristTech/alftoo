package main

import (
	"fmt"
	"os/exec"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AlftooApp struct {
	window   *sdl.Window
	renderer *sdl.Renderer

	command_text     string // The editing buffer of the command
	run_command_text string // If set, run this command after quitting SDL

	base_w int32
	base_h int32
	margin int32

	font       *ttf.Font
	font_fp    string
	font_size  int
	font_color sdl.Color
}

func (a *AlftooApp) SetDefaults() {
	a.base_w = 800
	a.base_h = 76
	a.margin = 16

	a.font_fp = "/usr/share/fonts/liberation/LiberationSans-Regular.ttf"
	a.font_size = 48

	a.font_color = sdl.Color{
		R: 16, G: 255, B: 64, A: 255,
	}
}

func (a *AlftooApp) Run() error {
	var run_error error

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		run_error = fmt.Errorf("Failed to initialize SDL:\n%w", err)
		goto QUIT_SDL
	}

	if err := ttf.Init(); err != nil {
		run_error = fmt.Errorf("Failed to initialize TTF:\n%w", err)
		goto QUIT_TTF
	}

	if window, err := sdl.CreateWindow(
		"alftoo",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		a.base_w+2*a.margin, a.base_h+2*a.margin,
		sdl.WINDOW_BORDERLESS,
	); err != nil {
		run_error = fmt.Errorf("Failed to create window:\n%w", err)
		goto QUIT
	} else {
		a.window = window
	}

	if renderer, err := sdl.CreateRenderer(
		a.window, -1, sdl.RENDERER_ACCELERATED,
	); err != nil {
		run_error = fmt.Errorf("Failed to create renderer:\n%w", err)
		goto QUIT
	} else {
		a.renderer = renderer
	}

	if err := a.SetFontPath(a.font_fp, a.font_size); err != nil {
		run_error = fmt.Errorf("Failed to set font: %s\n%w", err)
		goto QUIT
	}

	a.Draw()
	sdl.StartTextInput()

EVENT_LOOP:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				break EVENT_LOOP

			case *sdl.TextInputEvent:
				a.HandleInputString(ev.GetText())

			case *sdl.KeyboardEvent:
				a.HandleKeyboardEvent(ev)
			}
		}
	}

QUIT:
	sdl.StopTextInput()

	if a.font != nil {
		a.font.Close()
		a.font = nil
	}

	if a.window != nil {
		a.window.Destroy()
		a.window = nil
	}

	a.renderer = nil

QUIT_TTF:
	ttf.Quit()
QUIT_SDL:
	sdl.Quit()

	if run_error != nil {
		return run_error
	}

	if a.run_command_text != "" {
		cmd := exec.Command("sh", "-c", a.run_command_text)
		return cmd.Run()
	}

	return nil
}

func (a *AlftooApp) Quit() {
	sdl.PushEvent(&sdl.QuitEvent{
		Type: sdl.QUIT,
	})
}

func (a *AlftooApp) SetFontPath(fpath string, size_px int) error {
	if a.font != nil {
		a.font.Close()
		a.font = nil
	}

	if font, err := ttf.OpenFont(fpath, size_px*3/4); err != nil {
		return fmt.Errorf("Failed to load font at path: %s\n%w", err)
	} else {
		a.font = font
	}

	return nil
}

func (a *AlftooApp) DrawText(x, y int32, text string) error {
	texture, surface, err := a.RenderText(text)

	if texture != nil {
		defer texture.Destroy()
	}

	if surface != nil {
		defer surface.Free()
	}

	if err != nil {
		return fmt.Errorf("AlftooApp.DrawText:", err)
	}

	a.renderer.Copy(
		texture, nil,
		&sdl.Rect{
			X: x, Y: y,
			W: surface.W,
			H: surface.H,
		},
	)

	return nil
}

func (a *AlftooApp) RenderText(text string) (
	texture *sdl.Texture,
	surface *sdl.Surface,
	err error,
) {
	if surface, err = a.font.RenderUTF8Blended(
		a.command_text,
		a.font_color,
	); err != nil {
		return
	}

	if texture, err = a.renderer.CreateTextureFromSurface(surface); err != nil {
		return
	}

	return
}

func (a *AlftooApp) RenderTextWrapped(text string, wrap_length_px int) (
	texture *sdl.Texture,
	surface *sdl.Surface,
	err error,
) {
	if surface, err = a.font.RenderUTF8BlendedWrapped(
		a.command_text,
		a.font_color,
		wrap_length_px,
	); err != nil {
		return
	}

	if texture, err = a.renderer.CreateTextureFromSurface(surface); err != nil {
		return
	}

	return
}

func (a *AlftooApp) Draw() {
	var (
		err        error
		W, _       int32 = a.window.GetSize()
		margin     int   = int(a.margin)
		text_width int   = int(W) - 2*margin
		t_texture  *sdl.Texture
		t_surface  *sdl.Surface
	)

	t_texture, t_surface, err = a.RenderTextWrapped(a.command_text, text_width)

	if t_surface != nil {
		t_surface.Free()
	}

	if t_texture != nil {
		defer t_texture.Destroy()
	}

	if err == nil {
		width := a.base_w
		if t_surface.W > a.base_w {
			width = t_surface.W
		}

		height := a.base_h
		if t_surface.H > a.base_h {
			height = t_surface.H
		}

		var (
			window_width  = width + int32(2*margin)
			window_height = height + int32(2*margin)
		)
		a.ResizeWindow(window_width, window_height)

		a.renderer.SetDrawColor(30, 30, 30, 255)
		a.renderer.Clear()

		a.renderer.Copy(
			t_texture, nil,
			&sdl.Rect{
				X: int32(margin),
				Y: window_height/2 - t_surface.H/2,
				W: t_surface.W,
				H: t_surface.H,
			},
		)

	} else {
		a.renderer.SetDrawColor(30, 30, 30, 255)
		a.renderer.Clear()
	}

	a.renderer.Present()
}

func (a *AlftooApp) ResizeWindow(w, h int32) {
	var (
		original_w, original_h int32 = a.window.GetSize()
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
		a.window.SetSize(w, h)
		a.CenterWindow()
	}
}

func (a *AlftooApp) CenterWindow() error {
	if display_bounds, err := sdl.GetDisplayBounds(0); err != nil {
		return fmt.Errorf("AlftooApp.CenterWindow():\n%w", err)
	} else {
		var W, H int32 = a.window.GetSize()
		a.window.SetPosition(display_bounds.W/2-W/2, display_bounds.H/2-H/2)
		return nil
	}
}

func (a *AlftooApp) HandleInputString(input string) {
	a.SetCommandText(a.command_text + input)
}

func (a *AlftooApp) SetCommandText(t string) {
	a.command_text = t
	a.Draw()
}

func (a *AlftooApp) HandleKeyboardEvent(ev *sdl.KeyboardEvent) {
	if ev.Type != sdl.KEYDOWN {
		return
	}

	if ev.Keysym.Mod&sdl.KMOD_CTRL != 0 {
		a.HandleCTRLKeydownEvent(ev)
	}

	switch ev.Keysym.Sym {
	case sdl.K_ESCAPE:
		a.Quit()

	case sdl.K_BACKSPACE:
		if l := len(a.command_text); l > 0 {
			a.SetCommandText(a.command_text[:l-1])
		}

	case sdl.K_RETURN:
		a.run_command_text = a.command_text
		a.Quit()
	}
}

func (a *AlftooApp) HandleCTRLKeydownEvent(ev *sdl.KeyboardEvent) {
	switch ev.Keysym.Sym {
	case sdl.K_u:
		a.SetCommandText("")

	case sdl.K_v:
		if text, err := sdl.GetClipboardText(); err == nil {
			a.HandleInputString(text)
		}

	case sdl.K_d:
		a.Quit()
	}
}
