package controller

import (
	"github.com/chutommy/observer/engine"

	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
)

// initWindow sets window driver.
func (o *Observer) initWindow() {
	if o.cfg.Show {
		o.log.Info("Initializing capturing window driver")
		o.window = opencv.NewWindowDriver()
	}
}

// initServos initializes the blaster and connected servo motors.
func (o *Observer) initServos() {
	o.blaster.Start([]int64{o.cfg.ServoX.Pin, o.cfg.ServoY.Pin})

	o.log.Info("Starting servo motors")
	servoX := engine.NewServo(o.blaster, o.cfg.ServoX)
	servoY := engine.NewServo(o.blaster, o.cfg.ServoY)
	o.servos = engine.NewServos(servoX, servoY)
}

// initCamera turns on the camera and initializes the connection.
func (o *Observer) initCamera() {
	mat := gocv.NewMat()
	o.activeFrame.Store(mat)

	o.log.Info("Turning camera on (start recording)")
	_ = o.camera.On(opencv.Frame, func(data interface{}) {
		cam := data.(gocv.Mat)
		o.activeFrame.Store(cam)
	})
}

// checkFrequency validates the value of the period and edits if necessary.
func (o *Observer) checkFrequency() {
	if 1000/o.cfg.Period > o.cfg.MaxFPS {
		i := (1000 / o.cfg.MaxFPS) + 1

		o.log.Warnf("Period length is too short, extendeding to %d", i)
		o.cfg.Period = i
	}
}
