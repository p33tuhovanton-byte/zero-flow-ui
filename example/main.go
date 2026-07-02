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
		var glCtx gl.Context
		var screenWidth, screenHeight int
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				glCtx, _ = e.DrawContext.(gl.Context)
			case size.Event:
				screenWidth, screenHeight = e.WidthPx, e.HeightPx
			case touch.Event:
				if e.Type == touch.TypeBegin {
					TouchPulseEvent{StateHolder: holder}.Trigger()
				}
			case paint.Event:
				if glCtx == nil {
					continue
				}
				InitPeanoFactory{LimitX: screenWidth, LimitY: screenHeight, OnReady: GameLauncherAcceptor{GL: glCtx, Holder: holder}}.StartBuild()
				a.Publish()
			}
		}
	})
}

type InitPeanoFactory struct {
	LimitX, LimitY  int
	Current, SavedX Number
	OnReady         GameLauncherAcceptor
}

func (ipf InitPeanoFactory) StartBuild() { ipf.BuildX() }
func (ipf InitPeanoFactory) BuildX() {
	if ipf.LimitX <= 0 {
		InitPeanoFactory{LimitY: ipf.LimitY, Current: Zero{}, SavedX: ipf.Current, OnReady: ipf.OnReady}.BuildY()
		return
	}
	InitPeanoFactory{LimitX: ipf.LimitX - 1, LimitY: ipf.LimitY, Current: Successor{pred: ipf.Current}, OnReady: ipf.OnReady}.BuildX()
}
func (ipf InitPeanoFactory) BuildY() {
	if ipf.LimitY <= 0 {
		ipf.OnReady.Launch()
		return
	}
	InitPeanoFactory{LimitY: ipf.LimitY - 1, Current: Successor{pred: ipf.Current}, SavedX: ipf.SavedX, OnReady: ipf.OnReady}.BuildY()
}

type GameLauncherAcceptor struct {
	GL        gl.Context
	Holder    *CameraStateHolder
	WidthNum  Number
	HeightNum Number
}

func (gla GameLauncherAcceptor) Launch() {
	NativeGameRenderEvent{GL: gla.GL, Width: gla.WidthNum, Height: gla.HeightNum, Projection: gla.Holder.CurrentProjection}.Trigger()
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
	CanvasScanner{
		Step:    HorizontalRowStrategy{X: Zero{}, Y: Zero{}, MaxX: ngre.Width, MaxY: ngre.Height, ProjMethod: ngre.Projection},
		Canvas:  OpenGlCanvas{GlContext: ngre.GL},
		Storage: EmptySnapshot[GameColor]{},
	}.Scan()
}

type OpenGlCanvas struct {
	GlContext gl.Context
	Scene     Composited3DScene
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
}

func (psa PixelSaveAcceptor) AcceptColor() {
	// ЭКРАН СТАЛ БЕЛЫМ: Теперь дефолтный белый цвет инжектируется принудительно
	CanvasScanner{
		Step:   psa.Scanner.Step.AdvanceVector(),
		Canvas: psa.UpdatedCanvas,
		Storage: NodeSnapshot[GameColor]{
			tail:     psa.Scanner.Storage,
			NewPoint: SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: SolidWhiteColor{}},
		}.Accumulate(),
	}.Scan()
}

type DirectColorAction struct {
	Target *PixelSaveAcceptor
	Color  GameColor
}
func (dca DirectColorAction) IdentifyClass() {}
func (dca DirectColorAction) Execute()       { dca.Target.InjectedColor = dca.Color; dca.Target.AcceptColor() }

func (cs CanvasScanner) Scan() {
	saveAcceptor := PixelSaveAcceptor{Scanner: cs, UpdatedCanvas: OpenGlCanvas{GlContext: cs.Canvas.(OpenGlCanvas).GlContext}}
	
	// Вшиваем триггеры разрешения слоев (Куб -> Сетка -> Белый Фон)
	scene := Composited3DScene{
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

type HorizontalRowStrategy struct {
	X          Number
	Y          Number
	MaxX       Number
	MaxY       Number
	ProjMethod ProjectionStrategy
}

func (hrs HorizontalRowStrategy) IdentifyClass() {}

type VectorContainer struct{ Value Vector }

func (hrs HorizontalRowStrategy) AdvanceVector() Vector {
	container := &VectorContainer{}
	BranchFactory{
		Condition:   Zero{CompareTarget: hrs.X.Next()}.CheckEquality(),
		TrueBranch:  DirectVectorAction{Target: container, Result: HorizontalRowStrategy{X: Zero{}, Y: hrs.Y.Next(), MaxX: hrs.MaxX, MaxY: hrs.MaxY, ProjMethod: hrs.ProjMethod}},
		FalseBranch: DirectVectorAction{Target: container, Result: HorizontalRowStrategy{X: hrs.X.Next(), Y: hrs.Y, MaxX: hrs.MaxX, MaxY: hrs.MaxY, ProjMethod: hrs.ProjMethod}},
	}.Create().Select().Execute()
	return container.Value
}

type DirectVectorAction struct {
	Target *VectorContainer
	Result Vector
}

func (dva DirectVectorAction) IdentifyClass() {}
func (dva DirectVectorAction) Execute()       { dva.Target.Value = dva.Result }

func (hrs HorizontalRowStrategy) IsCanvasFinished() Bool   { return Zero{CompareTarget: hrs.Y}.CheckEquality() }
func (hrs HorizontalRowStrategy) IsGridIntersection() Bool { return hrs.X.IsMultipleOfGrid() }

type CubeIntersectionAcceptor struct {
	ScannerCoords  HorizontalRowStrategy
	ResultTarget   *BoolContainer
	ProjectedPoint Vector2D
}

func (cia CubeIntersectionAcceptor) AcceptProjection() {
	cia.ResultTarget.Value = cia.ProjectedPoint.U.CheckEquality()
}

type BoolContainer struct{ Value Bool }

func (hrs HorizontalRowStrategy) IsIntersecting3D() Bool {
	container := &BoolContainer{Value: False{}}

	var dynamicProjector ProjectionStrategy
	dynamicProjector = hrs.ProjMethod
	
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
