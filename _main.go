package main

import (
	"fmt"
	"image/color"
	"math"

	// gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH_FACTOR  = 16
	WINDOW_HEIGHT_FACTOR = 9
	FIELD_WIDTH_FACTOR   = 16
	FIELD_HEIGHT_FACTOR  = 9

	WINDOW_FRAGMENT = 90

	WINDOW_HEIGHT = WINDOW_HEIGHT_FACTOR * WINDOW_FRAGMENT
	WINDOW_WIDTH  = WINDOW_WIDTH_FACTOR * WINDOW_FRAGMENT

	TARGET_FPS = 60

	ARMS_COUNT = 5
)

var (
	DEBUG        = true
	BACKGROUND   = color.RGBA{0x18, 0x18, 0x18, 0xFF}
	IOSEVKA_FONT rl.Font
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "test")
	defer rl.CloseWindow()

	IOSEVKA_FONT = rl.LoadFont("src/IosevkaNerdFont-Regular.ttf")

	// rl.ToggleFullscreen()

	rl.SetTargetFPS(TARGET_FPS)

	main_texture := rl.LoadRenderTexture(WINDOW_WIDTH, WINDOW_HEIGHT)
	defer rl.UnloadRenderTexture(main_texture)

	graph_field := load_field()
	defer graph_field.unload()

	for !rl.WindowShouldClose() {
		// Events processing
		// ------------------------------------------------------------------------
		graph_field.process_events()

		if rl.IsKeyPressed(rl.KeyD) {
			DEBUG = !DEBUG
		}
		// ------------------------------------------------------------------------
		// !Events processing

		// Field Texture Drawing
		// ------------------------------------------------------------------------
		rl.BeginTextureMode(graph_field.texture)

		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(graph_field.camera)

		// for x := 0; x <= width; x += WINDOW_FRAGMENT * 2 {
		// 	for y := 0; y <= height; y += WINDOW_FRAGMENT * 4 {
		// 		snowflake(rl.Vector2{X: float32(x), Y: float32(y)}, WINDOW_FRAGMENT/2, 3, rl.White, 4)
		// 	}
		// }
		//
		// for x := WINDOW_FRAGMENT; x <= width; x += WINDOW_FRAGMENT * 2 {
		// 	for y := WINDOW_FRAGMENT * 2; y <= height; y += WINDOW_FRAGMENT * 4 {
		// 		snowflake(rl.Vector2{X: float32(x), Y: float32(y)}, WINDOW_FRAGMENT/2, 3, rl.White, 4)
		// 	}
		// }
		rl.EndMode2D()

		rl.EndTextureMode()
		// ------------------------------------------------------------------------
		// !Field Texture Drawing

		// Main Texture Drawing
		// ------------------------------------------------------------------------
		rl.BeginTextureMode(main_texture)
		rl.ClearBackground(BACKGROUND)

		rl.DrawTexturePro(
			graph_field.texture.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  graph_field.dimensions.Width,
				Height: graph_field.dimensions.Height,
			},
			graph_field.dimensions,
			rl.Vector2{X: 0, Y: 0},
			0,
			rl.White,
		)
		// rl.DrawTextureRec(
		// 	graph_field.texture.Texture,
		// 	rl.Rectangle{
		// 		X:      0,
		// 		Y:      0,
		// 		Width:  graph_field.dimensions.Width,
		// 		Height: graph_field.dimensions.Height,
		// 	},
		// 	rl.Vector2{
		// 		X: graph_field.dimensions.X,
		// 		Y: graph_field.dimensions.Y,
		// 	},
		// 	rl.White,
		// )

		if DEBUG {
			draw_debug_info()
		}

		rl.EndTextureMode()
		// ------------------------------------------------------------------------
		// !Main Texture Drawing

		// Window drawing
		// ------------------------------------------------------------------------
		rl.BeginDrawing()

		width := rl.GetRenderWidth()
		height := rl.GetRenderHeight()
		rl.DrawTextureRec(
			main_texture.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  float32(width),
				Height: -float32(height),
			},
			rl.Vector2{
				X: 0,
				Y: 0,
			},
			rl.White,
		)

		rl.EndDrawing()
		// !Window drawing
		// ------------------------------------------------------------------------
	}
}

func draw_debug_info() {
	display := rl.GetCurrentMonitor()
	width := rl.GetMonitorWidth(display)
	height := rl.GetMonitorHeight(display)

	msg := "DEBUG:\n"
	msg += fmt.Sprintf("- Dimensions: {width: %v, height: %v}\n", width, height)
	msg += fmt.Sprintf("- Mouse:\n")
	msg += fmt.Sprintf("  - On Screen: %#v\n", rl.GetMousePosition())

	rl.DrawTextEx(IOSEVKA_FONT, msg, rl.Vector2{X: float32(width) * 0.2, Y: float32(height) * 0.8}, 32, 0, rl.Red)
}

type field struct {
	dimensions rl.Rectangle
	camera     rl.Camera2D
	texture    rl.RenderTexture2D
}

func load_field() field {
	dimensions := rl.Rectangle{
		Width:  WINDOW_FRAGMENT * FIELD_WIDTH_FACTOR,
		Height: WINDOW_FRAGMENT * FIELD_HEIGHT_FACTOR,
		X:      WINDOW_FRAGMENT * (WINDOW_WIDTH_FACTOR - FIELD_WIDTH_FACTOR),
		Y:      0,
	}

	texture := rl.LoadRenderTexture(int32(dimensions.Width), int32(dimensions.Height))

	camera := rl.Camera2D{
		Offset:   rl.Vector2{},
		Target:   rl.Vector2{},
		Rotation: 0,
		Zoom:     1,
	}

	return field{
		dimensions: dimensions,
		camera:     camera,
		texture:    texture,
	}
}

func (self field) unload() {
	rl.UnloadRenderTexture(self.texture)
}
func (self field) process_events() {
	mouse_pos := rl.GetMousePosition()

	if !rl.CheckCollisionPointRec(mouse_pos, self.dimensions) {
		return
	}
	// if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1.0/self.camera.Zoom)

		self.camera.Target = rl.Vector2Add(self.camera.Target, delta)
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		// get the world point that is under the mouse
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), self.camera)

		// set the offset to where the mouse is
		self.camera.Offset = rl.GetMousePosition()

		// set the target to match, so that the camera maps the world space point under the cursor to the screen space point under the cursor at any zoom
		self.camera.Target = mouseWorldPos

		// zoom
		self.camera.Zoom += wheel * 0.125
		if self.camera.Zoom < 0.125 {
			self.camera.Zoom = 0.125
		}
	}

	if rl.IsKeyPressed(rl.KeyR) {
		self.camera.Zoom = 1
		self.camera.Target = rl.Vector2{X: 0, Y: 0}
		self.camera.Offset = rl.Vector2{
			X: 0,
			Y: 0,
		}
	}
}

func snowflake(center rl.Vector2, length float32, thickness float32, color color.RGBA, limit int) {
	if limit == 0 {
		return
	}

	var angle float64 = 2 * math.Pi / ARMS_COUNT

	for i := range ARMS_COUNT {
		child_center := rl.Vector2{
			X: center.X + length*float32(math.Cos(float64(i)*angle)),
			Y: center.Y + length*float32(math.Sin(float64(i)*angle)),
		}
		rl.DrawLineEx(center, child_center, thickness, color)

		snowflake(child_center, length*0.45, thickness*0.8, color, limit-1)
	}
}
