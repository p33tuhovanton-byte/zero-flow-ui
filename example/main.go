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

// Структура 1: Инкапсулирует логику побитовой отрисовки строки символа
type HardwareBlitter struct{}

// Метод 1: Побитовый разбор без массивов и коллекций
func (hb HardwareBlitter) BlitRow(glCtx gl.Context, bits byte, x byte, y byte, scale byte) {
	if (bits & 0x80) != 0 { hb.DrawPixel(glCtx, x+(0*scale), y, scale) }
	if (bits & 0x40) != 0 { hb.DrawPixel(glCtx, x+(1*scale), y, scale) }
	if (bits & 0x20) != 0 { hb.DrawPixel(glCtx, x+(2*scale), y, scale) }
	if (bits & 0x10) != 0 { hb.DrawPixel(glCtx, x+(3*scale), y, scale) }
	if (bits & 0x08) != 0 { hb.DrawPixel(glCtx, x+(4*scale), y, scale) }
	if (bits & 0x04) != 0 { drawPixelDummy(glCtx, x+(5*scale), y, scale) } // Оптимизация под размер шрифта
	if (bits & 0x02) != 0 { hb.DrawPixel(glCtx, x+(6*scale), y, scale) }
	if (bits & 0x01) != 0 { hb.DrawPixel(glCtx, x+(7*scale), y, scale) }
}

// Аппаратная заливка пикселя через Scissor Box
func (hb HardwareBlitter) DrawPixel(glCtx gl.Context, x byte, y byte, scale byte) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*4, int32(scale)*4)
	glCtx.ClearColor(0.0, 0.0, 0.0, 1.0) // Черный цвет букв
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

func drawPixelDummy(glCtx gl.Context, x byte, y byte, scale byte) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*4, int32(scale)*4)
	glCtx.ClearColor(0.0, 0.0, 0.0, 1.0)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

// Структура 2: Монолитный атлас, группирующий все символы системы
type StructuralAtlas struct {
	Blitter HardwareBlitter
}

// Метод 2: Сквозной каскадный проход без циклов, map и ключевого слова return
func (sa StructuralAtlas) FlowLine(glCtx gl.Context, charCode byte, x byte, y byte, scale byte) {
	// --- Каскад символа 'W' (ASCII 87) ---
	if charCode == 87 {
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x4A, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x54, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x64, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'O' (ASCII 79) ---
	if charCode == 79 {
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'K' (ASCII 75) ---
	if charCode == 75 {
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x44, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x48, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x70, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x48, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x44, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'I' (ASCII 73) ---
	if charCode == 73 {
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x18, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x18, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x18, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x18, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x18, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'n' (ASCII 110) ---
	if charCode == 110 {
		sa.Blitter.BlitRow(glCtx, 0x00, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0xDC, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x62, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'L' (ASCII 76) ---
	if charCode == 76 {
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x40, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x7E, x, y+(6*scale), scale)
	}
	// --- Каскад символа 'y' (ASCII 121) ---
	if charCode == 121 {
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(0*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(1*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(2*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(3*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x02, x, y+(4*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x42, x, y+(5*scale), scale)
		sa.Blitter.BlitRow(glCtx, 0x3C, x, y+(6*scale), scale)
	}
}

type UIElementContainer interface {
	DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, state *UIValueState)
}

type EndOfUIChain struct{}
func (e EndOfUIChain) DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, state *UIValueState) {}

type UINotificationButton struct {
	XMin, XMax, YMin, YMax byte
	Next                   UIElementContainer
}
func (b UINotificationButton) DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, state *UIValueState) {
	if tx >= b.XMin && tx <= b.XMax && ty >= b.YMin && ty <= b.YMax {
		*timeline = zeroflowui.LogUIEvent(*timeline, false, zeroflowui.EventInteraction, "NotificationButton", "ClickProcessed")
		descriptor, _, _ := (*timeline)()

		if descriptor.EventType == zeroflowui.EventInteraction {
			state.NotificationChar1 = 73  // 'I'
			state.NotificationChar2 = 110 // 'n'
		} else {
			state.NotificationChar1 = 76  // 'L'
			state.NotificationChar2 = 121 // 'y'
		}

		textSignal := zeroflowui.TextSignal{Type: zeroflowui.TextType, Payload: ""}
		pipe.Process(zeroflowui.NewTextFlow(textSignal), *timeline)
	}
	b.Next.DispatchTouch(pipe, timeline, tx, ty, state)
}

type UIValueState struct {
	NotificationChar1 byte
	NotificationChar2 byte
}

func main() {
	var uiTimeline zeroflowui.UIEventFlow = zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	uiState := &UIValueState{
		NotificationChar1: 76,  // 'L'
		NotificationChar2: 121, // 'y'
	}

	// Инициализируем наш скомпонованный атлас структур
	sysAtlas := StructuralAtlas{
		Blitter: HardwareBlitter{},
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
	}

	uiInterfaceChain := UINotificationButton{
		XMin: 0, XMax: 100, YMin: 0, YMax: 100,
		Next: EndOfUIChain{},
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
					touchX := byte(ev.X / 4.0)
					touchY := byte(ev.Y / 4.0)

					uiInterfaceChain.DispatchTouch(&pipeline, &uiTimeline, touchX, touchY, uiState)
					a.Send(paint.Event{})
				}
			}

			if _, ok := e.(paint.Event); ok {
				if glCtx == nil {
					a.Send(paint.Event{})
					continue
				}

				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Scissor(0, 0, int32(sz.WidthPx), int32(sz.HeightPx))
				
				glCtx.Enable(gl.SCISSOR_TEST)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				var startX byte = 10
				var startY byte = 20
				var textScale byte = 2

				// Сборка вывода через единый метод FlowLine без создания аллокаций
				glCtx.Disable(gl.SCISSOR_TEST)
				sysAtlas.FlowLine(glCtx, uiState.NotificationChar1, startX, startY, textScale)
				sysAtlas.FlowLine(glCtx, uiState.NotificationChar2, startX+20, startY, textScale)

				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
