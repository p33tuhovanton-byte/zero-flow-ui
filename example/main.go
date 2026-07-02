package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

type CameraStateHolder struct {
	CurrentProjection ProjectionStrategy
}

type SystemLifecycleLoop struct {
	AppInstance app.App
	EventsChan  <-chan any
	Holder      *CameraStateHolder
	GLContext   gl.Context
	WidthPeano  Number
	HeightPeano Number
	ScanState   Vector
}

func (sll SystemLifecycleLoop) LoopStep() {
	raw := <-sll.EventsChan
	
	evLifecycle, okLifecycle := raw.(lifecycle.Event)
	if okLifecycle {
		glCtx, _ := evLifecycle.DrawContext.(gl.Context)
		SystemLifecycleLoop{AppInstance: sll.AppInstance, EventsChan: sll.EventsChan, Holder: sll.Holder, GLContext: glCtx, WidthPeano: sll.WidthPeano, HeightPeano: sll.HeightPeano, ScanState: sll.ScanState}.LoopStep()
		return
	}
	
	evSize, okSize := raw.(size.Event)
	if okSize {
		newW := intToPeano(okSize.WidthPx, Zero{})
		newH := intToPeano(okSize.HeightPx, Zero{})
		newScan := WavefrontOrientedStrategy{X: Zero{}, Y: Zero{}, MaxX: newW, MaxY: newH, ProjMethod: sll.Holder.CurrentProjection, DirectionDelta: Zero{}.Next()}
		SystemLifecycleLoop{AppInstance: sll.AppInstance, EventsChan: sll.EventsChan, Holder: sll.Holder, GLContext: sll.GLContext, WidthPeano: newW, HeightPeano: newH, ScanState: newScan}.LoopStep()
		return
	}
	
	evTouch, okTouch := raw.(touch.Event)
	if okTouch {
		if evTouch.Type == touch.TypeBegin {
			TouchPulseEvent{StateHolder: sll.Holder}.Trigger()
		}
		SystemLifecycleLoop{AppInstance: sll.AppInstance, EventsChan: sll.EventsChan, Holder: sll.Holder, GLContext: sll.GLContext, WidthPeano: sll.WidthPeano, HeightPeano: sll.HeightPeano, ScanState: sll.ScanState}.LoopStep()
		return
	}
	
	_, okPaint := raw.(paint.Event)
	if okPaint {
		sll.PaintStep()
		return
	}
	
	SystemLifecycleLoop{AppInstance: sll.AppInstance, EventsChan: sll.EventsChan, Holder: sll.Holder, GLContext: sll.GLContext, WidthPeano: sll.WidthPeano, HeightPeano: sll.HeightPeano, ScanState: sll.ScanState}.LoopStep()
}

func (sll SystemLifecycleLoop) PaintStep() {
	ctx := sll.GLContext
	ctx.Enable(gl.SCISSOR_TEST)
	ctx.ClearColor(1.0, 1.0, 1.0, 1.0)
	ctx.Clear(gl.COLOR_BUFFER_BIT)

	container := UniversalContainer[Vector]{}
	NativeGameRenderEvent{GL: ctx, Width: sll.WidthPeano, Height: sll.HeightPeano, Projection: sll.Holder.CurrentProjection, CurrentVec: sll.ScanState, OutVec: container}.Trigger()

	sll.AppInstance.Publish()
	SystemLifecycleLoop{AppInstance: sll.AppInstance, EventsChan: sll.EventsChan, Holder: sll.Holder, GLContext: sll.GLContext, WidthPeano: sll.WidthPeano, HeightPeano: sll.HeightPeano, ScanState: container.Value}.LoopStep()
}

func main() {
	holder := &CameraStateHolder{CurrentProjection: TopViewProjection{}}
	app.Main(func(a app.App) {
		initialScan := WavefrontOrientedStrategy{X: Zero{}, Y: Zero{}, MaxX: Zero{}, MaxY: Zero{}, DirectionDelta: Zero{}.Next()}
		SystemLifecycleLoop{AppInstance: a, EventsChan: a.Events(), Holder: holder, ScanState: initialScan}.LoopStep()
	})
}

func intToPeano(n int, current Number) Number {
	if n <= 0 {
		return current
	}
	return intToPeano(n-1, Successor{pred: current})
}

