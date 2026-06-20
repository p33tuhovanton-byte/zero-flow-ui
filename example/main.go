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

// GlyphDecorator оперирует исключительно типом byte
type GlyphDecorator interface {
	RenderGlyph(rgba *image.RGBA, charCode byte, xHigh, xLow, yHigh, yLow, scale byte) bool
}

type EmptyGlyph struct{}
func (e EmptyGlyph) RenderGlyph(rgba *image.RGBA, charCode byte, xHigh, xLow, yHigh, yLow, scale byte) bool {
	return false 
}

type GlyphW struct {
	Next GlyphDecorator
}
func (g GlyphW) RenderGlyph(rgba *image.RGBA, charCode byte, xHigh, xLow, yHigh, yLow, scale byte) bool {
	if charCode == 87 { // ASCII 'W'
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (0 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (1 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (2 * scale), scale)
		drawRow(rgba, 0x4A, xHigh, xLow, yHigh, yLow + (3 * scale), scale)
		drawRow(rgba, 0x54, xHigh, xLow, yHigh, yLow + (4 * scale), scale)
		drawRow(rgba, 0x64, xHigh, xLow, yHigh, yLow + (5 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (6 * scale), scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, xHigh, xLow, yHigh, yLow, scale)
}

type GlyphO struct {
	Next GlyphDecorator
}
func (g GlyphO) RenderGlyph(rgba *image.RGBA, charCode byte, xHigh, xLow, yHigh, yLow, scale byte) bool {
	if charCode == 79 { // ASCII 'O'
		drawRow(rgba, 0x3C, xHigh, xLow, yHigh, yLow + (0 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (1 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (2 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (3 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (4 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (5 * scale), scale)
		drawRow(rgba, 0x3C, xHigh, xLow, yHigh, yLow + (6 * scale), scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, xHigh, xLow, yHigh, yLow, scale)
}

type GlyphK struct {
	Next GlyphDecorator
}
func (g GlyphK) RenderGlyph(rgba *image.RGBA, charCode byte, xHigh, xLow, yHigh, yLow, scale byte) bool {
	if charCode == 75 { // ASCII 'K'
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (0 * scale), scale)
		drawRow(rgba, 0x44, xHigh, xLow, yHigh, yLow + (1 * scale), scale)
		drawRow(rgba, 0x48, xHigh, xLow, yHigh, yLow + (2 * scale), scale)
		drawRow(rgba, 0x70, xHigh, xLow, yHigh, yLow + (3 * scale), scale)
		drawRow(rgba, 0x48, xHigh, xLow, yHigh, yLow + (4 * scale), scale)
		drawRow(rgba, 0x44, xHigh, xLow, yHigh, yLow + (5 * scale), scale)
		drawRow(rgba, 0x42, xHigh, xLow, yHigh, yLow + (6 * scale), scale)
		return true
	}
	return g.Next.RenderGlyph(rgba, charCode, xHigh, xLow, yHigh, yLow, scale)
}

func drawRow(rgba *image.RGBA, bits byte, xHigh, xLow, yHigh, yLow, scale byte) {
	if (bits & 0x80) != 0 { drawPixelBlock(rgba, xHigh, xLow + (0 * scale), yHigh, yLow, scale) }
	if (bits & 0x40) != 0 { drawPixelBlock(rgba, xHigh, xLow + (1 * scale), yHigh, yLow, scale) }
	if (bits & 0x20) != 0 { drawPixelBlock(rgba, xHigh, xLow + (2 * scale), yHigh, yLow, scale) }
	if (bits & 0x10) != 0 { drawPixelBlock(rgba, xHigh, xLow + (3 * scale), yHigh, yLow, scale) }
	if (bits & 0x08) != 0 { drawPixelBlock(rgba, xHigh, xLow + (4 * scale), yHigh, yLow, scale) }
	if (bits & 0x04) != 0 { drawPixelBlock(rgba, xHigh, xLow + (5 * scale), yHigh, yLow, scale) }
	if (bits & 0x02) != 0 { drawPixelBlock(rgba, xHigh, xLow + (6 * scale), yHigh, yLow, scale) }
	if (bits & 0x01) != 0 { drawPixelBlock(rgba, xHigh, xLow + (7 * scale), yHigh, yLow, scale) }
}

func drawPixelBlock(rgba *image.RGBA, xHigh, xLow, yHigh, yLow, scale byte) {
	var sy byte
	var sx byte
	for sy = 0; sy < scale; sy++ {
		for sx = 0; sx < scale; sx++ {
			finalX := (int(xHigh) << 8) + int(xLow + sx)
			finalY := (int(yHigh) << 8) + int(yLow + sy)
			
			pixOffset := (finalY * rgba.Stride) + (finalX * 4)

			rgba.Pix[pixOffset+0] = 0   // R
			rgba.Pix[pixOffset+1] = 0   // G
			rgba.Pix[pixOffset+2] = 0   // B
			rgba.Pix[pixOffset+3] = 255 // A
		}
	}
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

				// ИСПРАВЛЕНИЕ ОПЕЧАТКИ: Использованы точные системные имена sz.WidthPx и sz.HeightPx
				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				glCtx.Scissor(0, 0, sz.WidthPx, sz.HeightPx)
				
				glCtx.Enable(gl.SCISSOR_TEST)
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				var xHigh byte = 0
				var xLow1 byte = 40
				var xLow2 byte = 80
				var yHigh byte = 0
				var yLow byte = 120
				var textScale byte = 4

				atlasChain.RenderGlyph(rgba, uiState.Char1, xHigh, xLow1, yHigh, yLow, textScale)
				atlasChain.RenderGlyph(rgba, uiState.Char2, xHigh, xLow2, yHigh, yLow, textScale)

				statusBuffer.Upload()
				statusBuffer.Draw(sz, geom.Point{}, geom.Point{X: sz.WidthPt}, geom.Point{Y: sz.HeightPt}, rgba.Bounds())
				glCtx.Flush()
				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}
