package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom" // Добавлен импорт для работы с координатами экрана
	"golang.org/x/mobile/gl"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"zeroflowui"
)

func main() {
	// Декларативная сборка истории UI для логов
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "Zero-Collection UI запущен на Android OS!\nСтатус API: Активен\nПотоки логов: OK",
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
	}

	var glCtx gl.Context
	var images *glutil.Images
	var statusBuffer *glutil.Image
	var sz size.Event

	p := message.NewPrinter(language.English)

	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch x := a.Filter(e).(type) {
			case lifecycle.Event:
				switch x.To {
				case lifecycle.StageFocused:
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				case lifecycle.StageAlive:
					glCtx, _ = x.DrawContext.(gl.Context)
					if glCtx != nil {
						images = glutil.NewImages(glCtx)
					}
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
					// Текстура создается под физические пиксели
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					continue
				}

				glCtx.ClearColor(0, 0, 0, 1)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

				msg := p.Sprintf(textSignal.Payload)
				drawTextToRGBA(rgba, msg)

				statusBuffer.Upload()
				
				// Корректная отрисовка с использованием geom.Point и geom.Pt
				statusBuffer.Draw(
					sz,
					geom.Point{X: 0, Y: 0},
					geom.Point{X: sz.WidthPt, Y: 0},
					geom.Point{X: 0, Y: sz.HeightPt},
					rgba.Bounds(),
				)

				a.Publish()
			}
		}
	})
}

func drawTextToRGBA(rgba *image.RGBA, text string) {
	_ = text 
}
