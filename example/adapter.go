package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

// CameraStateHolder инкапсулирует текущую активную проекцию нативной камеры
type CameraStateHolder struct {
	CurrentProjection ProjectionStrategy
}

// ============================================================================
// АППАРАТНАЯ ТОЧКА ВХОДА (Изолированный нативный мост Android)
// ============================================================================

func main() {
	holder := &CameraStateHolder{CurrentProjection: TopViewProjection{}}
	var startScanVector Vector = WavefrontOrientedStrategy{
		X:              Zero{},
		Y:              Zero{},
		MaxX:           Zero{},
		MaxY:           Zero{},
		DirectionDelta: Zero{}.Next(),
	}

	app.Main(func(a app.App) {
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

// OpenGlPixelDriver осуществляет физический сброс бинарных команд в видеокарту.
// Он накапливает в себе два счетчика (X и Y) по цепочке Пеано без if-условий.
type OpenGlPixelDriver struct {
	GL       gl.Context
	CounterX int
	CounterY int
	ViewportW int
	ViewportH int
	ModeFlag  int // 0 - Scissor Grid, 1 - Scissor Object, 2 - Viewport Canvas
}

func (ogpd OpenGlPixelDriver) IdentifyClass() {}
func (ogpd OpenGlPixelDriver) IncrementPulse() HardwareIntegerDriver {
	// Базовый инкремент наращивает первый фокусный счетчик X
	return OpenGlPixelDriver{
		GL:        ogpd.GL,
		CounterX:  ogpd.CounterX + 1,
		CounterY:  ogpd.CounterY,
		ViewportW: ogpd.ViewportW,
		ViewportH: ogpd.ViewportH,
		ModeFlag:  ogpd.ModeFlag,
	}
}

func (ogpd OpenGlPixelDriver) IncrementSecondPulse() OpenGlPixelDriver {
	// Дополнительный инкремент наращивает счетчик Y
	return OpenGlPixelDriver{
		GL:        ogpd.GL,
		CounterX:  ogpd.CounterX,
		CounterY:  ogpd.CounterY + 1,
		ViewportW: ogpd.ViewportW,
		ViewportH: ogpd.ViewportH,
		ModeFlag:  ogpd.ModeFlag,
	}
}

func (ogpd OpenGlPixelDriver) IncrementWidth() OpenGlPixelDriver {
	return OpenGlPixelDriver{
		GL:        ogpd.GL,
		CounterX:  ogpd.CounterX,
		CounterY:  ogpd.CounterY,
		ViewportW: ogpd.ViewportW + 1,
		ViewportH: ogpd.ViewportH,
		ModeFlag:  ogpd.ModeFlag,
	}
}

func (ogpd OpenGlPixelDriver) IncrementHeight() OpenGlPixelDriver {
	return OpenGlPixelDriver{
		GL:        ogpd.GL,
		CounterX:  ogpd.CounterX,
		CounterY:  ogpd.CounterY,
		ViewportW: ogpd.ViewportW,
		ViewportH: ogpd.ViewportH + 1,
		ModeFlag:  ogpd.ModeFlag,
	}
}

func (ogpd OpenGlPixelDriver) ExecuteHardwarePulse() {
	// Логика выбора режима рендеринга изолирована на аппаратной границе сред
	if ogpd.ModeFlag == 2 {
		ogpd.GL.Viewport(0, 0, ogpd.ViewportW, ogpd.ViewportH)
		return
	}
	if ogpd.ModeFlag == 1 {
		ogpd.GL.Scissor(int32(ogpd.CounterX), int32(ogpd.CounterY), 2, 2)
		ogpd.GL.ClearColor(1.0, 0.0, 0.0, 1.0)
		ogpd.GL.Clear(gl.COLOR_BUFFER_BIT)
		return
	}
	ogpd.GL.Scissor(int32(ogpd.CounterX), int32(ogpd.CounterY), 2, 2)
	ogpd.GL.ClearColor(0.8, 0.8, 0.8, 1.0)
	ogpd.GL.Clear(gl.COLOR_BUFFER_BIT)
}
