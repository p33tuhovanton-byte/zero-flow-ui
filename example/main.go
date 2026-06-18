package main

import (
	"zeroflowui"
	"golang.org/x/mobile/app"
 "Golang.org/x/mobile/internal"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

func main() {
	// Декларативная сборка истории UI для логов
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "Zero-Collection UI запущен на Android OS!\n",
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
	}

	// Запуск жизненного цикла нативного Android-приложения
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch x := a.Filter(e).(type) {
			case lifecycle.Event:
				if x.To == lifecycle.StageFocused {
					// Выполняем логику конвейера при фокусе приложения
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				}
			case paint.Event:
				// Точка для отрисовки графики (очищаем экран)
				a.Publish()
			}
		}
	})
}
