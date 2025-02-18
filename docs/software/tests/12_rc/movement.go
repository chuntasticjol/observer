package main

import (
	"image"
	"log"
	"math"
)

var currentX = 90
var currentY = 90

var pxsPerDegree = math.Sqrt(float64(camWidth*camWidth)+float64(camHeight*camHeight)) / angleOfViewDig

func aimTarget(coor image.Point) {

	angleX := float64(coor.X-midPoint.X) / pxsPerDegree
	moveCam("axisX", angleX)

	angleY := float64(coor.Y-midPoint.Y) / pxsPerDegree
	moveCam("axisY", angleY)
}

func moveCam(direct string, angle float64) {

	switch direct {

	case "axisX":

		angle *= invertX
		angle *= calibrateX

		switch deltaX := currentX + int(math.Round(angle)); {
		case deltaX > 180:
			_ = servoX.Max()

			currentX = 180
		case deltaX < 0:
			_ = servoY.Min()

			currentX = 0
		default:
			_ = servoX.Move(uint8(deltaX))
			currentX = deltaX
		}

	case "axisY":

		angle *= invertY
		angle *= calibrateY

		switch deltaY := currentY + int(math.Round(angle)); {
		case deltaY > 180:
			_ = servoY.Max()

			currentY = 180
		case deltaY < 0:
			_ = servoY.Min()

			currentY = 0
		default:
			_ = servoY.Move(uint8(deltaY))
			currentY = deltaY
		}
	}
}

func calibrateServos() {
	log.Printf("Calibrating servomotors ...\n")
	centerServos()
	_ = servoX.Min()
	_ = servoY.Min()
	_ = servoX.Max()
	_ = servoY.Max()
	centerServos()
}

func centerServos() {
	_ = servoX.Center()
	_ = servoY.Center()
	currentX = 90
	currentY = 90
}
