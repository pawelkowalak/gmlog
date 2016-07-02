// Package gmlog implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It is anologous to standard library's log
// package, but prints lines of logs on gomobile application.
package gmlog

import (
	"fmt"
	"image"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/size"
	expFont "golang.org/x/mobile/exp/font"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"os"
	"sync"
)

// A Logger represents an active logging object that generates lines
// of output to glutil Images (useful in gomobile applications).
type Logger struct {
	mu     sync.Mutex
	images *glutil.Images
	m      *glutil.Image
	buf    []string
}

// New creates a Logger tied to the current GL images. Limit limits amount
// of lines to be displayed.
func New(images *glutil.Images, limit int) *Logger {
	return &Logger{
		images: images,
		buf:    make([]string, 0, limit),
	}
}

// Output adds formatted message to the buffer.
func (l *Logger) Output(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.buf) > 4 {
		l.buf = l.buf[1:]
	}
	l.buf = append(l.buf, msg)
	fmt.Println(msg)
}

// Printf adds new message format and optional arguments to logger buffer.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

// Print adds new message to logger buffer.
func (l *Logger) Print(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Draw draws all current logs at the top of the screen.
func (l *Logger) Draw(sz size.Event) {
	if sz.PixelsPerPt == 0 || sz.HeightPt == 0 || sz.WidthPt == 0 {
		fmt.Println("size.Event not occurred yet, can't draw logs")
		return
	}
	ttfBytes := expFont.Monospace()
	f, err := truetype.Parse(ttfBytes)
	if err != nil {
		fmt.Println(err)
	}
	if l.m != nil {
		l.m.Release()
	}
	l.m = l.images.NewImage(sz.WidthPx, sz.HeightPx)

	d := &font.Drawer{
		Dst: l.m.RGBA,
		Src: image.White,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    float64(12 * sz.PixelsPerPt),
			Hinting: font.HintingNone,
		}),
	}
	d.Dot = fixed.Point26_6{
		X: 0,
		Y: fixed.I(int(12 * sz.PixelsPerPt)),
	}
	l.mu.Lock()
	for _, m := range l.buf {
		d.DrawString(m)
		d.Dot.X = 0
		d.Dot.Y += fixed.I(int(12 * sz.PixelsPerPt))
	}
	l.mu.Unlock()
	l.m.Upload()
	l.m.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, sz.Bounds())
}

// Release releases image resources.
func (l *Logger) Release() {
	if l.m != nil {
		l.m.Release()
		l.m = nil
		l.images = nil
	}
}
