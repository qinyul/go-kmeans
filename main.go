package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
)

const (
	SAMPLE_SIZE = 0.5
	CAMERA_SPEED = SAMPLE_SIZE
	PHI = 3.14
)

type Sample3D struct {
	items *[]rl.Vector3
	count int
	capacity int
}

func generateCluster(center rl.Vector3,radius float32, count int, samples *Sample3D ) {
	r := rand.New(rand.NewSource(99))
	for i := 0; i < count; i++ {
		rad := r.Float32() * radius
		theta := r.Float32() * 2 * PHI
		phi := r.Float32() * 2 * PHI

		var sample rl.Vector3
		sample = rl.Vector3{
			X:	float32(math.Sin(float64(theta))) * float32(math.Cos(float64(phi))) * rad,
			Y:	float32(math.Sin(float64(theta))) * float32(math.Sin(float64(phi))) * rad,
			Z:	float32(math.Cos(float64(theta))) * rad,
		}
		*samples.items = append(*samples.items, rl.Vector3Add(sample,center))
	}
}

func main () {
	var set Sample3D
	clusterCount := 200
	clusterRadius := float32(20)

	initValue := []rl.Vector3{rl.Vector3Zero()}
	set = Sample3D{
		&initValue,
		clusterCount,
		int(10),
	}

	generateCluster(rl.Vector3{0,0,0},clusterRadius,set.count, &set)
	generateCluster(rl.Vector3{-clusterRadius,clusterRadius,0},clusterRadius / 2,set.count / 2, &set)
	generateCluster(rl.Vector3{clusterRadius,clusterRadius,0},clusterRadius / 2,set.count / 2, &set)
	
	cameraRad := float32(50)
	cameraTheta := float32(0.0)
	cameraPhi := float32(0.0)
	cameraVel := float32(0.0)

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(500,500,"kmeans")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose(){
		dt := float32(rl.GetFrameTime())
		cameraRad += cameraVel * dt
		if cameraRad < float32(0.0) {
			cameraRad = float32(0.0)
		}
		cameraVel -= rl.GetMouseWheelMove() * float32(20)
		cameraVel *= 0.9

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			delta := rl.GetMouseDelta()
			cameraTheta -= delta.X * 0.01
			cameraPhi -= delta.Y * 0.01
		}

		rl.BeginDrawing()
		camera := rl.Camera3D{
			Position:	rl.Vector3{
				X:	float32(math.Sin(float64(cameraTheta))) * float32(math.Cos(float64(cameraPhi))) * cameraRad,
				Y:	float32(math.Sin(float64(cameraTheta))) * float32(math.Sin(float64(cameraPhi))) * cameraRad,
				Z:	float32(math.Cos(float64(cameraTheta))) * cameraRad,
			},
			Target:		rl.Vector3{0,0,0},
			Up:			rl.Vector3{0,1,0},
			Fovy:		90.0,
			Projection:	rl.CameraPerspective,
		}

		if rl.IsKeyDown(rl.KeyW){
			camera.Position.Z -= CAMERA_SPEED * dt
		}
		if rl.IsKeyDown(rl.KeyS){
			camera.Position.Z += CAMERA_SPEED * dt
		}
		if rl.IsKeyDown(rl.KeyD){
			camera.Position.X += CAMERA_SPEED * dt
		}
		if rl.IsKeyDown(rl.KeyA){
			camera.Position.X -= CAMERA_SPEED * dt
		}

		rl.ClearBackground(rl.NewColor(18,18,18,18))
		rl.BeginMode3D(camera)

		for i := 0; i < len(*(set.items)); i++ {
			it := (*set.items)[i]
			rl.DrawSphere(
				it,
				SAMPLE_SIZE,
				rl.Red,
			)
		}

		rl.EndMode3D()
		rl.EndDrawing()
	}
}
