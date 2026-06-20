package main

import (
	"image"
	"image/color"
	"image/draw"
	"unsafe" // Прямое управление указателями памяти процессора

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

// Декларативный интерфейс блиттинга, оперирующий исключительно типом byte
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
		blitRow(dst, 0x3C, x, y + 6)
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
	rect := image.Rect(0, 0, 6, 6)
	blackSrc := &image.Uniform{color.Black}

	if (bits & 0x80) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+0*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x40) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+1*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x20) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+2*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x10) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+3*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x08) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+4*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x04) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+5*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x02) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+6*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
	if (bits & 0x01) != 0 { draw.Draw(dst, rect.Bounds().Add(image.Pt(int(startX+7*6), int(y*6))), blackSrc, image.Point{}, draw.Src) }
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
			// ЛИНЕЙНЫЙ ФИЛЬТР СОБЫТИЙ БЕЗ SWITCH И СВЯЗАННЫХ КОЛЛЕКЦИЙ

			// 1. Обработка Жизненного Цикла через unsafe (без проверки ctx, ok)
			if ev, ok := e.(lifecycle.Event); ok {
				if ev.To == lifecycle.StageAlive {
					// Читаем интерфейс DrawContext прямо из памяти по указателю (0 allocs)
					glCtx = *(*gl.Context)(unsafe.Pointer(&ev.DrawContext))
					images = glutil.NewImages(glCtx)
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageVisible {
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageFocused {
					glCtx = *(*gl.Context)(unsafe.Pointer(&ev.DrawContext))
					a.Send(paint.Event{})
				}
				if ev.To == lifecycle.StageDead {
					if statusBuffer != nil { statusBuffer.Release() }
					if images != nil { images.Release() }
					glCtx = nil
				}
			}

			// 2. Срез события размеров экрана (без создания промежуточных переменных)
			if ev, ok := e.(size.Event); ok {
				sz = ev
				if glCtx != nil && images != nil {
					if statusBuffer != nil { statusBuffer.Release() }
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})
			}

			// 3. Тапы по экрану
			if ev, ok := e.(touch.Event); ok {
				if ev.Type == touch.TypeBegin {
					uiState.Char1 = 79 // 'O'
					uiState.Char2 = 75 // 'K'
					a.Send(paint.Event{})
				}
			}

			// 4. Графический цикл Paint
			if _, ok := e.(paint.Event); ok {
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					continue
				}

				// Извлекаем размеры фреймбуфера напрямую из внутренней памяти Bounds
				// Это исключает ручное приведение int32(sz.WidthPx), убирая ошибку "cannot use h32"
				rectMax := statusBuffer.RGBA.Bounds().Max

				glCtx.Viewport(0, 0, rectMax.X, rectMax.Y)
				glCtx.Scissor(0, 0, *(*int32)(unsafe.Pointer(&rectMax.X)), *(*int32)(unsafe.Pointer(&rectMax.Y)))
				glCtx.Disable(gl.SCISSOR_TEST)

				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				var startX byte = 20
				var startY byte = 60

				// Сквозной выкат цепочки декораторов по byte-координатам без return
				atlasChain.RenderGlyph(rgba, uiState.Char1, startX, startY)
				atlasChain.RenderGlyph(rgba, uiState.Char2, startX + 12, startY)

				statusBuffer.Upload()
				statusBuffer.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, rgba.Bounds())
				
				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
