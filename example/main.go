package main

import (
  "unsafe"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

// Декларативный интерфейс сквозного прохода атласа символов
type GlyphDecorator interface {
	RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte)
}

type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {}

type GlyphW struct{ Next GlyphDecorator }
func (g GlyphW) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 87 { // 'W'
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x4A, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x54, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x64, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphO struct{ Next GlyphDecorator }
func (g GlyphO) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 79 { // 'O'
		blitRow(glCtx, 0x3C, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x3C, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphK struct{ Next GlyphDecorator }
func (g GlyphK) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 75 { // 'K'
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x44, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x48, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x70, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x48, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x44, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphI struct{ Next GlyphDecorator }
func (g GlyphI) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 73 { // 'I'
		blitRow(glCtx, 0x3C, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x18, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x18, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x18, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x18, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x18, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x3C, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphN struct{ Next GlyphDecorator }
func (g GlyphN) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 110 { // 'n'
		blitRow(glCtx, 0x00, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0xDC, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x62, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphL struct{ Next GlyphDecorator }
func (g GlyphL) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 76 { // 'L'
		blitRow(glCtx, 0x40, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x40, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x40, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x40, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x40, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x40, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x7E, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

type GlyphY struct{ Next GlyphDecorator }
func (g GlyphY) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte, isDisaster byte) {
	if charCode == 121 { // 'y'
		blitRow(glCtx, 0x42, x, y+(0*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(1*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(2*scale), scale, isDisaster)
		blitRow(glCtx, 0x3C, x, y+(3*scale), scale, isDisaster)
		blitRow(glCtx, 0x02, x, y+(4*scale), scale, isDisaster)
		blitRow(glCtx, 0x42, x, y+(5*scale), scale, isDisaster)
		blitRow(glCtx, 0x3C, x, y+(6*scale), scale, isDisaster)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale, isDisaster)
}

func blitRow(glCtx gl.Context, bits byte, startX byte, y byte, scale byte, isDisaster byte) {
	if (bits & 0x80) != 0 { drawHWBlock(glCtx, startX+(0*scale), y, scale, isDisaster) }
	if (bits & 0x40) != 0 { drawHWBlock(glCtx, startX+(1*scale), y, scale, isDisaster) }
	if (bits & 0x20) != 0 { drawHWBlock(glCtx, startX+(2*scale), y, scale, isDisaster) }
	if (bits & 0x10) != 0 { drawHWBlock(glCtx, startX+(3*scale), y, scale, isDisaster) }
	if (bits & 0x08) != 0 { drawHWBlock(glCtx, startX+(4*scale), y, scale, isDisaster) }
	if (bits & 0x04) != 0 { drawHWBlock(glCtx, startX+(5*scale), y, scale, isDisaster) }
	if (bits & 0x02) != 0 { drawHWBlock(glCtx, startX+(6*scale), y, scale, isDisaster) }
	if (bits & 0x01) != 0 { drawHWBlock(glCtx, startX+(7*scale), y, scale, isDisaster) }
}

func drawHWBlock(glCtx gl.Context, x byte, y byte, scale byte, isDisaster byte) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*4, int32(scale)*4)
	if isDisaster == 1 {
		glCtx.ClearColor(1.0, 0.0, 0.0, 1.0) 
	} else {
		glCtx.ClearColor(0.0, 0.0, 0.0, 1.0)
	}
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

// ==========================================
// СТРУКТУРНЫЙ АНАЛИЗАТОР ЦЕПОЧКИ СОБЫТИЙ UI
// ==========================================

type StructuralAtlas struct {
	Chain GlyphDecorator
}

func (sa StructuralAtlas) InterpretUILoopScreen(glCtx gl.Context, flow zeroflowui.UIEventFlow, edgeX byte, currentY byte) {
	if flow == nil {
		return
	}

	descriptor, nextFlow, isEnd := flow()
	if isEnd {
		return
	}

	var disasterFlag byte = 0
	if descriptor.ActionDetails == "Crash" || descriptor.ActionDetails == "Panic" {
		disasterFlag = 1
	}

	var smallScale byte = 1

	if descriptor.EventType == zeroflowui.EventInteraction {
		sa.Chain.RenderGlyph(glCtx, 73, edgeX, currentY, smallScale, disasterFlag)   // 'I'
		sa.Chain.RenderGlyph(glCtx, 110, edgeX+4, currentY, smallScale, disasterFlag) // 'n'
	} else {
		sa.Chain.RenderGlyph(glCtx, 76, edgeX, currentY, smallScale, disasterFlag)   // 'L'
		sa.Chain.RenderGlyph(glCtx, 121, edgeX+4, currentY, smallScale, disasterFlag) // 'y'
	}

	sa.InterpretUILoopScreen(glCtx, nextFlow, edgeX, currentY-10)
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

		*timeline = zeroflowui.LogUIEvent(*timeline, false, zeroflowui.EventInteraction, "NotificationButton", "ClickProcessed")
		pipe.Process(zeroflowui.NewTextFlow(*signal), *timeline)
	}
	b.Next.DispatchTouch(pipe, timeline, tx, ty, signal)
}

func main() {
	var uiTimeline zeroflowui.UIEventFlow = zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

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
var nextStream zeroflowui.StringIterator
var isEnd bool
charStr, nextStream, isEnd = charStream()
if !isEnd && charStr != "" {
     // ИСПРАВЛЕНИЕ: Поэтапное, плоское раскрытие указателя для прохождения чекера gomobile
     strPtr := unsafe.Pointer(&charStr)
     dataPtr := *(*unsafe.Pointer)(strPtr)
     var rawByte1 byte = *(*byte)(dataPtr)
     sysAtlas.Chain.RenderGlyph(glCtx, rawByte1, startX, startY, textScale, 0)
     charStr, nextStream, isEnd = nextStream()
     if !isEnd && charStr != "" {
         strPtr2 := unsafe.Pointer(&charStr)
         dataPtr2 := *(*unsafe.Pointer)(strPtr2)
         var rawByte2 byte = *(*byte)(dataPtr2)
         sysAtlas.Chain.RenderGlyph(glCtx, rawByte2, startX+20, startY, textScale, 0)
      }
}
var topRightX byte = byte(sz.WidthPx>>2) - 15
var topRightY byte = byte(sz.HeightPx>>2) - 15
sysAtlas.InterpretUILoopScreen(glCtx, uiTimeline, topRightX, topRightY)
glCtx.Disable(gl.SCISSOR_TEST)
glCtx.Flush()
a.Publish()
a.Send(paint.Event{})
}
}
})
}