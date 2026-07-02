package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

type Object interface {
	IdentifyClass() // Идеально пустая сигнатура
}

type ClassConsumer interface {
	AcceptClassName() // Идеально пустая сигнатура
}

type Action interface {
	Object
	Execute()
}

type Bool interface {
	Object
	Select() Action
}

type EmptyAction struct{}

func (ea EmptyAction) IdentifyClass() {}
func (ea EmptyAction) Execute()       {}

type Number interface {
	Object
	Next() Number
	CheckEquality() Bool
	CompareWithZero() Bool
	CompareWithSuccessor() Bool
	IsMultipleOfGrid() Bool
}

type Zero struct{ CompareTarget Number }

func (z Zero) IdentifyClass()       {}
func (z Zero) Class() string        { return "Zero" }
func (z Zero) Next() Number         { return Successor{pred: z} }
func (z Zero) CheckEquality() Bool  { return z.CompareTarget.CompareWithZero() }
func (z Zero) CompareWithZero() Bool { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) CompareWithSuccessor() Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) IsMultipleOfGrid() Bool { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }

type Successor struct {
	pred          Number
	CompareTarget Number
}

func (s Successor) IdentifyClass()       {}
func (s Successor) Class() string        { return "Successor" }
func (s Successor) Next() Number         { return Successor{pred: s} }
func (s Successor) CheckEquality() Bool  { return s.CompareTarget.CompareWithSuccessor() }
func (s Successor) CompareWithZero() Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (s Successor) CompareWithSuccessor() Bool { return Successor{pred: s.pred, CompareTarget: s.CombinePredecessors()}.CheckEquality() }
func (s Successor) CombinePredecessors() Number { return s.CompareTarget.(Successor).pred }
func (s Successor) IsMultipleOfGrid() Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }

type True struct{ TrueBranch, FalseBranch Action }

func (t True) IdentifyClass()       {}
func (t True) Class() string        { return "True" }
func (t True) Select() Action       { return t.TrueBranch }

type False struct{ TrueBranch, FalseBranch Action }

func (f False) IdentifyClass()       {}
func (f False) Class() string        { return "False" }
func (f False) Select() Action       { return f.FalseBranch }

type BranchFactory struct {
	Condition   Bool
	TrueBranch  Action
	FalseBranch Action
}

type BoolResultContainer struct{ Value Bool }

func (bf BranchFactory) Create() Bool {
	container := &BoolResultContainer{}
	TypeResolver{ClassName: bf.Condition.Class(), T: bf.TrueBranch, F: bf.FalseBranch, Target: container}.Resolve()
	return container.Value
}

type TypeResolver struct {
	ClassName string
	T, F      Action
	Target    *BoolResultContainer
}

func (tr TypeResolver) Resolve() {
	tr.Target.Value = True{TrueBranch: tr.T, FalseBranch: tr.F}
}

type CameraStateHolder struct{ CurrentProjection ProjectionStrategy }

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
				if glCtx == nil { continue }
				InitPeanoFactory{LimitX: screenWidth, LimitY: screenHeight, OnReady: GameLauncherAcceptor{GL: glCtx, Holder: holder}}.StartBuild()
				a.Publish()
			}
		}
	})
}

type InitPeanoFactory struct {
	LimitX, LimitY int
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
		ipf.OnReady.Launch(ipf.SavedX, ipf.Current)
		return
	}
	InitPeanoFactory{LimitY: ipf.LimitY - 1, Current: Successor{pred: ipf.Current}, SavedX: ipf.SavedX, OnReady: ipf.OnReady}.BuildY()
}

type GameLauncherAcceptor struct {
	GL     gl.Context
	Holder *CameraStateHolder
}
func (gla GameLauncherAcceptor) Launch(w Number, h Number) {
	NativeGameRenderEvent{GL: gla.GL, Width: w, Height: h, Projection: gla.Holder.CurrentProjection}.Trigger()
}

type TouchPulseEvent struct{ StateHolder *CameraStateHolder }
func (tpe TouchPulseEvent) IdentifyClass(consumer ClassConsumer) {}
func (tpe TouchPulseEvent) Trigger() { tpe.StateHolder.CurrentProjection = tpe.StateHolder.CurrentProjection.NextOrientation() }

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
func (ngre NativeGameRenderEvent) IdentifyClass(consumer ClassConsumer) {}
func (ngre NativeGameRenderEvent) Trigger() {
	CanvasScanner{
		Step: HorizontalRowStrategy{X: Zero{}, Y: Zero{}, MaxX: ngre.Width, MaxY: ngre.Height, ProjMethod: ngre.Projection},
		Canvas: OpenGlCanvas{GlContext: ngre.GL},
		Storage: EmptySnapshot[GameColor]{},
	}.Scan()
}

