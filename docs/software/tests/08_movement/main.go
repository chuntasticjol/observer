package main

import (
	"image"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/opencv"
	"gobot.io/x/gobot/platforms/raspi"
	"gocv.io/x/gocv"
)

var robotName = "Big Brother"
var cascade = "../haarcascades/haarcascade_frontalface_default.xml"

var cameraSource = 0
var servoXpin = "servoMotX"
var servoYpin = "servoMotY"

var camWidth = 640
var camHeight = 480
var angleOfView = 72.4

var aimArea = 0.25

var img atomic.Value

var raspiAdaptor = raspi.NewAdaptor()

var window = opencv.NewWindowDriver()
var camera = opencv.NewCameraDriver(cameraSource)
var servoX = gpio.NewServoDriver(raspiAdaptor, servoXpin)
var servoY = gpio.NewServoDriver(raspiAdaptor, servoYpin)

func main() {

	work := func() {
		mat := gocv.NewMat()
		defer func(mat *gocv.Mat) {
			_ = mat.Close()
		}(&mat)

		img.Store(mat)

		_ = camera.On(opencv.Frame, func(data interface{}) {
			i := data.(gocv.Mat)

			img.Store(i)
		})

		gobot.Every(1*time.Millisecond, func() {

			i := img.Load().(gocv.Mat)
			if i.Empty() {
				return
			}

			objects := opencv.DetectObjects(cascade, i)

			target := getNearestObject(objects)
			if (target != image.Rectangle{}) {

				lock := getCoordinates(target)

				if !lock.In(midRect) {
					aimTarget(lock)
				}

				opencv.DrawRectangles(i, []image.Rectangle{midRect}, 0, 255, 0, 5)
				window.ShowImage(i)
				window.WaitKey(1)
			}
		})
	}

	connections := []gobot.Connection{raspiAdaptor}
	devices := []gobot.Device{camera, window, servoX, servoY}

	robot := gobot.NewRobot(
		robotName,
		connections,
		devices,
		work,
	)

	_ = robot.Start()
}
