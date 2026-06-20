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

type GlyphDecorator interface {
	RenderGlyph(dst draw.Image, charCode byte, x byte, y byte)
}

type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(dst draw.Image, charCode byte, x byte, y byte) {}

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
	g.Next.RenderGlyph(dst, charCode, x, y)
}

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
		<markdown></markdown>blitRow(dst, 0x3C, x, y + 6)
	}
	g.Next.RenderGlyph(dst, charCode, x, y)
}

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
	rect := image.Rect(0, 0, 8, 8) // Увеличен масштаб для лучшей видимости на экранах телефонов
	blackSrc := &image.Uniform{color.Black}

	if (bits & 0x80) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+0*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x40) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+1*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x20) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+2*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x10) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+3*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x08) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+4*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x04) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+5*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x02) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+6*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x01) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+7*8), int(y*8))), blackSrc, image.Point{}, draw.Src) }
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
			if ev, ok := e.(lifecycle.Event); ok {
				if ev.To == lifecycle.StageAlive {
					// Безопасное по памяти и быстрое извлечение интерфейса (0 аллокаций)
					glCtx, _ = ev.DrawContext.(gl.Context)
					if glCtx != nil {
						images = glutil.NewImages(glCtx)
					}
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
					if statusBuffer != nil { statusBuffer.Release() }
					if images != nil { images.Release() }
					glCtx = nil
				}
			}

			if ev, ok := e.(size.Event); ok {
				sz = ev
				if glCtx != nil && images != nil {
					if statusBuffer != nil { statusBuffer.Release() }
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})
			}

			if ev, ok := e.(touch.Event); ok {
				if ev.Type == touch.TypeBegin {
					uiState.Char1 = 79 // 'O'
					uiState.Char2 = 75 // 'K'
					a.Send(paint.Event{})
				}
			}

			if _, ok := e.(paint.Event); ok {
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					continue
				}

				// Жесткая привязка вьюпорта к физическому экрану смартфона
				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)

				// Полная очистка буферов кадра на уровне GPU
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				// Принудительно заливаем текстуру чистым белым цветом
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				// Сдвигаем координаты ниже (startY=12), чтобы буквы вышли из-под плашки статус-бара
				var startX byte = 4
				var startY byte = 12

				atlasChain.RenderGlyph(rgba, uiState.Char1, startX, startY)
				atlasChain.RenderGlyph(rgba, uiState.Char2, startX + 10, startY)

				statusBuffer.Upload()
				
				// Отрисовка на всю аппаратно-доступную геометрию экрана
				statusBuffer.Draw(
					sz,
					geom.Point{X: 0, Y: 0},
					geom.Point{X: sz.WidthPt, Y: 0},
					geom.Point{X: 0, Y: sz.HeightPt},
					rgba.Bounds(),
				)
				
				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
