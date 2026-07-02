package main

type WavefrontIntersectionAcceptor struct {
	ResultTarget   UniversalContainer[Bool]
	ProjectedPoint Vector2D
}
func (wia WavefrontIntersectionAcceptor) AcceptProjection() { wia.ResultTarget.Value = wia.ProjectedPoint.U.CheckEquality() }

func (wos WavefrontOrientedStrategy) IsIntersecting3D() Bool {
	cubeStart := Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()
	cubeEnd := cubeStart.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()

	uContainer := UniversalContainer[Number]{}
	wos.X.Differentiate()
	isAfterStart := uContainer.Value.CompareWithZero()
	
	finalContainer := UniversalContainer[Bool]{Value: False{}}
	BranchFactory{
		Condition: isAfterStart,
		TrueBranch: DirectAction[Bool]{Target: finalContainer, Result: isAfterStart},
		FalseBranch: DirectAction[Bool]{Target: finalContainer, Result: False{}},
	}.Create()
	
	wos.ProjMethod.InjectContinuation()
	wos.ProjMethod.Project()
	
	return finalContainer.Value
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