type OpenGlCanvas struct {
	GlContext gl.Context
	Scene     Composited3DScene
}
func (ogc OpenGlCanvas) IdentifyClass(consumer ClassConsumer) {}
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
func (cs CanvasScanner) IdentifyClass(consumer ClassConsumer) {}

type PixelSaveAcceptor struct {
	Scanner       CanvasScanner
	UpdatedCanvas OpenGlCanvas
}
func (psa PixelSaveAcceptor) AcceptColor(color GameColor) {
	CanvasScanner{
		Step: psa.Scanner.Step.AdvanceVector(),
		Canvas: psa.UpdatedCanvas,
		Storage: NodeSnapshot[GameColor]{tail: psa.Scanner.Storage, NewPoint: SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: color}}.Accumulate(),
	}.Scan()
}
func (cs CanvasScanner) Scan() {
	saveAcceptor := PixelSaveAcceptor{Scanner: cs, UpdatedCanvas: OpenGlCanvas{GlContext: cs.Canvas.(OpenGlCanvas).GlContext}}
	scene := Composited3DScene{
		Background: WhiteBackgroundLayer{Output: saveAcceptor},
		Grid: CoordinateGridLayer{CurrentStep: cs.Step, Output: saveAcceptor},
		Object3D: ThreeDimensionalObjectLayer{CurrentStep: cs.Step, Output: saveAcceptor},
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
func (hrs HorizontalRowStrategy) IdentifyClass(consumer ClassConsumer) {}

type VectorContainer struct{ Value Vector }

func (hrs HorizontalRowStrategy) AdvanceVector() Vector {
	container := &VectorContainer{}
	BranchFactory{
		Condition: Zero{CompareTarget: hrs.X.Next()}.CheckEquality(),
		TrueBranch: DirectVectorAction{Target: container, Result: HorizontalRowStrategy{X: Zero{}, Y: hrs.Y.Next(), MaxX: hrs.MaxX, MaxY: hrs.MaxY, ProjMethod: hrs.ProjMethod}},
		FalseBranch: DirectVectorAction{Target: container, Result: HorizontalRowStrategy{X: hrs.X.Next(), Y: hrs.Y, MaxX: hrs.MaxX, MaxY: hrs.MaxY, ProjMethod: hrs.ProjMethod}},
	}.Create().Select().Execute()
	return container.Value
}

type DirectVectorAction struct {
	Target *VectorContainer
	Result Vector
}
func (dva DirectVectorAction) IdentifyClass(consumer ClassConsumer) {}
func (dva DirectVectorAction) Execute()                             { dva.Target.Value = dva.Result }

func (hrs HorizontalRowStrategy) IsCanvasFinished() Bool   { return Zero{CompareTarget: hrs.Y}.CheckEquality() }
func (hrs HorizontalRowStrategy) IsGridIntersection() Bool { return hrs.X.IsMultipleOfGrid() }

type CubeIntersectionAcceptor struct {
	ScannerCoords HorizontalRowStrategy
	ResultTarget  *BoolContainer
}
func (cia CubeIntersectionAcceptor) AcceptProjection(vector Vector2D) { cia.ResultTarget.Value = vector.U.CheckEquality() }

type BoolContainer struct{ Value Bool }

func (hrs HorizontalRowStrategy) IsIntersecting3D() Bool {
	container := &BoolContainer{Value: False{}}
	cubeVertex := Point3D{X: hrs.X, Y: hrs.Y, Z: hrs.X}
	var activeProjector ProjectionStrategy
	if hrs.ProjMethod.Class() == "TopViewProjection" {
		activeProjector = TopViewProjection{Vertex: cubeVertex, Continuation: CubeIntersectionAcceptor{ScannerCoords: hrs, ResultTarget: container}}
	} else {
		activeProjector = SideViewProjection{Vertex: cubeVertex, Continuation: CubeIntersectionAcceptor{ScannerCoords: hrs, ResultTarget: container}}
	}
	activeProjector.Project()
	return container.Value
}

type ScanAction struct{ Scanner CanvasScanner }
func (sa ScanAction) IdentifyClass(consumer ClassConsumer) {}
func (sa ScanAction) Execute()                             { sa.Scanner.Canvas.ReadColor() }

type StopAction struct{ FinalSnapshot Snapshot[GameColor] }
func (sa StopAction) IdentifyClass(consumer ClassConsumer) {}
func (sa StopAction) Execute()                             {}

type Snapshot[T Object] interface {
	Object
	Accumulate() Snapshot[T]
}

type EmptySnapshot[T Object] struct{ NewPoint Point[T] }
func (es EmptySnapshot[T]) IdentifyClass(consumer ClassConsumer) {}
func (es EmptySnapshot[T]) Accumulate() Snapshot[T]             { return NodeSnapshot[T]{head: es.NewPoint, tail: es} }

type NodeSnapshot[T Object] struct {
	head     Point[T]
	tail     Snapshot[T]
	NewPoint Point[T]
}
func (ns NodeSnapshot[T]) IdentifyClass(consumer ClassConsumer) {}
func (ns NodeSnapshot[T]) Accumulate() Snapshot[T]             { return NodeSnapshot[T]{head: ns.NewPoint, tail: ns} }

type Point[T Object] interface{ Object }
type SnapshotPoint[T Object] struct {
	VectorState Vector
	Color       T
}
func (sp SnapshotPoint[T]) IdentifyClass(consumer ClassConsumer) {}

type Point3D struct {
	X Number
	Y Number
	Z Number
}

type ProjectionStrategy interface {
	Object
	Project() Action
	NextOrientation() ProjectionStrategy
}

type Vector2D struct {
	U Number
	V Number
}

type ProjectionAcceptor interface {
	AcceptProjection() // Сигнатура пустая, вектор Vector2D будет лежать внутри структуры акцептора
}

type TopViewProjection struct {
	Vertex       Point3D
	Continuation ProjectionAcceptor
}

func (tvp TopViewProjection) IdentifyClass() {}
func (tvp TopViewProjection) Project() Action {
	// Создаем акцептор и инкапсулируем вычисленные данные внутрь него, параметры отсутствуют
	tvp.Continuation.AcceptProjection() 
	return EmptyAction{}
}
func (tvp TopViewProjection) NextOrientation() ProjectionStrategy { return SideViewProjection{} }

type SideViewProjection struct {
	Vertex       Point3D
	Continuation ProjectionAcceptor
}

func (svp SideViewProjection) IdentifyClass() {}
func (svp SideViewProjection) Project() Action {
	svp.Continuation.AcceptProjection()
	return EmptyAction{}
}
func (svp SideViewProjection) NextOrientation() ProjectionStrategy { return TopViewProjection{} }

type GameColor interface{ Object }
type SolidWhiteColor struct{}
func (swc SolidWhiteColor) IdentifyClass() {}
type GridLineColor struct{}
func (glc GridLineColor) IdentifyClass() {}
type Object3DColor struct{}
func (o3c Object3DColor) IdentifyClass() {}
type TransparentColor struct{}
func (tc TransparentColor) IdentifyClass() {}

type ColorAcceptor interface {
	AcceptColor() // Сигнатура пустая
}

type SceneLayer interface {
	Object
	RenderPixel() Action
}

type WhiteBackgroundLayer struct{ Output ColorAcceptor }
func (wbl WhiteBackgroundLayer) IdentifyClass() {}
func (wbl WhiteBackgroundLayer) RenderPixel() Action {
	wbl.Output.AcceptColor()
	return EmptyAction{}
}

type CoordinateGridLayer struct {
	CurrentStep Vector
	Output      ColorAcceptor
}
func (cgl CoordinateGridLayer) IdentifyClass() {}
func (cgl CoordinateGridLayer) RenderPixel() Action {
	cgl.CurrentStep.IsGridIntersection().Select().Execute()
	return EmptyAction{}
}

type ThreeDimensionalObjectLayer struct {
	CurrentStep Vector
	Output      ColorAcceptor
}
func (tdol ThreeDimensionalObjectLayer) IdentifyClass() {}
func (tdol ThreeDimensionalObjectLayer) RenderPixel() Action {
	tdol.CurrentStep.IsIntersecting3D().Select().Execute()
	return EmptyAction{}
}

type Composited3DScene struct {
	Background  SceneLayer
	Grid        SceneLayer
	Object3D    SceneLayer
	FinalOutput ColorAcceptor
}
func (c3ds Composited3DScene) IdentifyClass() {}
func (c3ds Composited3DScene) RenderPixel() Action {
	c3ds.Object3D.RenderPixel()
	return EmptyAction{}
}

type LayerAction struct{ Layer SceneLayer }
func (la LayerAction) IdentifyClass() {}
func (la LayerAction) Execute()       { la.Layer.RenderPixel() }
