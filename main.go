package main

import (
	"fmt"
	"math"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	field_graph := graph{}

	screenWidth := int32(1280)
	screenHeight := int32(688)

	rl.InitWindow(screenWidth, screenHeight, "graph_roller")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	field_texture := rl.LoadTexture("src/field-texture.png")
	defer rl.UnloadTexture(field_texture)

	// graph_roller: controls initialization
	var instrumentsToggleGroupActive int32 = 0
	nodeSelectorScrollView := rl.Rectangle{X: 0, Y: 0, Width: 0, Height: 0}
	nodeSelectorScrollOffset := rl.Vector2{X: 0, Y: 0}
	nodeSelectorBoundsOffset := rl.Vector2{X: 0, Y: 0}
	tempPanelScrollView := rl.Rectangle{X: 0, Y: 0, Width: 0, Height: 0}
	tempPanelScrollOffset := rl.Vector2{X: 0, Y: 0}
	tempPanelBoundsOffset := rl.Vector2{X: 0, Y: 0}

	layoutRecs := map[string]rl.Rectangle{
		"main_panel":               {X: 0, Y: 0, Width: 1280, Height: 720},
		"graph_field":              {X: 256, Y: 0, Width: 1024, Height: 680},
		"graph_texture":            {X: 257, Y: 1, Width: 1022, Height: -678},
		"instruments_toggle_group": {X: 264, Y: 688, Width: 24, Height: 24},
		"node_selector":            {X: 0, Y: 24, Width: 256, Height: 344},
		"temp_panel":               {X: 0, Y: 368, Width: 256, Height: 352},
	}

	graph_field := rl.LoadRenderTexture(int32(layoutRecs["graph_field"].Width)-2, int32(layoutRecs["graph_field"].Height)-2)
	defer rl.UnloadRenderTexture(graph_field)

	fieldWidth := float32(5000)
	fieldHeight := float32(5000)

	camera_initial_target := rl.Vector2{X: fieldWidth / 2, Y: fieldHeight / 2}
	camera_initial_offset := rl.Vector2{X: float32(screenWidth) / 2, Y: float32(screenHeight) / 2}

	fieldCamera := rl.Camera2D{
		Target:   camera_initial_target,
		Offset:   camera_initial_offset,
		Rotation: 0,
		Zoom:     1,
	}

	for !rl.WindowShouldClose() {
		// Update logic
		rl.SetMouseCursor(rl.MouseCursorDefault)

		var (
			mousePos           rl.Vector2 = rl.GetMousePosition()
			mousePosOnFieldInt rl.Vector2 = GetMousePositionOnField(mousePos, layoutRecs, fieldCamera)

			frameTime float32 = rl.GetFrameTime()

			scalePivot rl.Vector2 = mousePosOnFieldInt
			deltaScale float32    = 0

			wheelMoved = false
		)

		if rl.CheckCollisionPointRec(mousePos, layoutRecs["graph_field"]) {
			if mouseWheelMove := rl.GetMouseWheelMoveV().Y; mouseWheelMove != 0 {
				wheelMoved = true

				deltaScale += mouseWheelMove
			}

			if math.Abs(float64(deltaScale)) > 0.5 {
				p0 := rl.Vector2{
					X: (scalePivot.X - fieldCamera.Target.X) / fieldCamera.Zoom,
					Y: (scalePivot.Y - fieldCamera.Target.Y) / fieldCamera.Zoom,
				}

				newZoom := fieldCamera.Zoom + deltaScale*frameTime
				if newZoom < 0.5 {
					newZoom = 0.5
				} else if newZoom > 3.0 {
					newZoom = 3.0
				}
				fieldCamera.Zoom = newZoom

				p1 := rl.Vector2{
					X: (scalePivot.X - fieldCamera.Target.X) / fieldCamera.Zoom,
					Y: (scalePivot.Y - fieldCamera.Target.Y) / fieldCamera.Zoom,
				}

				fieldCamera.Target = rl.Vector2Add(fieldCamera.Target, rl.Vector2Subtract(p0, p1))
			}

			if rl.IsMouseButtonDown(rl.MouseButtonMiddle) || wheelMoved {
				delta := rl.GetMouseDelta()

				if !wheelMoved {
					rl.SetMouseCursor(rl.MouseCursorResizeAll)
					delta = rl.Vector2Scale(delta, -1.0/fieldCamera.Zoom)
				} else {
					delta = rl.Vector2Scale(delta, 1.0/fieldCamera.Zoom)
				}

				newTarget := rl.Vector2Add(fieldCamera.Target, delta)

				halfScreenWidth := (float32(screenWidth) / fieldCamera.Zoom) / 2
				halfScreenHeight := (float32(screenHeight) / fieldCamera.Zoom) / 2

				newTarget.X = rl.Clamp(newTarget.X, halfScreenWidth, fieldWidth-halfScreenWidth)
				newTarget.Y = rl.Clamp(newTarget.Y, halfScreenHeight, fieldHeight-halfScreenHeight)

				fieldCamera.Target = newTarget
			}

			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && instrumentsToggleGroupActive == 0 {
				mouse_pos := rl.GetMousePosition()
				offset := rl.Vector2{
					X: layoutRecs["graph_field"].X,
					Y: layoutRecs["graph_field"].Y,
				}

				on_field_pos := rl.GetScreenToWorld2D(rl.Vector2Subtract(mouse_pos, offset), fieldCamera)

				field_graph.addVertex(on_field_pos)
			}

			if rl.IsKeyPressed(rl.KeyR) {
				fieldCamera.Target = camera_initial_target
				fieldCamera.Offset = camera_initial_offset
				fieldCamera.Zoom = 1
				fieldCamera.Rotation = 0
			}

		}

		// Draw
		rl.BeginTextureMode(graph_field)
		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		rl.BeginMode2D(fieldCamera)

		for x := float32(0); x < fieldWidth; x += float32(field_texture.Width) {
			for y := float32(0); y < fieldHeight; y += float32(field_texture.Height) {
				rl.DrawTextureV(field_texture, rl.Vector2{X: x, Y: y}, rl.White)
			}
		}

		for _, vertex := range field_graph.vertices {
			rl.DrawCircleV(vertex.position, 32, rl.Brown)
			rl.DrawText(vertex.name, int32(vertex.position.X-5), int32(vertex.position.Y-8), 16, rl.Black)
		}

		rl.EndMode2D()

		rl.DrawText("Use middle mouse button to pan", 10, 10, 20, rl.Red)
		rl.DrawText("Use mouse wheel to zoom", 10, 40, 20, rl.Red)
		rl.DrawText(fmt.Sprintf("Zoom: %.2f", fieldCamera.Zoom), 10, 70, 20, rl.Red)

		rl.EndTextureMode()

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		// raygui: controls drawing
		gui.Panel(layoutRecs["main_panel"], "Graph Roller")
		gui.Panel(layoutRecs["graph_field"], "")

		rl.DrawTextureRec(
			graph_field.Texture,
			rl.Rectangle{X: 0, Y: 0, Width: 1022, Height: -678},
			rl.Vector2{
				X: 257,
				Y: 1,
			},
			rl.White,
		)

		instrumentsToggleGroupActive = gui.ToggleGroup(
			layoutRecs["instruments_toggle_group"],
			"A;D;E",
			instrumentsToggleGroupActive,
		)
		gui.ScrollPanel(
			rl.Rectangle{
				X:      layoutRecs["node_selector"].X,
				Y:      layoutRecs["node_selector"].Y,
				Width:  layoutRecs["node_selector"].Width - nodeSelectorBoundsOffset.X,
				Height: layoutRecs["node_selector"].Height - nodeSelectorBoundsOffset.Y,
			},
			// "Nodes Selector",
			"",
			layoutRecs["node_selector"],
			&nodeSelectorScrollOffset,
			&nodeSelectorScrollView,
		)
		gui.ScrollPanel(
			rl.Rectangle{
				X:      layoutRecs["temp_panel"].X,
				Y:      layoutRecs["temp_panel"].Y,
				Width:  layoutRecs["temp_panel"].Width - tempPanelBoundsOffset.X,
				Height: layoutRecs["temp_panel"].Height - tempPanelBoundsOffset.Y,
			},
			// "Temp",
			"",
			layoutRecs["temp_panel"],
			&tempPanelScrollOffset,
			&tempPanelScrollView,
		)
		rl.EndDrawing()
	}
}

type graph struct {
	vertices []*vertex
}

func (self *graph) addVertex(position rl.Vector2) {
	new_vertex := vertex{
		name:     fmt.Sprint(len(self.vertices) + 1),
		position: position,
	}

	self.vertices = append(self.vertices, &new_vertex)
}

type vertex struct {
	name     string
	position rl.Vector2
}

func GetMousePositionOnField(mousePos rl.Vector2, layoutRecs map[string]rl.Rectangle, fieldCamera rl.Camera2D) rl.Vector2 {
	return rl.GetScreenToWorld2D(
		rl.Vector2Subtract(
			mousePos,
			rl.Vector2{
				X: layoutRecs["graph_field"].X,
				Y: layoutRecs["graph_field"].Y,
			}),
		fieldCamera,
	)
}
