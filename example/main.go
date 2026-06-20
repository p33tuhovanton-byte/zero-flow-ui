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
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	// Управляем интерфейсом через штатный Payload
	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "CMD:SET_THEME_WHITE|MODE:FULLSCREEN|MSG:Zero-Collection UI активен!",
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
					// Принудительно извлекаем контекст драйвера Android
					if ctx, ok := x.DrawContext.(gl.Context); ok {
						glCtx = ctx
						images = glutil.NewImages(glCtx)
					}
					a.Send(paint.Event{}) // Инициируем первый цикл отрисовки
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
				// Пересоздаем буфер строго под новые размеры экрана
				if glCtx != nil && images != nil {
					if statusBuffer != nil {
						statusBuffer.Release()
					}
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})
			case paint.Event:
				// Если Android еще не подготовил контекст, запрашиваем paint повторно
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					continue
				}

				payloadStr := textSignal.Payload
				isWhiteTheme := strings.Contains(payloadStr, "SET_THEME_WHITE")
				isFullScreen := strings.Contains(payloadStr, "FULLSCREEN")

				// 1. Жёсткая очистка основного экрана Android через glCtx напрямую
				if isWhiteTheme {
					glCtx.ClearColor(1.0, 1.0, 1.0, 1.0) // Белый
				} else {
					glCtx.ClearColor(0.0, 0.0, 0.0, 1.0) // Черный
				}
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				// 2. Синхронизируем пиксельную карту фреймбуфера
				rgba := statusBuffer.RGBA
				bgColor := color.Black
				if isWhiteTheme {
					bgColor = color.White
				}
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

				// Загружаем текстуру в память GPU
				statusBuffer.Upload()
				
				// 3. Вывод на полный экран с принудительным игнорированием системных отступов
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
				
				// Зацикливаем paint.Event, чтобы Android не успевал подменить холст на дефолтный серый
				a.Send(paint.Event{})
			}
		}
	})
}
