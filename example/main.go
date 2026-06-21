package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "",
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
	}

	var glCtx gl.Context
	var sz size.Event

	app.Main(func(a app.App) {
		for e := range a.Events() {
			if ev, ok := e.(lifecycle.Event); ok {
				if ev.To == lifecycle.StageAlive {
					glCtx, _ = ev.DrawContext.(gl.Context)
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageVisible {
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageFocused {
					glCtx, _ = ev.DrawContext.(gl.Context)
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageDead {
					glCtx = nil
				}
			}

			if ev, ok := e.(size.Event); ok {
				sz = ev
				a.Send(paint.Event{})
			}

			if ev, ok := e.(touch.Event); ok {
				if ev.Type == touch.TypeBegin {
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
					a.Send(paint.Event{})
				}
			}

			if _, ok := e.(paint.Event); ok {
				if glCtx == nil {
					a.Send(paint.Event{})
					continue
				}

				// АППАРАТНОЕ УНИЧТОЖЕНИЕ ТЕМЫ ANDROID НА УРОВНЕ ВИДЕОКАРТЫ
				// Передаем системные параметры размеров экрана напрямую в драйвер
				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Scissor(0, 0, int32(sz.WidthPx), int32(sz.HeightPx))
				
				// Полностью выключаем тест отсечения окна, чтобы стереть черную верхнюю плашку
				glCtx.Disable(gl.SCISSOR_TEST)

				// Принудительно заливаем абсолютно весь экран чистым белым цветом (1.0, 1.0, 1.0, 1.0)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				// Синхронизируем кадры GPU
				glCtx.Flush()
				
				// Команда Publish мгновенно выводит белый буфер из Go прямо на дисплей смартфона
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
