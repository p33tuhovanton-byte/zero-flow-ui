package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch" // Пакет для обработки тапов
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"zeroflowui"
)

// Фиксированная матрица шрифта 8x8 для Zero-Alloc растеризации.
// Содержит базовые символы для демонстрации статусов логов (Z, e, r, o, F, l, w, U, I, :, K).
// 1 - пиксель текста, 0 - фон. Массив статичен и не аллоцирует память.
var zeroFont8x8 = map[rune][8]byte{
	'Z': {0x7E, 0x0C, 0x18, 0x30, 0x60, 0x42, 0x7E, 0x00},
	'e': {0x3C, 0x42, 0x7E, 0x40, 0x40, 0x42, 0x3C, 0x00},
	'r': {0x5C, 0x62, 0x40, 0x40, 0x40, 0x40, 0xE0, 0x00},
	'o': {0x3C, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'F': {0x7E, 0x40, 0x40, 0x78, 0x40, 0x40, 0x40, 0x00},
	'l': {0x60, 0x60, 0x60, 0x60, 0x60, 0x60, 0x7E, 0x00},
	'w': {0x42, 0x42, 0x42, 0x4A, 0x54, 0x64, 0x42, 0x00},
	'U': {0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'I': {0x3C, 0x18, 0x18, 0x18, 0x18, 0x18, 0x3C, 0x00},
	':': {0x00, 0x18, 0x18, 0x00, 0x18, 0x18, 0x00, 0x00},
	'O': {0x3C, 0x42, 0x42, 0x42, 0x42, 0x42, 0x3C, 0x00},
	'K': {0x42, 0x44, 0x48, 0x70, 0x48, 0x44, 0x42, 0x00},
	' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
}

func main() {
	uiTimeline := zeroflowui.EndOfUI()
	uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Rendered")

	// Начальный статус
	textSignal := zeroflowui.TextSignal{
		Type:    zeroflowui.TextType,
		Payload: "ZeroFlowUI: Wait...",
	}

	pipeline := zeroflowui.SystemPipelineDecorator{
		Next: zeroflowui.TerminalProcessor{},
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
				case lifecycle.StageFocused:
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
				case lifecycle.StageAlive:
					if ctx, ok := x.DrawContext.(gl.Context); ok {
						glCtx = ctx
						images = glutil.NewImages(glCtx)
					}
					a.Send(paint.Event{})
				case lifecycle.StageDead:
					if statusBuffer != nil {
						statusBuffer.Release()
						statusBuffer = nil
					}
					if images != nil {
						images.Release()
						images = nil
					}
					glCtx = nil
				}

			case size.Event:
				sz = x
				if glCtx != nil && images != nil {
					if statusBuffer != nil {
						statusBuffer.Release()
					}
					statusBuffer = images.NewImage(sz.WidthPx, sz.HeightPx)
				}
				a.Send(paint.Event{})

			case touch.Event:
				// ОБРАБОТКА ТАПОВ: Меняем полезную нагрузку сигнала по нажатию на экран
				if x.Type == touch.TypeBegin {
					// Имитируем декларативное уведомление от конвейера
					textSignal.Payload = "ZeroFlowUI: OK"
					// Прокачиваем событие через вашу библиотеку
					uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "AndroidMainWindow", "Touched")
					pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
					// Принудительно запрашиваем перерисовку экрана
					a.Send(paint.Event{})
				}

			case paint.Event:
				if glCtx == nil || images == nil || statusBuffer == nil {
					a.Send(paint.Event{})
					continue
				}

				// Аппаратно фиксируем область вывода на полный экран
				glCtx.Viewport(0, 0, sz.WidthPx, sz.HeightPx)
				
				// Полная очистка холста GPU в чистый белый цвет
				glCtx.ClearColor(1.0, 1.0, 1.0, 1.0)
				glCtx.Clear(gl.COLOR_BUFFER_BIT)

				rgba := statusBuffer.RGBA
				// Быстрая заливка подложки текстуры белым цветом
				draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

				// ZERO-ALLOC РЕНДЕРИНГ: Выводим текст из Payload на белый фон черным цветом
				drawZeroAllocText(rgba, textSignal.Payload, 40, 100, 4) // Масштаб шрифта x4 для читаемости

				statusBuffer.Upload()
				
				// Отрисовка текстуры фреймбуфера
				statusBuffer.Draw(
					sz,
					geom.Point{X: 0, Y: 0},
					geom.Point{X: sz.WidthPt, Y: 0},
					geom.Point{X: 0, Y: sz.HeightPt},
					rgba.Bounds(),
				)

				a.Publish()
				a.Send(paint.Event{})
			}
		}
	})
}

// drawZeroAllocText попиксельно переносит символы напрямую в массив байт фреймбуфера.
// Не выделяет память (0 allocations/op).
func drawZeroAllocText(rgba *image.RGBA, text string, startX, startY, scale int) {
	currentX := startX
	
	for _, char := range text {
		bitmap, exists := zeroFont8x8[char]
		if !exists {
			bitmap = zeroFont8x8[' '] // Если символ не найден — рисуем пробел
		}

		// Перенос битовой карты символа в пиксели image.RGBA
		for row := 0; row < 8; row++ {
			bits := bitmap[row]
			for col := 0; col < 8; col++ {
				// Проверяем, установлен ли бит пикселя текста
				if (bits & (1 << (7 - uint(col)))) != 0 {
					// Отрисовываем пиксель с учетом масштабирования (scale)
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							px := currentX + col*scale + sx
							py := startY + row*scale + sy
							
							// Красим пиксель в черный цвет
							rgba.SetRGBA(px, py, color.RGBA{R: 0, G: 0, B: 0, A: 255})
						}
					}
				}
			}
		}
		// Сдвигаем координату X для следующего символа (8 пикселей + 2 пикселя отступ)
		currentX += 10 * scale
	}
}
