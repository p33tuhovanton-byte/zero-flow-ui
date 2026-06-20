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

// Декларативный интерфейс конвейера атласа. Метод больше ничего не возвращает (void/без return).
type GlyphDecorator interface {
	RenderGlyph(dst draw.Image, charCode byte, x byte, y byte)
}

// Терминальный декоратор — просто тупик конвейера, завершающий сквозной проход
type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(dst draw.Image, charCode byte, x byte, y byte) {
	// Конец цепочки. Нет return, поток просто завершается.
}

// Декоратор символа 'W'
type GlyphW struct {
	Next GlyphDecorator
}
func (g GlyphW) RenderGlyph(dst draw.Image, charCode byte, x byte, y byte) {
	if charCode == 87 { // ASCII 'W'
		blitRow(dst, 0x42, x, y + 0)
		blitRow(dst, 0x42, x, y + 1)
		blitRow(dst, 0x42, x, y + 2)
		blitRow(dst, 0x4A, x, y + 3)
		blitRow(dst, 0x54, x, y + 4)
		blitRow(dst, 0x64, x, y + 5)
		blitRow(dst, 0x42, x, y + 6)
	}
	// Безусловный сквозной проход дальше по конвейеру без прерывания
	g.Next.RenderGlyph(dst, charCode, x, y)
}

// Декоратор символа 'O'
type GlyphO struct {
	Next GlyphDecorator
}
func (g GlyphO) RenderGlyph(dst draw.Image, charCode byte, x byte, y byte) {
	if charCode == 79 { // ASCII 'O'
		blitRow(dst, 0x3C, x, y + 0)
		blitRow(dst, 0x42, x, y + 1)
		blitRow(dst, 0x42, x, y + 2)
		blitRow(dst, 0x42, x, y + 3)
		blitRow(dst, 0x42, x, y + 4)
		blitRow(dst, 0x42, x, y + 5)
		blitRow(dst, 0x3C, x, y + 6)
	}
	g.Next.RenderGlyph(dst, charCode, x, y)
}

// Декоратор символа 'K'
type GlyphK struct {
	Next GlyphDecorator
}
func (g GlyphK) RenderGlyph(dst draw.Image, charCode byte, x byte, y byte) {
	if charCode == 75 { // ASCII 'K'
		blitRow(dst, 0x42, x, y + 0)
		blitRow(dst, 0x44, x, y + 1)
		blitRow(dst, 0x48, x, y + 2)
		blitRow(dst, 0x70, x, y + 3)
		blitRow(dst, 0x48, x, y + 4)
		blitRow(dst, 0x44, x, y + 5)
		blitRow(dst, 0x42, x, y + 6)
	}
	g.Next.RenderGlyph(dst, charCode, x, y)
}

func blitRow(dst draw.Image, bits byte, startX byte, y byte) {
	rect := image.Rect(0, 0, 4, 4)
	blackSrc := &image.Uniform{color.Black}

	if (bits & 0x80) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+0*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x40) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+1*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x20) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+2*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x10) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+3*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x08) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+4*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x04) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+5*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x02) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+6*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x01) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+7*4), int(y*4))), blackSrc, image.Point{}, draw.Src) }
}

type UIValueState struct {
	Char1 byte
	Char2 byte
}

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	uiState := &UIValueState{
		Char1: 87, // 'W'
		Char2: 87, // 'W'
	}

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
				glCtx.Scissor(0, 0, int32(sz.WidthPx), int32(sz.HeightPx))
				glCtx.Enable(gl.SCISSOR_TEST)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				var startX byte = 10
				var startY byte = 30

				// Вызовы выполняются сквозным образом, гарантируя отсутствие зависаний стека
				atlasChain.RenderGlyph(rgba, uiState.Char1, startX, startY)
				atlasChain.RenderGlyph(rgba, uiState.Char2, startX + 10, startY)

				statusBuffer.Upload()
				statusBuffer.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, rgba.Bounds())
				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