type OpenGlPixelDriver struct {
	GL       gl.Context
	CounterX int
	CounterY int
	ModeFlag int
}

func (ogpd OpenGlPixelDriver) IdentifyClass() {}
func (ogpd OpenGlPixelDriver) IncrementPulse() HardwareIntegerDriver {
	return OpenGlPixelDriver{GL: ogpd.GL, CounterX: ogpd.CounterX + 1, CounterY: ogpd.CounterY, ModeFlag: ogpd.ModeFlag}
}
func (ogpd OpenGlPixelDriver) ExecuteHardwarePulse() {
	ogpd.GL.Scissor(int32(ogpd.CounterX), int32(ogpd.CounterY), 2, 2)
	ogpd.GL.ClearColor(1.0, 0.0, 0.0, 1.0)
	ogpd.GL.Clear(gl.COLOR_BUFFER_BIT)
}

type GenericColorAcceptor[T GameColor] struct {
	Target UniversalContainer[GameColor]
	Result T
}

func (gca GenericColorAcceptor[T]) IdentifyClass() {}
func (gca GenericColorAcceptor[T]) Execute()       {} 
func (gca GenericColorAcceptor[T]) AcceptColor()   { gca.Execute() }

type TouchPulseEvent struct{ StateHolder *CameraStateHolder }
func (tpe TouchPulseEvent) IdentifyClass() {}
func (tpe TouchPulseEvent) Trigger() {
	tpe.StateHolder.CurrentProjection = tpe.StateHolder.CurrentProjection.NextOrientation()
}

type Canvas interface {
	Object
	ReadColor() Action
}

type NativeGameRenderEvent struct {
	GL         gl.Context
	Width      Number
	Height     Number
	Projection ProjectionStrategy
	CurrentVec Vector
	OutVec     UniversalContainer[Vector]
}

func (ngre NativeGameRenderEvent) IdentifyClass() {}
func (ngre NativeGameRenderEvent) Trigger() {
	CameraNode{
		Projection: ngre.Projection,
		ChildNode: CanvasNode{
			ScanStrategy: ngre.CurrentVec,
			HardwareGL:   ngre.GL,
			OutVec:       ngre.OutVec,
		},
	}.ProcessNode().Execute()
}

type CanvasNode struct {
	ScanStrategy Vector
	HardwareGL   gl.Context
	OutVec       UniversalContainer[Vector]
}
func (cn CanvasNode) IdentifyClass() {}
func (cn CanvasNode) ProcessNode() Action {
	var initialSnapshot Snapshot[GameColor] = EmptySnapshot[GameColor]{}
	CanvasScanner{Step: cn.ScanStrategy, Canvas: OpenGlCanvas{GlContext: cn.HardwareGL}, Storage: initialSnapshot, OutVec: cn.OutVec}.Scan()
	return EmptyAction{}
}

type OpenGlCanvas struct {
	GlContext gl.Context
	Scene     SceneNode
}

func (ogc OpenGlCanvas) IdentifyClass() {}
func (ogc OpenGlCanvas) ReadColor() Action {
	ogc.Scene.RenderPixel()
	return EmptyAction{}
}

type CanvasScanner struct {
	Step    Vector
	Canvas  Canvas
	Storage Snapshot[GameColor]
	OutVec  UniversalContainer[Vector]
}

func (cs CanvasScanner) IdentifyClass() {}

type PixelSaveAcceptor struct {
	Scanner       CanvasScanner
	UpdatedCanvas OpenGlCanvas
	InjectedColor GameColor
}

func (psa PixelSaveAcceptor) IdentifyClass() {}
func (psa PixelSaveAcceptor) Execute()       {}
func (psa PixelSaveAcceptor) AcceptColor() {
	container := UniversalContainer[Snapshot[GameColor]]{}
	activeSnapshot := psa.Scanner.Storage
	
	if node, ok := activeSnapshot.(NodeSnapshot[GameColor]); ok {
		node.TargetNewPoint = SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: psa.InjectedColor}
		node.AcceptorTarget = container
		activeSnapshot = node
	} else if empty, ok := activeSnapshot.(EmptySnapshot[GameColor]); ok {
		empty.TargetNewPoint = SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: psa.InjectedColor}
		empty.AcceptorTarget = container
		activeSnapshot = empty
	}

	activeSnapshot.MutateSnapshotState()

	CanvasScanner{
		Step:    psa.Scanner.Step.AdvanceVector(),
		Canvas:  psa.UpdatedCanvas,
		Storage: container.Value,
		OutVec:  psa.Scanner.OutVec,
	}.Scan()
}

