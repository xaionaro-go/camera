package camera

import (
	"fmt"
	"reflect"
	"sync"
)

type Registry struct {
	locker            sync.Mutex
	platforms         []Platform
	alreadyRegistered map[reflect.Type]struct{}
}

var defaultRegistry = NewRegistry()

func DefaultRegistry() *Registry {
	return defaultRegistry
}

func NewRegistry() *Registry {
	return &Registry{
		locker:            sync.Mutex{},
		platforms:         []Platform{},
		alreadyRegistered: map[reflect.Type]struct{}{},
	}
}

func (r *Registry) RegisterPlatform(plat Platform) {
	r.locker.Lock()
	defer r.locker.Unlock()

	t := reflect.TypeOf(plat)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if _, ok := r.alreadyRegistered[t]; ok {
		panic(fmt.Errorf("type '%T' is already registered", t))
	}

	r.platforms = append(r.platforms, plat)
}

type DevicePathAndPlatform struct {
	DevicePath DevicePath
	Platform   Platform
}

func (d DevicePathAndPlatform) ListFormats() (Formats, error) {
	return d.Platform.ListFormats(d.DevicePath)
}

func (d DevicePathAndPlatform) OpenCamera(format Format) (Camera, error) {
	return d.Platform.OpenCamera(d.DevicePath, format)
}

func (r *Registry) ListCameras() ([]DevicePathAndPlatform, error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	var result []DevicePathAndPlatform
	for _, plat := range r.platforms {
		cameras, _ := plat.ListCameras()
		for _, devicePath := range cameras {
			result = append(result, DevicePathAndPlatform{
				DevicePath: devicePath,
				Platform:   plat,
			})
		}
	}
	return result, nil
}
