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

func main() {
	holder := &CameraStateHolder{CurrentProjection: TopViewProjection{}}
	app.Main(func(a app.App) {
		runLifecycleLoop(a.Events(), holder, nil, Zero{}, Zero{})
	})
}

func runLifecycleLoop(events <-chan any, holder *CameraStateHolder, ctx gl.Context, w Number, h Number) {
	raw := <-events
	evLifecycle, okLifecycle := raw.(lifecycle.Event)
	if okLifecycle {
		glCtx, _ := evLifecycle.DrawContext.(gl.Context)
		AppLifecycleLoop{StateHolder: holder, GLContext: glCtx, WidthNum: w, HeightNum: h}.DispatchLifecycle()
		runLifecycleLoop(events, holder, glCtx, w, h)
		return
	}
	evSize, okSize := raw.(size.Event)
	if okSize {
		newW := intToPeano(evSize.WidthPx, Zero{})
		newH := intToPeano(evSize.HeightPx, Zero{})
		runLifecycleLoop(events, holder, ctx, newW, newH)
		return
	}
	evTouch, okTouch := raw.(touch.Event)
	if okTouch {
		if evTouch.Type == touch.TypeBegin {
			TouchPulseEvent{StateHolder: holder}.Trigger()
		}
		runLifecycleLoop(events, holder, ctx, w, h)
		return
	}
	_, okPaint := raw.(paint.Event)
	if okPaint {
		AppLifecycleLoop{StateHolder: holder, GLContext: ctx, WidthNum: w, HeightNum: h}.DispatchPaint()
		runLifecycleLoop(events, holder, ctx, w, h)
		return
	}
	runLifecycleLoop(events, holder, ctx, w, h)
}

func intToPeano(n int, current Number) Number {
	if n <= 0 {
		return current
	}
	return intToPeano(n-1, Successor{pred: current})
}

type AppLifecycleLoop struct {
	StateHolder *CameraStateHolder
	GLContext   gl.Context
	WidthNum    Number
	HeightNum   Number
}

func (all AppLifecycleLoop) DispatchLifecycle() {}
func (all AppLifecycleLoop) DispatchSize()      {}
func (all AppLifecycleLoop) DispatchPaint() {
	if all.GLContext != nil {
		NativeGameRenderEvent{
			GL:         all.GLContext,
			Width:      all.WidthNum,
			Height:     all.HeightNum,
			Projection: all.StateHolder.CurrentProjection,
		}.Trigger()
	}
}

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
}

func (ngre NativeGameRenderEvent) IdentifyClass() {}
func (ngre NativeGameRenderEvent) Trigger() {
	ngre.GL.ClearColor(1.0, 1.0, 1.0, 1.0)
	ngre.GL.Clear(gl.COLOR_BUFFER_BIT)

	canvasNode := CanvasNode{
		ScanStrategy: WavefrontOrientedStrategy{X: Zero{}, Y: Zero{}, MaxX: ngre.Width, MaxY: ngre.Height, ProjMethod: ngre.Projection, DirectionDelta: Zero{}.Next()},
		HardwareGL:   ngre.GL,
	}
	
	engineTree := CameraNode{
		Projection: ngre.Projection,
		ChildNode:  canvasNode,
	}

	engineTree.ProcessNode().Execute()
}

type CanvasNode struct {
	ScanStrategy Vector
	HardwareGL   gl.Context
}
func (cn CanvasNode) IdentifyClass() {}
func (cn CanvasNode) ProcessNode() Action {
	CanvasScanner{
		Step:    cn.ScanStrategy,
		Canvas:  OpenGlCanvas{GlContext: cn.HardwareGL},
		Storage: EmptySnapshot[GameColor]{},
	}.Scan()
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

type Scanner interface {
	Object
	Scan()
}

type CanvasScanner struct {
	Step    Vector
	Canvas  Canvas
	Storage Snapshot[GameColor]
}

func (cs CanvasScanner) IdentifyClass() {}

type PixelSaveAcceptor struct {
	Scanner       CanvasScanner
	UpdatedCanvas OpenGlCanvas
	InjectedColor GameColor
}

func (psa PixelSaveAcceptor) AcceptColor() {
	psa.InjectedColor.PaintHardwarePixel()

	CanvasScanner{
		Step:   psa.Scanner.Step.AdvanceVector(),
		Canvas: psa.UpdatedCanvas,
		Storage: NodeSnapshot[GameColor]{
			tail:     psa.Scanner.Storage,
			NewPoint: SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: psa.InjectedColor},
		}.Accumulate(),
	}.Scan()
}

