package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "graph_roller")
	defer rl.CloseWindow()

	// graph_roller: controls initialization
	nodesSelectorScrollView := rl.Rectangle{X: 0, Y: 0, Width: 0, Height: 0}
	nodesSelectorScrollOffset := rl.Vector2{X: 0, Y: 0}
	nodesSelectorBoundsOffset := rl.Vector2{X: 0, Y: 0}

	rl.SetTargetFPS(60)

	graph_field := rl.LoadRenderTexture(638, 406)
	defer rl.UnloadRenderTexture(graph_field)

	field_camera := rl.Camera2D{
		Offset:   rl.Vector2{},
		Target:   rl.Vector2{},
		Rotation: 0,
		Zoom:     1,
	}

	for !rl.WindowShouldClose() {
		// Update
		// TODO: Implement required update logic
		mouse_pos := rl.GetMousePosition()

		if rl.CheckCollisionPointRec(mouse_pos, rl.Rectangle{
			X:      241,
			Y:      25,
			Width:  638,
			Height: 406,
		}) {
			// if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				delta := rl.GetMouseDelta()
				delta = rl.Vector2Scale(delta, -1.0/field_camera.Zoom)

				field_camera.Target = rl.Vector2Add(field_camera.Target, delta)
			}

			wheel := rl.GetMouseWheelMove()
			if wheel != 0 {
				// get the world point that is under the mouse
				mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), field_camera)

				// set the offset to where the mouse is
				field_camera.Offset = rl.GetMousePosition()

				// set the target to match, so that the camera maps the world space point under the cursor to the screen space point under the cursor at any zoom
				field_camera.Target = mouseWorldPos

				// zoom
				field_camera.Zoom += wheel * 0.125
				if field_camera.Zoom < 0.125 {
					field_camera.Zoom = 0.125
				}
			}

			if rl.IsKeyPressed(rl.KeyR) {
				field_camera.Zoom = 1
				field_camera.Target = rl.Vector2{X: 0, Y: 0}
				field_camera.Offset = rl.Vector2{
					X: 0,
					Y: 0,
				}
			}
		}

		// Draw

		rl.BeginTextureMode(graph_field)
		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		rl.BeginMode2D(field_camera)

		rl.DrawRectangle(0, 0, 100, 100, rl.Red)

		rl.EndMode2D()
		rl.EndTextureMode()

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		// raygui: controls drawing
		gui.ScrollPanel(
			rl.Rectangle{X: 0, Y: 24, Width: 240 - nodesSelectorBoundsOffset.X, Height: 408 - nodesSelectorBoundsOffset.Y},
			"scroll panel",
			rl.Rectangle{X: 0, Y: 24, Width: 240, Height: 408},
			&nodesSelectorScrollOffset,
			&nodesSelectorScrollView,
		)
		gui.Label(rl.Rectangle{X: 0, Y: 0, Width: 240, Height: 24}, "Graph Roller")

		gui.Panel(rl.Rectangle{X: 240, Y: 24, Width: 640, Height: 408}, "")
		rl.DrawTextureRec(
			graph_field.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  638,
				Height: -406,
			},
			rl.Vector2{
				X: 241,
				Y: 25,
			},
			rl.White,
		)

		gui.GroupBox(rl.Rectangle{X: 0, Y: 440, Width: 880, Height: 40}, "Controls")

		rl.EndDrawing()
	}

}
