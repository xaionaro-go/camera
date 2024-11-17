package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"math"
	"net/http"
	_ "net/http/pprof"

	"log"
	"os"
	"time"

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

	var plat camera.Platform
	var devicePath camera.DevicePath
	if *platformFlag != "" {
		plat = allplatforms.Get(*platformFlag)
		if plat == nil {
			panic(fmt.Errorf("platform '%s' is unknown", *platformFlag))
		}
		availableCameras, err := plat.ListCameras()
		if err != nil {
			panic(fmt.Errorf("unable to list cameras: %w", err))
		}
		for _, c := range availableCameras {
			if c == *deviceFlag {
				devicePath = c
				break
			}
		}
		if devicePath == "" {
			panic(fmt.Errorf("camera with path '%s' is not found (available: %#+v)", *deviceFlag, availableCameras))
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
			panic(fmt.Errorf("camera with path '%s' is not found (available: %#+v)", *deviceFlag, availableCameras))
		}
	}

	formats, err := plat.ListFormats(devicePath)
	if err != nil {
		panic(fmt.Errorf("unable to list the formats: %w", err))
	}
	if len(formats) == 0 {
		panic(fmt.Errorf("the list of available formats is empty"))
	}

	var buf bytes.Buffer
	jsonEnc := json.NewEncoder(&buf)
	jsonEnc.SetIndent("", " ")
	jsonEnc.Encode(formats)
	log.Printf("available formats:\n%s", buf.Bytes())

	if *pixFmtFlag != "" {
		pixFmt := camera.PixelFormatByName(*pixFmtFlag)
		if pixFmt == camera.PixelFormatUndefined {
			panic(fmt.Errorf("unknown pixel format name '%s'", *pixFmtFlag))
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
		panic("no appropriate formats available")
	}

	format := formats.BestResolution()

	log.Printf("requesting format %#+v", format)
	camera, err := plat.OpenCamera(devicePath, format)
	if err != nil {
		panic(fmt.Errorf("unable to open the camera: %w", err))
	}
	defer camera.Close()

	log.Printf("starting streaming")
	err = camera.StartStreaming()
	if err != nil {
		panic(fmt.Errorf("unable to initiate the streaming on the camera: %w", err))
	}
	defer camera.StopStreaming()
	ctx := context.Background()

	frameReadCtx, cancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFn()

	log.Printf("getting a frame")
	frame, err := camera.GetFrame(frameReadCtx)
	if err != nil {
		panic(fmt.Errorf("unable to get a video frame: %w", err))
	}

	log.Printf("releasing the memory buffer of the frame")
	err = camera.ReleaseFrame(frame)
	if err != nil {
		panic(fmt.Errorf("unable to release frame %d: %w", frame, err))
	}

	log.Printf("getting the second frame")
	frame, err = camera.GetFrame(frameReadCtx)
	if err != nil {
		panic(fmt.Errorf("unable to get a video frame: %w", err))
	}

	log.Printf("encoding the picture into PNG")
	err = png.Encode(os.Stdout, frame.Image())
	if err != nil {
		panic(fmt.Errorf("unable to encode the frame into the PNG file: %w", err))
	}

	log.Printf("releasing the memory buffer of the frame")
	err = camera.ReleaseFrame(frame)
	if err != nil {
		panic(fmt.Errorf("unable to release frame %d: %w", frame, err))
	}
}
