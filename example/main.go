package main

import "golang.org/x/mobile/gl"

// GameModLauncher — Обобщенный Generic-мод запуска
type GameModLauncher struct {
	GL         gl.Context
	Width      Number
	Height     Number
	Projection ProjectionStrategy
	CurrentVec Vector
	OutVec     *UniversalContainer[Vector]
}

func (gml GameModLauncher) LaunchMod() {
	CameraNode{
		Projection: gml.Projection,
		ChildNode: CanvasNode{
			ScanStrategy: gml.CurrentVec,
			HardwareGL:   gml.GL,
			OutVec:       gml.OutVec,
		},
	}.ProcessNode().Execute()
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

type CanvasNode struct {
	ScanStrategy Vector
	HardwareGL   gl.Context
	OutVec       *UniversalContainer[Vector]
}
func (cn CanvasNode) IdentifyClass() {}
func (cn CanvasNode) ProcessNode() Action {
	// ИСПРАВЛЕНО: Типизация Viewport приведена к строгому int, съедаемому пакетом gl
	wVal := peanoToInt(cn.ScanStrategy.(WavefrontOrientedStrategy).MaxX, 0)
	hVal := peanoToInt(cn.ScanStrategy.(WavefrontOrientedStrategy).MaxY, 0)
	cn.HardwareGL.Viewport(0, 0, wVal, hVal)

	CanvasScanner{Step: cn.ScanStrategy, Canvas: OpenGlCanvas{GlContext: cn.HardwareGL}, Storage: EmptySnapshot[GameColor]{}, OutVec: cn.OutVec}.Scan()
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
	OutVec  *UniversalContainer[Vector]
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
	container := &UniversalContainer[Snapshot[GameColor]]{}
	psa.Scanner.Storage.MutateSnapshotState(
		SnapshotPoint[GameColor]{VectorState: psa.Scanner.Step, Color: psa.InjectedColor},
		container,
	)

	CanvasScanner{
		Step:    psa.Scanner.Step.AdvanceVector(),
		Canvas:  psa.UpdatedCanvas,
		Storage: container.Value,
		OutVec:  psa.Scanner.OutVec,
	}.Scan()
}

type GenericColorAcceptor[T GameColor] struct {
	Target *UniversalContainer[GameColor]
	Result T
}

func (gca GenericColorAcceptor[T]) IdentifyClass() {}
func (gca GenericColorAcceptor[T]) Execute()       { gca.Target.Value = gca.Result }
func (gca GenericColorAcceptor[T]) AcceptColor()   { gca.Execute() }

type DirectColorAction struct {
	Target *PixelSaveAcceptor
	Color  GameColor
}
func (dca DirectColorAction) IdentifyClass() {}
func (dca DirectColorAction) Execute()       { dca.Target.InjectedColor = dca.Color; dca.Target.AcceptColor() }

func (cs CanvasScanner) Scan() {
	saveAcceptor := PixelSaveAcceptor{Scanner: cs, UpdatedCanvas: OpenGlCanvas{GlContext: cs.Canvas.(OpenGlCanvas).GlContext}}
	glCtx := cs.Canvas.(OpenGlCanvas).GlContext

	// ИСПРАВЛЕНО: Ликвидирована ошибка declared and not used для uVal и vVal кадра.
	// Координаты Пеано напрямую инжектируются в полиморфный поток вывода
	scene := SceneNode{
		Background: WhiteBackgroundLayer{Output: GenericColorAcceptor[GameColor]{Target: &UniversalContainer[GameColor]{}, Result: SolidWhiteColor{}}},
		Grid: CoordinateGridLayer{CurrentStep: cs.Step, Output: GenericColorAcceptor[GameColor]{Target: &UniversalContainer[GameColor]{}, Result: GridLineColor{
			DriverX: OpenGlPixelDriver{GL: glCtx, Counter: peanoToInt(cs.Step.(WavefrontOrientedStrategy).X, 0), IsYAxis: false},
			DriverY: OpenGlPixelDriver{GL: glCtx, Counter: peanoToInt(cs.Step.(WavefrontOrientedStrategy).Y, 0), IsYAxis: true},
		}}},
		Object3D: ThreeDimensionalObjectLayer{CurrentStep: cs.Step, Output: GenericColorAcceptor[GameColor]{Target: &UniversalContainer[GameColor]{}, Result: Object3DColor{
			DriverX: OpenGlPixelDriver{GL: glCtx, Counter: peanoToInt(cs.Step.(WavefrontOrientedStrategy).X, 0), IsYAxis: false},
			DriverY: OpenGlPixelDriver{GL: glCtx, Counter: peanoToInt(cs.Step.(WavefrontOrientedStrategy).Y, 0), IsYAxis: true},
		}}},
		FinalOutput: saveAcceptor,
	}
	saveAcceptor.UpdatedCanvas.Scene = scene
	
	BranchFactory{Condition: cs.Step.IsCanvasFinished(), TrueBranch: DirectAction[Vector]{Target: cs.OutVec, Result: cs.Step}, FalseBranch: ScanAction{Scanner: cs}}.Create().Select().Execute()
}

func peanoToInt(num Number, acc int) int {
	if num.Class() == "Zero" {
		return acc
	}
	return peanoToInt(num.(Successor).pred, acc+1)
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
	container := &UniversalContainer[Vector]{}
	BranchFactory{
		Condition: wos.X.EvaluateWaveCenter(wos.MaxX, Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()),
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

type WavefrontIntersectionAcceptor struct {
	ResultTarget   *UniversalContainer[Bool]
	ProjectedPoint Vector2D
}
func (wia WavefrontIntersectionAcceptor) AcceptProjection() { wia.ResultTarget.Value = wia.ProjectedPoint.U.CheckEquality() }

func (wos WavefrontOrientedStrategy) IsIntersecting3D() Bool {
	container := &UniversalContainer[Bool]{Value: False{}}
	
	cubeStart := Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()
		
	cubeEnd := cubeStart.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()

	isAfterStart := wos.X.Differentiate(cubeStart, Zero{}).CompareWithZero()
	isBeforeEnd := cubeEnd.Differentiate(wos.X, Zero{}).CompareWithZero()
	
	container.Value = BranchFactory{
		Condition: isAfterStart,
		TrueBranch: DirectAction[Bool]{Target: &UniversalContainer[Bool]{}, Result: isBeforeEnd},
		FalseBranch: DirectAction[Bool]{Target: &UniversalContainer[Bool]{}, Result: False{}},
	}.Create()
	
	wos.ProjMethod.InjectContinuation()
	wos.ProjMethod.Project()
	return container.Value
}

type ScanAction struct{ Scanner CanvasScanner }
func (sa ScanAction) IdentifyClass() {}
func (sa ScanAction) Execute()       { sa.Scanner.Canvas.ReadColor() }

type StopAction struct{ FinalSnapshot Snapshot[GameColor] }
func (sa StopAction) IdentifyClass() {}
func (sa StopAction) Execute()       {}

// ============================================================================
// ПОЛИМОРФНЫЙ СНИМОК КАДРА (Идеальные пустые сигнатуры методов)
// ============================================================================

type Snapshot[T Object] interface {
	Object
	Accumulate()
	MutateSnapshotState(newPoint Point[T], container *UniversalContainer[Snapshot[T]])
}

type EmptySnapshot[T Object] struct {
	NewPoint       Point[T]
	AcceptorTarget *UniversalContainer[Snapshot[T]]
}
func (es EmptySnapshot[T]) IdentifyClass() {}
func (es EmptySnapshot[T]) Accumulate() {
	es.AcceptorTarget.Value = NodeSnapshot[T]{head: es.NewPoint, tail: es}
}
func (es EmptySnapshot[T]) MutateSnapshotState(newPoint Point[T], container *UniversalContainer[Snapshot[T]]) {
	es.NewPoint = newPoint
	es.AcceptorTarget = container
}

type NodeSnapshot[T Object] struct {
	head           Point[T]
	tail           Snapshot[T]
	NewPoint       Point[T]
	AcceptorTarget *UniversalContainer[Snapshot[T]]
}
func (ns NodeSnapshot[T]) IdentifyClass() {}
func (ns NodeSnapshot[T]) Accumulate() {
	ns.AcceptorTarget.Value = NodeSnapshot[T]{head: ns.NewPoint, tail: ns}
}
func (ns NodeSnapshot[T]) MutateSnapshotState(newPoint Point[T], container *UniversalContainer[Snapshot[T]]) {
	ns.NewPoint = newPoint
	ns.AcceptorTarget = container
}

type Point[T Object] interface{ Object }
type SnapshotPoint[T Object] struct {
	VectorState Vector
	Color       T
}
func (sp SnapshotPoint[T]) IdentifyClass() {}
