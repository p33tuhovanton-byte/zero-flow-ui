package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

// Интерфейс декоратора символа. Никаких массивов в аргументах.
type GlyphDecorator interface {
	RenderGlyph(rgba *image.RGBA, charCode byte, x, y, scale int) bool
}

// Терминальный декоратор (конец цепочки атласа)
type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(rgba *image.RGBA, charCode byte, x, y, scale int) bool {
	return false // Символ не найден в атласе
}

// Структурный декоратор для символа 'W'
type GlyphW struct {
	Next GlyphDecorator
}
func (g GlyphW) RenderGlyph(rgba *image.RGBA, charCode byte, x, y, scale int) bool {
	if charCode == 87 { // ASCII 'W'
		// Построчный накат битовой маски 8x8 через примитивные вызовы (0 аллокаций)
		drawRow(rgba, 0x42, x, y+0*scale, scale)
		drawRow(rgba, 0x42, x, y+1*scale, scale)
		drawRow(rgba, 0x42, x, y+2*scale, scale)
		drawRow(rgba, 0x4A, x, y+3*scale, scale)
		drawRow(rgba, 0x54, x, y+4*scale, scale)
		drawRow(rgba, 0x64, x, y+5*scale, scale)
		drawRow(rgba, 0x42, x, y+6*scale, scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, x, y, scale)
}

// Структурный декоратор для символа 'O'
type GlyphO struct {
	Next GlyphDecorator
}
func (g GlyphO) RenderGlyph(rgba *image.RGBA, charCode byte, x, y, scale int) bool {
	if charCode == 79 { // ASCII 'O'
		drawRow(rgba, 0x3C, x, y+0*scale, scale)
		drawRow(rgba, 0x42, x, y+1*scale, scale)
		drawRow(rgba, 0x42, x, y+2*scale, scale)
		drawRow(rgba, 0x42, x, y+3*scale, scale)
		drawRow(rgba, 0x42, x, y+4*scale, scale)
		drawRow(rgba, 0x42, x, y+5*scale, scale)
		drawRow(rgba, 0x3C, x, y+6*scale, scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, x, y, scale)
}

// Структурный декоратор для символа 'K'
type GlyphK struct {
	Next GlyphDecorator
}
func (g GlyphK) RenderGlyph(rgba *image.RGBA, charCode byte, x, y, scale int) bool {
	if charCode == 75 { // ASCII 'K'
		drawRow(rgba, 0x42, x, y+0*scale, scale)
		drawRow(rgba, 0x44, x, y+1*scale, scale)
		drawRow(rgba, 0x48, x, y+2*scale, scale)
		drawRow(rgba, 0x70, x, y+3*scale, scale)
		drawRow(rgba, 0x48, x, y+4*scale, scale)
		drawRow(rgba, 0x44, x, y+5*scale, scale)
		drawRow(rgba, 0x42, x, y+6*scale, scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, x, y, scale)
}

// Вспомогательная функция отрисовки битов строки (0 аллокаций)
func drawRow(rgba *image.RGBA, bits byte, startX, y, scale int) {
	if (bits & 0x80) != 0 { drawPixelBlock(rgba, startX+0*scale, y, scale) }
	if (bits & 0x40) != 0 { drawPixelBlock(rgba, startX+1*scale, y, scale) }
	if (bits & 0x20) != 0 { drawPixelBlock(rgba, startX+2*scale, y, scale) }
	if (bits & 0x10) != 0 { drawPixelBlock(rgba, startX+3*scale, y, scale) }
	if (bits & 0x08) != 0 { drawPixelBlock(rgba, startX+4*scale, y, scale) }
	if (bits & 0x04) != 0 { drawPixelBlock(rgba, startX+5*scale, y, scale) }
	if (bits & 0x02) != 0 { drawPixelBlock(rgba, startX+6*scale, y, scale) }
	if (bits & 0x01) != 0 { drawPixelBlock(rgba, startX+7*scale, y, scale) }
}

func drawPixelBlock(rgba *image.RGBA, startX, startY, scale int) {
	for sy := 0; sy < scale; sy++ {
		for sx := 0; sx < scale; sx++ {
			rgba.SetRGBA(startX+sx, startY+sy, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		}
	}
}

// Управляющая структура состояния UI без использования массивов
type UIValueState struct {
	Char1 byte
	Char2 byte
}

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	// Статическое состояние UI-символов
	uiState := &UIValueState{
		Char1: 87, // 'W'
		Char2: 87, // 'W'
	}

	// Собираем структурный атлас через декораторы
	atlasChain := GlyphW{
		Next: GlyphO{
			Next: GlyphK{
				Next: EmptyGlyph{},
			},
		},
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
				case lifecycle.StageAlive:
					if ctx, ok := x.DrawContext.(gl.Context); ok {
						glCtx = ctx
						images = glutil.NewImages(glCtx)
					}
					a.Send(paint.Event{})
				case lifecycle.StageDead:
					if statusBuffer != nil { statusBuffer.Release() }
					if images != nil { images.Release() }
					glCtx = nil
				}
			case size.Event:
				sz = x
				if glCtx != nil && images != nil {
					if statusBuffer != nil { statusBuffer.Release() }
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})
			case touch.Event:
				if x.Type == touch.TypeBegin {
					// Атомарное изменение состояния ячеек структуры полей
					uiState.Char1 = 79 // 'O'
					uiState.Char2 = 75 // 'K'
					a.Send(paint.Event{})
				}
			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					continue
				}

				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Scissor(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Enable(gl.SCISSOR_TEST)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				// Программный вызов цепочки декораторов по отдельным полям
				atlasChain.RenderGlyph(rgba, uiState.Char1, 40, 120, 4)
				atlasChain.RenderGlyph(rgba, uiState.Char2, 80, 120, 4)

				statusBuffer.Upload()
				statusBuffer.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, rgba.Bounds())
				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