func (cs CanvasScanner) Scan() {
	saveAcceptor := PixelSaveAcceptor{Scanner: cs, UpdatedCanvas: OpenGlCanvas{GlContext: cs.Canvas.(OpenGlCanvas).GlContext}}
	glCtx := cs.Canvas.(OpenGlCanvas).GlContext

	scene := SceneNode{
		Background: WhiteBackgroundLayer{Output: GenericColorAcceptor[GameColor]{Target: UniversalContainer[GameColor]{}, Result: SolidWhiteColor{}}},
		Grid: CoordinateGridLayer{CurrentStep: cs.Step, Output: GenericColorAcceptor[GameColor]{Target: UniversalContainer[GameColor]{}, Result: GridLineColor{
			DriverX: OpenGlPixelDriver{GL: glCtx},
			DriverY: OpenGlPixelDriver{GL: glCtx},
		}}},
		Object3D: ThreeDimensionalObjectLayer{CurrentStep: cs.Step, Output: GenericColorAcceptor[GameColor]{Target: UniversalContainer[GameColor]{}, Result: Object3DColor{
			DriverX: OpenGlPixelDriver{GL: glCtx},
			DriverY: OpenGlPixelDriver{GL: glCtx},
		}}},
		FinalOutput: saveAcceptor,
	}
	saveAcceptor.UpdatedCanvas.Scene = scene
	
	BranchFactory{Condition: cs.Step.IsCanvasFinished(), TrueBranch: DirectAction[Vector]{Target: cs.OutVec, Result: cs.Step}, FalseBranch: ScanAction{Scanner: cs}}.Create().Select().Execute()
}

type Vector interface {
	Object
	AdvanceVector() Vector
	IsCanvasFinished() Bool
	IsGridIntersection() Bool
	IsIntersecting3D() Bool
}

type WavefrontOrientedStrategy struct {
	X              Number
	Y              Number
	MaxX           Number
	MaxY           Number
	ProjMethod     ProjectionStrategy
	DirectionDelta Number 
}

func (wos WavefrontOrientedStrategy) IdentifyClass() {}

func (wos WavefrontOrientedStrategy) AdvanceVector() Vector {
	container := UniversalContainer[Vector]{}
	midThreshold := Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()

	BranchFactory{
		Condition: wos.X.EvaluateWaveCenter(wos.MaxX, midThreshold),
		TrueBranch: DirectAction[Vector]{Target: container, Result: WavefrontOrientedStrategy{X: wos.X.Next(), Y: wos.Y, MaxX: wos.MaxX, MaxY: wos.MaxY, ProjMethod: wos.ProjMethod, DirectionDelta: Zero{}.Next()}},
		FalseBranch: BranchFactory{
			Condition:  Zero{CompareTarget: wos.X.Next()}.CheckEquality(),
			TrueBranch: DirectAction[Vector]{Target: container, Result: WavefrontOrientedStrategy{X: Zero{}, Y: wos.Y.Next(), MaxX: wos.MaxX, MaxY: wos.MaxY, ProjMethod: wos.ProjMethod, DirectionDelta: Zero{}.Next()}},
			FalseBranch: DirectAction[Vector]{Target: container, Result: WavefrontOrientedStrategy{X: wos.X.Next(), Y: wos.Y, MaxX: wos.MaxX, MaxY: wos.MaxY, ProjMethod: wos.ProjMethod, DirectionDelta: wos.DirectionDelta}},
		}.Create().Select(),
	}.Create().Select().Execute()
	return container.Value
}

func (wos WavefrontOrientedStrategy) IsCanvasFinished() Bool   { return Zero{CompareTarget: wos.Y}.CheckEquality() }
func (wos WavefrontOrientedStrategy) IsGridIntersection() Bool { return wos.X.IsMultipleOfGrid() }
