package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
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
	var sz size.Event

	// Запуск жизненного цикла нативного Android-приложения
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch x := a.Filter(e).(type) {
			case lifecycle.Event:
				switch x.To {
				case lifecycle.StageFocused:
					// Выполняем логику конвейера при фокусе приложения
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				case lifecycle.StageAlive:
					glCtx, _ = x.DrawContext.(gl.Context)
					if glCtx != nil {
						images = glutil.NewImages(glCtx)
					}
				case lifecycle.StageDead:
					if images != nil {
						images.Release()
					}
					glCtx = nil
				}
			case size.Event:
				sz = x
			case paint.Event:
				if glCtx == nil || images == nil {
					continue
				}

				// Очистка экрана (черный фон для консольного UI статусов)
				glCtx.ClearColor(0, 0, 0, 1)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				// Отрисовка UI статусов методами вашей библиотеки zeroflowui
				// Передаем контекст GL и текущий размер экрана для вывода текста сигналов
				zeroflowui.RenderNativeTextStatus(glCtx, images, sz, textSignal, uiTimeline)

				a.Publish()
			}
		}
	})
}
