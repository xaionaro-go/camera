package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/pflag"
	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/allplatforms"
)

func main() {
	availableCameras, err := camera.ListCameras()
	if err != nil {
		panic(fmt.Errorf("unable to get the list of cameras: %w", err))
	}
	if len(availableCameras) == 0 {
		panic("no cameras found")
	}

	netPprofAddr := pflag.String("net-pprof-addr", "", "")
	widthFlag := pflag.Uint64("width", 0, "")
	fpsFlag := pflag.Float64("fps", math.NaN(), "")
	pixFmtFlag := pflag.String("pixel-format", "", "")
	platformFlag := pflag.String("platform", "", "")
	deviceFlag := pflag.String("device", availableCameras[0].DevicePath, "")
	pflag.Parse()

	if *netPprofAddr != "" {
		go func() {
			log.Println(http.ListenAndServe(*netPprofAddr, nil))
		}()
	}

	a := app.New()
	w := a.NewWindow("CameraView")
	w.Show()
	defer func() { processRecover(w, recover()) }()

	var plat camera.Platform
	var devicePath camera.DevicePath
	if *platformFlag != "" {
		plat = allplatforms.Get(*platformFlag)
		if plat == nil {
			panicInUI(w, fmt.Errorf("platform '%s' is unknown", *platformFlag))
		}
		availableCameras, err := plat.ListCameras()
		if err != nil {
			panicInUI(w, fmt.Errorf("unable to list cameras: %w", err))
		}
		for _, c := range availableCameras {
			if c == *deviceFlag {
				devicePath = c
				break
			}
		}
		if devicePath == "" {
			panicInUI(w, fmt.Errorf("camera with path '%s' is not found (available: %#+v)", *deviceFlag, availableCameras))
		}
	} else {
		var cameraSelector camera.DevicePathAndPlatform
		for _, c := range availableCameras {
			if c.DevicePath == *deviceFlag {
				cameraSelector = c
				break
			}
		}
		if cameraSelector.Platform == nil {
			panicInUI(w, fmt.Errorf("camera with path '%s' is not found (available: %#+v)", *deviceFlag, availableCameras))
		}
		plat = cameraSelector.Platform
		devicePath = cameraSelector.DevicePath
	}

	formats, err := plat.ListFormats(devicePath)
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to list the formats: %w", err))
	}
	if len(formats) == 0 {
		panicInUI(w, fmt.Errorf("the list of available formats is empty"))
	}

	var buf bytes.Buffer
	jsonEnc := json.NewEncoder(&buf)
	jsonEnc.SetIndent("", " ")
	jsonEnc.Encode(formats)
	log.Printf("available formats:\n%s", buf.Bytes())

	if *pixFmtFlag != "" {
		pixFmt := camera.PixelFormatByName(*pixFmtFlag)
		if pixFmt == camera.PixelFormatUndefined {
			panicInUI(w, fmt.Errorf("unknown pixel format name '%s'", *pixFmtFlag))
		}
		formats = formats.FilterByPixelFormat(pixFmt)
		buf.Reset()
		jsonEnc.Encode(formats)
		log.Printf("available formats for the pixel format %v:\n%s", *pixFmtFlag, buf.Bytes())
	}

	if *widthFlag != 0 {
		formats = formats.FilterByWidth(*widthFlag)
		buf.Reset()
		jsonEnc.Encode(formats)
		log.Printf("available formats for width %d:\n%s", *widthFlag, buf.Bytes())
	}

	if !math.IsNaN(*fpsFlag) {
		formats = formats.FilterByFPS(*fpsFlag)
		buf.Reset()
		jsonEnc.Encode(formats)
		log.Printf("available formats for FPS %f:\n%s", *fpsFlag, buf.Bytes())
	}

	if len(formats) == 0 {
		panicInUI(w, fmt.Errorf("no appropriate formats available"))
	}

	format := formats.BestResolution()

	log.Printf("requesting format %#+v", format)
	camera, err := plat.OpenCamera(devicePath, format)
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to open the camera: %w", err))
	}
	defer camera.Close()

	log.Printf("starting streaming")
	err = camera.StartStreaming()
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to initiate the streaming on the camera: %w", err))
	}
	defer camera.StopStreaming()
	ctx := context.Background()

	frameReadCtx, cancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFn()

	log.Printf("getting a frame")
	frame, err := camera.GetFrame(frameReadCtx)
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to get a video frame: %w", err))
	}

	log.Printf("releasing the memory buffer of the frame")
	err = camera.ReleaseFrame(frame)
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to release frame %d: %w", frame, err))
	}

	log.Printf("getting the second frame")
	frame, err = camera.GetFrame(frameReadCtx)
	if err != nil {
		panicInUI(w, fmt.Errorf("unable to get a video frame: %w", err))
	}

	img := canvas.NewImageFromImage(frame.Image())
	img.FillMode = canvas.ImageFillOriginal
	img.ScaleMode = canvas.ImageScaleFastest
	w.Canvas().SetContent(img)

	prevFrame := frame
	go func() {
		defer func() { processRecover(w, recover()) }()
		for {
			frame, err := camera.GetFrame(ctx)
			if err != nil {
				panicInUI(w, fmt.Errorf("unable to get a video frame: %w", err))
			}
			img.Image = frame.Image()
			img.Refresh()

			err = camera.ReleaseFrame(prevFrame)
			if err != nil {
				panicInUI(w, fmt.Errorf("unable to release frame %d: %w", frame, err))
			}
			prevFrame = frame
		}
	}()

	a.Run()
}

func panicInUI(
	w fyne.Window,
	err error,
) {
	log.Println(err.Error())
	text := widget.NewLabel(err.Error())
	text.Wrapping = fyne.TextWrapWord
	w.SetContent(text)
	w.ShowAndRun()
	<-context.Background().Done()
}

func processRecover(
	w fyne.Window,
	r any,
) {
	if r == nil {
		return
	}

	debug.PrintStack()
	panicInUI(w, fmt.Errorf("%v\n\n%s", r, debug.Stack()))
}
