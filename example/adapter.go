package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

// ============================================================================
// АППАРАТНАЯ ТОЧКА ВХОДА (Изолированный нативный мост Android)
// ============================================================================

func main() {
	// Создаем стартовую конфигурацию камеры и волнового вектора сканера
	holder := &CameraStateHolder{CurrentProjection: TopViewProjection{}}
	var startScanVector Vector = WavefrontOrientedStrategy{
		X:              Zero{},
		Y:              Zero{},
		MaxX:           Zero{},
		MaxY:           Zero{},
		DirectionDelta: Zero{}.Next(),
	}

	app.Main(func(a app.App) {
		// Запускаем вечный реактивный автомат, полностью изолирующий канал событий
		runHardwareLifecycleLoop(a, a.Events(), holder, nil, Zero{}, Zero{}, startScanVector)
	})
}

func runHardwareLifecycleLoop(a app.App, events <-chan any, holder *CameraStateHolder, ctx gl.Context, w Number, h Number, scanState Vector) {
	raw := <-events
	evLifecycle, okLifecycle := raw.(lifecycle.Event)
	if okLifecycle {
		glCtx, _ := evLifecycle.DrawContext.(gl.Context)
		runHardwareLifecycleLoop(a, events, holder, glCtx, w, h, scanState)
		return
	}
	evSize, okSize := raw.(size.Event)
	if okSize {
		newW := intToPeano(evSize.WidthPx, Zero{})
		newH := intToPeano(evSize.HeightPx, Zero{})
		updatedScan := WavefrontOrientedStrategy{X: Zero{}, Y: Zero{}, MaxX: newW, MaxY: newH, ProjMethod: holder.CurrentProjection, DirectionDelta: Zero{}.Next()}
		runHardwareLifecycleLoop(a, events, holder, ctx, newW, newH, updatedScan)
		return
	}
	evTouch, okTouch := raw.(touch.Event)
	if okTouch {
		if evTouch.Type == touch.TypeBegin {
			TouchPulseEvent{StateHolder: holder}.Trigger()
		}
		runHardwareLifecycleLoop(a, events, holder, ctx, w, h, scanState)
		return
	}
	_, okPaint := raw.(paint.Event)
	if okPaint {
		if ctx != nil {
			ctx.Enable(gl.SCISSOR_TEST)
			ctx.ClearColor(1.0, 1.0, 1.0, 1.0)
			ctx.Clear(gl.COLOR_BUFFER_BIT)

			// ПОЛИМОРФНЫЙ ПЕРЕХОД: Передаем управление в Generic-мод нашей игры.
			// Контекст OpenGL упаковывается в безопасный типизированный контейнер.
			container := &UniversalContainer[Vector]{}
			GameModLauncher{
				GL:         ctx,
				Width:      w,
				Height:     h,
				Projection: holder.CurrentProjection,
				CurrentVec: scanState,
				OutVec:     container,
			}.LaunchMod()

			a.Publish()
			runHardwareLifecycleLoop(a, events, holder, ctx, w, h, container.Value)
			return
		}
		runHardwareLifecycleLoop(a, events, holder, ctx, w, h, scanState)
		return
	}
	runHardwareLifecycleLoop(a, events, holder, ctx, w, h, scanState)
}

func intToPeano(n int, current Number) Number {
	if n <= 0 {
		return current
	}
	return intToPeano(n-1, Successor{pred: current})
}

// OpenGlPixelDriver осуществляет физический сброс бинарных команд в видеокарту
type OpenGlPixelDriver struct {
	GL      gl.Context
	Counter int
	IsYAxis bool 
}

func (ogpd OpenGlPixelDriver) IdentifyClass() {}
func (ogpd OpenGlPixelDriver) IncrementPulse() HardwareIntegerDriver {
	return OpenGlPixelDriver{GL: ogpd.GL, Counter: ogpd.Counter + 1, IsYAxis: ogpd.IsYAxis}
}
func (ogpd OpenGlPixelDriver) ExecuteHardwarePulse() {
	ogpd.GL.Scissor(int32(ogpd.Counter), int32(ogpd.Counter), 2, 2)
	ogpd.GL.ClearColor(1.0, 0.0, 0.0, 1.0)
	ogpd.GL.Clear(gl.COLOR_BUFFER_BIT)
}
