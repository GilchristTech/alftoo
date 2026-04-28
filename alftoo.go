package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	alf_window   *sdl.Window
	alf_renderer *sdl.Renderer

	alf_command_text     string // The editing buffer of the command
	alf_run_command_text string // If set, run this command after quitting SDL

	alf_base_w int32
	alf_base_h int32
	alf_margin int32

	alf_font       *ttf.Font
	alf_font_name  string
	alf_font_fp    string
	alf_font_size  int
	alf_font_color sdl.Color

	alf_background_color sdl.Color
)

func SetDefaults() error {
	alf_base_w = 800
	alf_base_h = 76
	alf_margin = 16

	alf_font_name = "Sans"

	if fp, err := FontFindPath(alf_font_name); err != nil {
		return fmt.Errorf("alftoo.SetDefaults: could not find font\n - %w", err)
	} else {
		alf_font_fp = fp
	}

	alf_font_size = 48

	alf_font_color = sdl.Color{
		R: 16, G: 255, B: 64, A: 255,
	}

	alf_background_color = sdl.Color{
		R: 30, G: 30, B: 30, A: 255,
	}

	return nil
}

func Run() error {
	var run_error error

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		run_error = fmt.Errorf("alftoo.Run: Failed to initialize SDL\n - %w", err)
		goto QUIT_SDL
	}

	if err := ttf.Init(); err != nil {
		run_error = fmt.Errorf("alftoo.Run: Failed to initialize TTF\n - %w", err)
		goto QUIT_TTF
	}

	if window, err := sdl.CreateWindow(
		"alftoo",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		alf_base_w+2*alf_margin, alf_base_h+2*alf_margin,
		sdl.WINDOW_BORDERLESS,
	); err != nil {
		run_error = fmt.Errorf("alftoo.Run: Failed to create window\n - %w", err)
		goto QUIT
	} else {
		alf_window = window
	}

	if renderer, err := sdl.CreateRenderer(
		alf_window, -1, sdl.RENDERER_ACCELERATED,
	); err != nil {
		run_error = fmt.Errorf("alftoo.Run: Failed to create renderer\n - %w", err)
		goto QUIT
	} else {
		alf_renderer = renderer
	}

	if err := SetFontPath(alf_font_fp, alf_font_size); err != nil {
		run_error = fmt.Errorf("alftoo.Run: Failed to set font\n - %w", err)
		goto QUIT
	}

	Draw()
	sdl.StartTextInput()

EVENT_LOOP:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				break EVENT_LOOP

			case *sdl.TextInputEvent:
				HandleInputString(ev.GetText())

			case *sdl.KeyboardEvent:
				HandleKeyboardEvent(ev)
			}
		}
	}

QUIT:
	sdl.StopTextInput()

	if alf_font != nil {
		alf_font.Close()
		alf_font = nil
	}

	if alf_window != nil {
		alf_window.Destroy()
		alf_window = nil
	}

	alf_renderer = nil

QUIT_TTF:
	ttf.Quit()
QUIT_SDL:
	sdl.Quit()

	if run_error != nil {
		return run_error
	}

	if c := alf_run_command_text; c != "" {
		err := RunCommand(alf_run_command_text)
		alf_run_command_text = ""
		return err
	}

	return nil
}

func Quit() {
	sdl.PushEvent(&sdl.QuitEvent{
		Type: sdl.QUIT,
	})
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

func DrawText(x, y int32, text string) error {
	texture, surface, err := RenderText(text)

	if texture != nil {
		defer texture.Destroy()
	}

	if surface != nil {
		defer surface.Free()
	}

	if err != nil {
		return fmt.Errorf("alftoo.DrawText\n - %w", err)
	}

	alf_renderer.Copy(
		texture, nil,
		&sdl.Rect{
			X: x, Y: y,
			W: surface.W,
			H: surface.H,
		},
	)

	return nil
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

func Draw() {
	var (
		err        error
		W, _       int32 = alf_window.GetSize()
		margin     int   = int(alf_margin)
		text_width int   = int(W) - 2*margin
		t_texture  *sdl.Texture
		t_surface  *sdl.Surface
	)

	t_texture, t_surface, err = RenderTextWrapped(alf_command_text, text_width)

	if t_surface != nil {
		t_surface.Free()
	}

	if t_texture != nil {
		defer t_texture.Destroy()
	}

	if err == nil {
		width := alf_base_w
		if t_surface.W > alf_base_w {
			width = t_surface.W
		}

		height := alf_base_h
		if t_surface.H > alf_base_h {
			height = t_surface.H
		}

		var (
			window_width  = width + int32(2*alf_margin)
			window_height = height + int32(2*alf_margin)
		)
		ResizeWindow(window_width, window_height)

		alf_renderer.SetDrawColor(
			alf_background_color.R,
			alf_background_color.G,
			alf_background_color.B,
			alf_background_color.A,
		)

		alf_renderer.Clear()

		alf_renderer.Copy(
			t_texture, nil,
			&sdl.Rect{
				X: int32(alf_margin),
				Y: window_height/2 - t_surface.H/2,
				W: t_surface.W,
				H: t_surface.H,
			},
		)

	} else {
		alf_renderer.SetDrawColor(
			alf_background_color.R,
			alf_background_color.G,
			alf_background_color.B,
			alf_background_color.A,
		)

		alf_renderer.Clear()
	}

	alf_renderer.Present()
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

func HandleInputString(input string) {
	SetCommandText(alf_command_text + input)
}

func CommandText() string {
	return alf_command_text
}

func SetCommandText(t string) {
	alf_command_text = t
	Draw()
}

func HandleKeyboardEvent(ev *sdl.KeyboardEvent) {
	if ev.Type != sdl.KEYDOWN {
		return
	}

	if ev.Keysym.Mod&sdl.KMOD_CTRL != 0 {
		HandleCTRLKeydownEvent(ev)
	}

	switch ev.Keysym.Sym {
	case sdl.K_ESCAPE:
		Quit()

	case sdl.K_BACKSPACE:
		if l := len(alf_command_text); l > 0 {
			SetCommandText(alf_command_text[:l-1])
		}

	case sdl.K_RETURN:
		command_name := getWord(alf_command_text, 0)
		command_args := alf_command_text[len(command_name):]

		if strings.HasPrefix(command_name, ":") {
			if colon_command := colon_commands[command_name]; colon_command == nil {
				fmt.Fprintf(
					os.Stderr,
					`alftoo.HandleKeyboardEvent: Colon command "%s" does not exist\n`,
					command_name,
				)

			} else if err := colon_command.Run(command_name, command_args); err != nil {
				fmt.Fprintf(
					os.Stderr,
					"alftoo.HandleKeyboardEvent: Error running colon command\n - %v\n",
					err,
				)
			} else {
				SetCommandText("")
			}

		} else {
			alf_run_command_text = strings.TrimSpace(alf_command_text)
			Quit()
		}
	}
}

func HandleCTRLKeydownEvent(ev *sdl.KeyboardEvent) {
	switch ev.Keysym.Sym {
	case sdl.K_u:
		SetCommandText("")

	case sdl.K_v:
		if text, err := sdl.GetClipboardText(); err == nil {
			HandleInputString(text)
		}

	case sdl.K_d:
		Quit()
	}
}

func RunCommand(c string) error {
	cmd := exec.Command("sh", "-c", alf_run_command_text)
	return cmd.Run()
}
