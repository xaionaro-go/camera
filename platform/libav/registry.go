package libav

import (
	"github.com/xaionaro-go/camera"
)

func init() {
	camera.DefaultRegistry().RegisterPlatform(Platform{})
}
