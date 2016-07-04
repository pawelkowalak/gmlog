Simple logging package for gomobile applications. It supports Print, Printf, Fatal and Fatalf methods.

`go get github.com/viru/gmlog`

Typical usage (simplified code to show the idea):

```go
package main

import (
	"github.com/viru/gmlog"
	// other imports left out for readibility
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					onStart(glctx, sz) // Initialise logger here.
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop() // Release logger here.
				}
			case size.Event:
				sz = e
			case paint.Event:
				onPaint(glctx, sz) // Make sure logger Draw func is called here.
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}

var (
	images    *glutil.Images
	eng       sprite.Engine
	log       *gmlog.Logger // Global log handler.
)

func onStart(glctx gl.Context, sz size.Event) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	log = gmlog.New(images, 5)) // Pass your images to logger along with lines limit.
	log.Print("Initialised app") // Now you can print to your logs.
}

func onStop() {
	eng.Release()
	log.Release() // Release logger before releasing images.
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(0, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	eng.Render(scene, now, sz)
	log.Draw(sz) // Call Draw to display the logs.
}
```

Logs are added on top of the screen (not configurable at the moment):

[[https://github.com/viru/gmlog.wiki/blob/master/screenshot.png|alt=screenshot]]