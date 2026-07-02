package main

// WavefrontIntersectionAcceptor используется для обработки проекции куба
type WavefrontIntersectionAcceptor struct {
	ResultTarget   UniversalContainer[Bool]
	ProjectedPoint Vector2D
}
func (wia WavefrontIntersectionAcceptor) AcceptProjection() { 
	wia.ResultTarget.Value = wia.ProjectedPoint.U.CheckEquality() 
}

func (wos WavefrontOrientedStrategy) IsIntersecting3D() Bool {
	// Длинная цепочка Next Пеано выстроена в линию, безопасную для go fmt
	cubeStart := Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()
	cubeEnd := cubeStart.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()

	container := UniversalContainer[Bool]{Value: False{}}
	
	// Вычисляем границы 3D куба Пеано полностью через полиморфную дифференциальную цепочку
	container.Value = BranchFactory{
		Condition: wos.X.Differentiate(cubeStart, Zero{}).CompareWithZero(),
		TrueBranch: DirectAction[Bool]{Target: UniversalContainer[Bool]{}, Result: cubeEnd.Differentiate(wos.X, Zero{}).CompareWithZero()},
		FalseBranch: DirectAction[Bool]{Target: UniversalContainer[Bool]{}, Result: False{}},
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
// ПОЛИМОРФНЫЙ СНИМОК КАДРА (Чистые пустые сигнатуры методов)
// ============================================================================

type Snapshot[T Object] interface {
	Object
	Accumulate()
	MutateSnapshotState()
}

type EmptySnapshot[T Object] struct {
	NewPoint       Point[T]
	AcceptorTarget UniversalContainer[Snapshot[T]]
	TargetNewPoint Point[T] 
}
func (es EmptySnapshot[T]) IdentifyClass() {}
func (es EmptySnapshot[T]) Accumulate() {
	es.AcceptorTarget.Value = NodeSnapshot[T]{head: es.NewPoint, tail: es}
}
func (es EmptySnapshot[T]) MutateSnapshotState() {
	es.NewPoint = es.TargetNewPoint
}

type NodeSnapshot[T Object] struct {
	head           Point[T]
	tail           Snapshot[T]
	NewPoint       Point[T]
	AcceptorTarget UniversalContainer[Snapshot[T]]
	TargetNewPoint Point[T]
}
func (ns NodeSnapshot[T]) IdentifyClass() {}
func (ns NodeSnapshot[T]) Accumulate() {
	ns.AcceptorTarget.Value = NodeSnapshot[T]{head: ns.NewPoint, tail: ns}
}
func (ns NodeSnapshot[T]) MutateSnapshotState() {
	ns.NewPoint = ns.TargetNewPoint
}

type Point[T Object] interface{ Object }
type SnapshotPoint[T Object] struct {
	VectorState Vector
	Color       T
}
func (sp SnapshotPoint[T]) IdentifyClass() {}
