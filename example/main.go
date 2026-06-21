package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"unsafe"
	"zeroflowui"
)

// Декларативный интерфейс сквозного прохода атласа символов (принимает сырые байты цвета)
type GlyphDecorator interface {
	RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte)
}

type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {}

type GlyphW struct{ Next GlyphDecorator }
func (g GlyphW) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 87 {
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x4A, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x54, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x64, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphO struct{ Next GlyphDecorator }
func (g GlyphO) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 79 {
		blitRow(glCtx, 0x3C, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x3C, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphK struct{ Next GlyphDecorator }
func (g GlyphK) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 75 {
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x44, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x48, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x70, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x48, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x44, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphI struct{ Next GlyphDecorator }
func (g GlyphI) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 73 {
		blitRow(glCtx, 0x3C, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x18, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x18, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x18, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x18, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x18, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x3C, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphN struct{ Next GlyphDecorator }
func (g GlyphN) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 110 {
		blitRow(glCtx, 0x00, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0xDC, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x62, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphL struct{ Next GlyphDecorator }
func (g GlyphL) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 76 {
		blitRow(glCtx, 0x40, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x40, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x40, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x40, x, y+(3*scale), scale, r, g, b)
		blitRow(glCtx, 0x40, x, y+(4*scale), scale, r, g, b)
		blitRow(glCtx, 0x40, x, y+(5*scale), scale, r, g, b)
		blitRow(glCtx, 0x7E, x, y+(6*scale), scale, r, g, b)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

type GlyphY struct{ Next GlyphDecorator }
func (g GlyphY) RenderGlyph(glCtx gl.Context, charCode, x, y, scale, r, g, b byte) {
	if charCode == 121 {
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, r, g, b)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, r, g, b)
		blitRow(glCtx, 0x3C, x, y+(3*scale), scale, isDisaster) // Ошибка компиляции исправлена ниже
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, r, g, b)
}

func blitRow(glCtx gl.Context, bits byte, startX byte, y byte, scale byte, r, g, b byte) {
	if (bits & 0x80) != 0 { drawHWBlock(glCtx, startX+(0*scale), y, scale, r, g, b) }
	if (bits & 0x40) != 0 { drawHWBlock(glCtx, startX+(1*scale), y, scale, r, g, b) }
	if (bits & 0x20) != 0 { drawHWBlock(glCtx, startX+(2*scale), y, scale, r, g, b) }
	if (bits & 0x10) != 0 { drawHWBlock(glCtx, startX+(3*scale), y, scale, r, g, b) }
	if (bits & 0x08) != 0 { drawHWBlock(glCtx, startX+(4*scale), y, scale, r, g, b) }
	if (bits & 0x04) != 0 { drawHWBlock(glCtx, startX+(5*scale), y, scale, r, g, b) }
	if (bits & 0x02) != 0 { drawHWBlock(glCtx, startX+(6*scale), y, scale, r, g, b) }
	if (bits & 0x01) != 0 { drawHWBlock(glCtx, startX+(7*scale), y, scale, r, g, b) }
}

func drawHWBlock(glCtx gl.Context, x byte, y byte, scale byte, r, g, b byte) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*4, int32(scale)*4)
	
	// Переводим значения байта (0 или 1) в вещественные коэффициенты цвета для GPU
	glCtx.ClearColor(float32(r), float32(g), float32(b), 1.0)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

type StructuralAtlas struct {
	Chain GlyphDecorator
}

func (sa StructuralAtlas) InterpretUILoopScreen(glCtx gl.Context, flow zeroflowui.UIEventFlow, edgeX byte, currentY byte, scrH byte) {
	if flow == nil {
		return
	}
	descriptor, nextFlow, isEnd := flow()
	if isEnd {
		return
	}

	// Извлекаем кастомные байты цвета, зашитые во Fluent-дескриптор действия вашей библиотеки
	// Допустим, 0 - черный, 1 - цветной (настраивается на стороне ActionBuilder)
	rByte := descriptor.ColorR
	gByte := descriptor.ColorG
	bByte := descriptor.ColorB

	var smallScale byte = 1
	realY := scrH - currentY

	if descriptor.EventType == zeroflowui.EventInteraction {
		sa.Chain.RenderGlyph(glCtx, 73, edgeX, realY, smallScale, rByte, gByte, bByte)
		sa.Chain.RenderGlyph(glCtx, 110, edgeX+4, realY, smallScale, rByte, gByte, bByte)
	} else {
		sa.Chain.RenderGlyph(glCtx, 76, edgeX, realY, smallScale, rByte, gByte, bByte)
		sa.Chain.RenderGlyph(glCtx, 121, edgeX+4, realY, smallScale, rByte, gByte, bByte)
	}
	sa.InterpretUILoopScreen(glCtx, nextFlow, edgeX, currentY+10, scrH)
}

type UIElementContainer interface {
	DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, signal *zeroflowui.TextSignal)
}

type EndOfUIChain struct{}
func (e EndOfUIChain) DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, signal *zeroflowui.TextSignal) {}

type UINotificationButton struct {
	XMin, XMax, YMin, YMax byte
	Next                   UIElementContainer
}

func (b UINotificationButton) DispatchTouch(pipe *zeroflowui.SystemPipelineDecorator, timeline *zeroflowui.UIEventFlow, tx, ty byte, signal *zeroflowui.TextSignal) {
	if tx >= b.XMin && tx <= b.XMax && ty >= b.YMin && ty <= b.YMax {
		signal.Payload = "OK"

		// ОПЕРИРОВАНИЕ ФУНКЦИОНАЛЬНЫМ ЦВЕТОМ ЛОГОВ ЧЕРЕЗ ВАШ ACTIONBUILDER
		// Задаем цвет лога: Сине-зеленый (R=0, G=1, B=1)
		action := zeroflowui.CreateAction().
			SetComponent("NotificationButton", false).
			SetEvent(zeroflowui.EventInteraction, "ClickProcessed").
			SetColorRGB(0, 1, 1) // Вызываем метод динамической покраски строки лога вашей системы

		action = action.Listen(func(desc zeroflowui.UIStateDescriptor) {})
		*timeline = action.Emit(*timeline)

		pipe.Process(zeroflowui.NewTextFlow(*signal), *timeline)
	}
	b.Next.DispatchTouch(pipe, timeline, tx, ty, signal)
}

func main() {
	var uiTimeline zeroflowui.UIEventFlow = zeroflowui.EndOfUI()
	
	// Стартовый системный лог красим в темно-зеленый цвет (R=0, G=1, B=0)
	uiTimeline = zeroflowui.LogUIEventColored(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered", 0, 1, 0)

	textSignal := &zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "WW",
	}

	glyphChain := GlyphW{
		Next: GlyphO{
			Next: GlyphK{
				Next: GlyphI{
					Next: GlyphN{
						Next: GlyphL{
							Next: GlyphY{
								Next: EmptyGlyph{},
							},
						},
					},
				},
			},
		},
	}

	sysAtlas := StructuralAtlas{Chain: glyphChain}

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
					uiInterfaceChain.DispatchTouch(&pipeline, &uiTimeline, touchX, touchY, textSignal)
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
glCtx.Disable(gl.SCISSOR_TEST)
charStream := zeroflowui.MakeStream(textSignal.Payload)
var charStr string
var nextStream zeroflowui.StringIteratorvar 
isEnd bool
charStr, nextStream, isEnd = charStream()
if !isEnd && charStr != "" {
   strPtr := unsafe.Pointer(&charStr)
   dataPtr := *(*unsafe.Pointer)(strPtr)
   var rawByte1 byte = *(*byte)(dataPtr)
   // Главный статус OK/WW всегда выводим черным цветом (0, 0, 0)
   sysAtlas.Chain.RenderGlyph(glCtx, rawByte1, startX, startY, textScale, 0, 0, 0)
   charStr, nextStream, isEnd = nextStream()
   if !isEnd && charStr != "" {
      strPtr2 := unsafe.Pointer(&charStr)
      dataPtr2 := *(*unsafe.Pointer)(strPtr2)
      var rawByte2 byte = *(*byte)(dataPtr2)
      sysAtlas.Chain.RenderGlyph(glCtx, rawByte2, startX+20, startY, textScale, 0, 0, 0)
   }
}
var topRightX byte = byte(sz.WidthPx>>2) - 15
var screenHeightByte byte = byte(sz.HeightPx >> 2)
// Выводим кастомные цветные логи в правый верхний угол дисплея
sysAtlas.InterpretUILoopScreen(glCtx, uiTimeline, topRightX, 15, screenHeightByte)
glCtx.Disable(gl.SCISSOR_TEST)
glCtx.Flush()
a.Publish()
a.Send(paint.Event{})
}
}
})
}