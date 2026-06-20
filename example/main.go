package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

var zeroFont8x8 = map[rune]byte{
	'Z': {0x7E, 0x0C, 0x18, 0x30, 0x60, 0x42, 0x7E, 0x00},
	'e': {0x3C, 0x42, 0x7E, 0x40, 0x40, 0x42, 0x3C, 0x00},
	'r': {0x5C, 0x62, 0x40, 0x40, 0x40, 0x40, 0xE0, 0x00},
	'o': {0x3C, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'F': {0x7E, 0x40, 0x40, 0x78, 0x40, 0x40, 0x40, 0x00},
	'l': {0x60, 0x60, 0x60, 0x60, 0x60, 0x60, 0x7E, 0x00},
	'w': {0x42, 0x42, 0x42, 0x4A, 0x54, 0x64, 0x42, 0x00},
	'U': {0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'I': {0x3C, 0x18, 0x18, 0x18, 0x18, 0x18, 0x3C, 0x00},
	':': {0x00, 0x18, 0x18, 0x00, 0x18, 0x18, 0x00, 0x00},
	'O': {0x3C, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'K': {0x42, 0x44, 0x48, 0x70, 0x48, 0x44, 0x42, 0x00},
	' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
}

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "ZeroFlowUI: Wait...",
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
	}

	var glCtx gl.Context
	var images *glutil.Images
	var statusBuffer *glutil.Image
	var sz size.Event

	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch x := a.Filter(e).(type) {
			case lifecycle.Event:
				switch x.To {
				case lifecycle.StageFocused:
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				case lifecycle.StageAlive:
					if ctx, ok := x.DrawContext.(gl.Context); ok {
						glCtx = ctx
						images = glutil.NewImages(glCtx)
					}
					a.Send(paint.Event{})
				case lifecycle.StageDead:
					if statusBuffer != nil {
						statusBuffer.Release()
						statusBuffer = nil
					}
					if images != nil {
						images.Release()
						images = nil
					}
					glCtx = nil
				}

			case size.Event:
				sz = x
				if glCtx != nil && images != nil {
					if statusBuffer != nil {
						statusBuffer.Release()
					}
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})

			case touch.Event:
				if x.Type == touch.TypeBegin {
					textSignal.Payload = "ZeroFlowUI: OK"
					uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Touched")
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
					a.Send(paint.Event{})
				}

			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					return
				}

				// НАДЁЖНАЯ АППАРАТНАЯ ИНИЦИАЛИЗАЦИЯ И ОЧИСТКА ЭКРАНА
				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Scissor(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Enable(gl.SCISSOR_TEST)

				// Задаём белый цвет фона
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				// Выводим текст логов черным цветом
				drawZeroAllocText(rgba, textSignal.Payload, 40, 120, 4)

				statusBuffer.Upload()
				
				statusBuffer.Draw(
					sz,
					geom.Point{X: 0, Y: 0},
					geom.Point{X: sz.WidthPt, Y: 0},
					geom.Point{X: 0, Y: sz.HeightPt},
					rgba.Bounds(),
				)

				// Принудительно заставляем GPU сбросить буфер и обновить картинку на экране
				glCtx.Flush()

				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}

func drawZeroAllocText(rgba *image.RGBA, text string, startX, startY, scale int) {
	currentX := startX
	for _, char := range text {
		bitmap, exists := zeroFont8x8[char]
		if !exists {
			bitmap = zeroFont8x8[' ']
		}
		for row := 0; row < 8; row++ {
			bits := bitmap[row]
			for col := 0; col < 8; col++ {
				if (bits & (1 << (7 - uint(col)))) != 0 {
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							px := currentX + col*scale + sx
							py := startY + row*scale + sy
							rgba.SetRGBA(px, py, color.RGBA{R: 0, G: 0, B: 0, A: 255})
						}
					}
				}
			}
		}
		currentX += 10 * scale
	}
}
