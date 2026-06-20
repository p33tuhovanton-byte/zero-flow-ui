package main

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

func main() {
	// 1. Используем штатную сборку истории UI из вашей библиотеки
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	// 2. Управляем интерфейсом через структуру TextSignal "из коробки"
	// Используем специальный синтаксис внутри Payload для передачи команд действия (Action)
	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "CMD:SET_THEME_WHITE|MODE:FULLSCREEN|MSG:Zero-Collection UI активен на Android 11!\n",
	}

	// 3. Стандартный декоратор конвейера из вашей библиотеки
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
					// Передаем сигнал в конвейер обработки zeroflowui без изменений кода библиотеки
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				case lifecycle.StageAlive:
					glCtx, _ = x.DrawContext.(gl.Context)
					if glCtx != nil {
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
			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					continue
				}

				// Читаем управляющие параметры интерфейса из существующей структуры textSignal.Payload
				payloadStr := textSignal.Payload
				
				// Флаги конфигурации UI (заполняются на основе Payload библиотеки)
				isWhiteTheme := strings.Contains(payloadStr, "SET_THEME_WHITE")
				isFullScreen := strings.Contains(payloadStr, "FULLSCREEN")

				// Извлекаем чистый текст сообщения для вывода
				displayMsg := payloadStr
				if idx := strings.Index(payloadStr, "MSG:"); idx != -1 {
					displayMsg = payloadStr[idx+4:]
				}

				// Применяем белый цвет очистки экрана на основе структуры сигнала
				if isWhiteTheme {
					glCtx.ClearColor(1.0, 1.0, 1.0, 1.0) // Чистый белый фон
				} else {
					glCtx.ClearColor(0.0, 0.0, 0.0, 1.0) // Дефолтный черный
				}
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				// Заливка текстуры фреймбуфера
				rgba := statusBuffer.RGBA
				bgColor := color.Black
				if isWhiteTheme {
					bgColor = color.White
				}
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

				// Эмуляция ввода полученного текста сообщения из Payload
				_ = displayMsg 

				statusBuffer.Upload()
				
				// Если в сигнале передан режим FULLSCREEN, растягиваем буфер на весь экран
				if isFullScreen {
					statusBuffer.Draw(
						sz,
						geom.Point{X: 0, Y: 0},
						geom.Point{X: sz.WidthPt, Y: 0},
						geom.Point{X: 0, Y: sz.HeightPt},
						rgba.Bounds(),
					)
				}

				a.Publish()
				a.Send(paint.Event{}) // Запрос постоянной перерисовки кадра
			}
		}
	})
}
