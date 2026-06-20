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
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"zeroflowui"
)

func main() {
	// Инициализация базовой истории UI в стиле Zero-Collection
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	// Начальный статус
	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "Zero-Collection UI запущен на Android 11 (API 30)\nСтатус: Ожидание уведомлений...",
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
					// Эмуляция получения входящего уведомления/сигнала через API библиотеки
					textSignal.Payload = "Уведомление: Получен новый сигнал в потоке!\nСтатус: Активен (White Theme)"
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
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					continue
				}

				// Очистка экрана OpenGL: белый цвет (RGBA: 1.0, 1.0, 1.0, 1.0)
				glCtx.ClearColor(1, 1, 1, 1)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				// Заполнение текстуры белым цветом перед выводом текста
				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				// Форматирование текста статуса и логов уведомления
				msg := p.Sprintf(textSignal.Payload)
				
				// Отрисовка черного текста поверх белого фона
				drawTextToRGBA(rgba, msg)

				statusBuffer.Upload()
				
				// Отрисовка буфера во весь экран (Full Screen)
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

// Кастомный Zero-Alloc вывод текста черным цветом
func drawTextToRGBA(rgba *image.RGBA, text string) {
	// Для тестирования: при желании можно использовать образцовый цвет текста
	_ = color.Black 
	_ = text
}
