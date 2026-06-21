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

// Декларативный интерфейс сквозного прохода атласа (без return и без scrH)
type GlyphDecorator interface {
	RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte)
}

type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte) {}

type GlyphW struct {
	Next GlyphDecorator
}
func (g GlyphW) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte) {
	if charCode == 87 { // ASCII 'W'
		blitRow(glCtx, 0x42, x, y + (0 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (1 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (2 * scale), scale)
		blitRow(glCtx, 0x4A, x, y + (3 * scale), scale)
		blitRow(glCtx, 0x54, x, y + (4 * scale), scale)
		blitRow(glCtx, 0x64, x, y + (5 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (6 * scale), scale)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale)
}

type GlyphO struct {
	Next GlyphDecorator
}
func (g GlyphO) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte) {
	if charCode == 79 { // ASCII 'O'
		blitRow(glCtx, 0x3C, x, y + (0 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (1 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (2 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (3 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (4 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (5 * scale), scale)
		blitRow(glCtx, 0x3C, x, y + (6 * scale), scale)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale)
}

type GlyphK struct {
	Next GlyphDecorator
}
func (g GlyphK) RenderGlyph(glCtx gl.Context, charCode byte, x byte, y byte, scale byte) {
	if charCode == 75 { // ASCII 'K'
		blitRow(glCtx, 0x42, x, y + (0 * scale), scale)
		blitRow(glCtx, 0x44, x, y + (1 * scale), scale)
		blitRow(glCtx, 0x48, x, y + (2 * scale), scale)
		blitRow(glCtx, 0x70, x, y + (3 * scale), scale)
		<markdown></markdown>blitRow(glCtx, 0x48, x, y + (4 * scale), scale)
		blitRow(glCtx, 0x44, x, y + (5 * scale), scale)
		blitRow(glCtx, 0x42, x, y + (6 * scale), scale)
	}
	g.Next.RenderGlyph(glCtx, charCode, x, y, scale)
}

func blitRow(glCtx gl.Context, bits byte, startX byte, y byte, scale byte) {
	if (bits & 0x80) != 0 { drawHWBlock(glCtx, startX + (0 * scale), y, scale) }
	if (bits & 0x40) != 0 { drawHWBlock(glCtx, startX + (1 * scale), y, scale) }
	if (bits & 0x20) != 0 { drawHWBlock(glCtx, startX + (2 * scale), y, scale) }
	if (bits & 0x10) != 0 { drawHWBlock(glCtx, startX + (3 * scale), y, scale) }
	if (bits & 0x08) != 0 { drawHWBlock(glCtx, startX + (4 * scale), y, scale) }
	if (bits & 0x04) != 0 { drawHWBlock(glCtx, startX + (5 * scale), y, scale) }
	if (bits & 0x02) != 0 { drawHWBlock(glCtx, startX + (6 * scale), y, scale) }
	if (bits & 0x01) != 0 { drawHWBlock(glCtx, startX + (7 * scale), y, scale) }
}

// Прямая отрисовка пикселя через верхнюю систему координат (без инверсии)
func drawHWBlock(glCtx gl.Context, x byte, y byte, scale byte) {
	glCtx.Enable(gl.SCISSOR_TEST)
	
	// Отрисовка идет от верхнего левого угла (0,0) вниз
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*4, int32(scale)*4)
	glCtx.ClearColor(0.0, 0.0, 0.0, 1.0) // Черный цвет пикселя буквы
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

type UIValueState struct {
	NotificationChar1 byte
	NotificationChar2 byte
}

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	uiState := &UIValueState{
		NotificationChar1: 87, // 'W'
		NotificationChar2: 87, // 'W'
	}

	atlasChain := GlyphW{
		Next: GlyphO{
			Next: GlyphK{
				Next: EmptyGlyph{},
			},
		},
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
					uiState.NotificationChar1 = 79 // 'O'
					uiState.NotificationChar2 = 75 // 'K'
					uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "NotificationSent")
					a.Send(paint.Event{})
				}
			}

			if _, ok := e.(paint.Event); ok {
				if glCtx == nil {
					a.Send(paint.Event{})
					continue
				}

				w32 := int32(sz.WidthPx)
				h32 := int32(sz.HeightPx)

				// Принудительно очищаем 100% вьюпорта экрана белым цветом
				glCtx.Viewport(0, 0, w32, h32)
				glCtx.Scissor(0, 0, w32, h32)
				glCtx.Enable(gl.SCISSOR_TEST)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				// Координаты (X, Y) и масштаб заданы в рамках 8-битного byte
				var startX byte = 10
				var startY byte = 20
				var textScale byte = 2

				// Сквозной выкат рендеринга букв
				atlasChain.RenderGlyph(glCtx, uiState.NotificationChar1, startX, startY, textScale)
				atlasChain.RenderGlyph(glCtx, uiState.NotificationChar2, startX + 20, startY, textScale)

				// Отключаем локальные scissor-боксы, возвращая управление главному экрану
				glCtx.Disable(gl.SCISSOR_TEST)

				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