// Избавились от кастомного DirectColorAction, заменив на обобщенную команду DirectAction
func (cs CanvasScanner) Scan() {
	saveAcceptor := PixelSaveAcceptor{Scanner: cs, UpdatedCanvas: OpenGlCanvas{GlContext: cs.Canvas.(OpenGlCanvas).GlContext}}

	scene := SceneNode{
		Background:  WhiteBackgroundLayer{Output: saveAcceptor},
		Grid:        CoordinateGridLayer{CurrentStep: cs.Step, Output: saveAcceptor},
		Object3D:    ThreeDimensionalObjectLayer{CurrentStep: cs.Step, Output: saveAcceptor},
		FinalOutput: saveAcceptor,
	}
	saveAcceptor.UpdatedCanvas.Scene = scene
	BranchFactory{Condition: cs.Step.IsCanvasFinished(), TrueBranch: StopAction{FinalSnapshot: cs.Storage}, FalseBranch: ScanAction{Scanner: cs}}.Create().Select().Execute()
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
	// ИСПРАВЛЕНО: Кастомный VectorContainer заменен на UniversalContainer[Vector]
	container := &UniversalContainer[Vector]{}
	midThreshold := Zero{}.Next().Next().Next().Next().Next()

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

// ИСПРАВЛЕНО: Кастомный CubeIntersectionAcceptor полностью удален за ненадобностью!
type WavefrontIntersectionAcceptor struct {
	ResultTarget   *UniversalContainer[Bool]
	ProjectedPoint Vector2D
}
func (wia WavefrontIntersectionAcceptor) AcceptProjection() { wia.ResultTarget.Value = wia.ProjectedPoint.U.CheckEquality() }

func (wos WavefrontOrientedStrategy) IsIntersecting3D() Bool {
	container := &UniversalContainer[Bool]{Value: False{}}
	cubeEdgeDistance := wos.X.Differentiate(Zero{}.Next().Next().Next(), Zero{})
	container.Value = cubeEdgeDistance.CheckEquality()
	
	var dynamicProjector ProjectionStrategy
	dynamicProjector = wos.ProjMethod
	dynamicProjector.InjectContinuation()
	dynamicProjector.Project()
	return container.Value
}

type ScanAction struct{ Scanner CanvasScanner }
func (sa ScanAction) IdentifyClass() {}
func (sa ScanAction) Execute()       { sa.Scanner.Canvas.ReadColor() }

type StopAction struct{ FinalSnapshot Snapshot[GameColor] }
func (sa StopAction) IdentifyClass() {}
func (sa StopAction) Execute()       {}

type Snapshot[T Object] interface {
	Object
	Accumulate() Snapshot[T]
}

type EmptySnapshot[T Object] struct{ NewPoint Point[T] }
func (es EmptySnapshot[T]) IdentifyClass()         {}
func (es EmptySnapshot[T]) Accumulate() Snapshot[T] { return NodeSnapshot[T]{head: es.NewPoint, tail: es} }

type NodeSnapshot[T Object] struct {
	head     Point[T]
	tail     Snapshot[T]
	NewPoint Point[T]
}
func (ns NodeSnapshot[T]) IdentifyClass()         {}
func (ns NodeSnapshot[T]) Accumulate() Snapshot[T] { return NodeSnapshot[T]{head: ns.NewPoint, tail: ns} }

type Point[T Object] interface{ Object }
type SnapshotPoint[T Object] struct {
	VectorState Vector
	Color       T
}
func (sp SnapshotPoint[T]) IdentifyClass() {}
