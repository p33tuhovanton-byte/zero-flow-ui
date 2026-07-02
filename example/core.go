package main

// ============================================================================
// ОБОБЩЕННЫЕ КОРНЕВЫЕ КОНТРАКТЫ (Generics без any)
// ============================================================================

type Object interface {
	IdentifyClass()
}

type Node interface {
	Object
	ProcessNode() Action
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

// UniversalContainer инкапсулирует результат вычисления любого типа T, расширяющего Object
type UniversalContainer[T Object] struct {
	Value T
}

// DirectAction — обобщенная CPS-команда, избавляющая от десятков мелких структур
type DirectAction[T Object] struct {
	Target *UniversalContainer[T]
	Result T
}

func (da DirectAction[T]) IdentifyClass() {}
func (da DirectAction[T]) Execute()       { da.Target.Value = da.Result }

// ============================================================================
// АРИФМЕТИКА ПЕАНО НА УНИВЕРСАЛЬНЫХ ШАБЛОНАХ
// ============================================================================

type Number interface {
	Object
	Next() Number
	CheckEquality() Bool
	CompareWithZero() Bool
	CompareWithSuccessor() Bool
	IsMultipleOfGrid() Bool
	EvaluateGridStep(currentStep Number) Bool
	Differentiate(other Number, accumulator Number) Number
	EvaluateWaveCenter(maxX Number, threshold Number) Bool
}

type Zero struct{ CompareTarget Number }

func (z Zero) IdentifyClass()       {}
func (z Zero) Class() string        { return "Zero" }
func (z Zero) Next() Number         { return Successor{pred: z} }
func (z Zero) CheckEquality() Bool  { return z.CompareTarget.CompareWithZero() }
func (z Zero) CompareWithZero() Bool { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) CompareWithSuccessor() Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) IsMultipleOfGrid() Bool { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) EvaluateGridStep(currentStep Number) Bool { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) Differentiate(other Number, accumulator Number) Number { return accumulator }
func (z Zero) EvaluateWaveCenter(maxX Number, threshold Number) Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }

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

func (s Successor) IsMultipleOfGrid() Bool {
	gridInterval := Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next()
	return s.EvaluateGridStep(gridInterval)
}

func (s Successor) EvaluateGridStep(currentStep Number) Bool {
	container := &UniversalContainer[Bool]{}
	BranchFactory{
		Condition: currentStep.CompareWithZero(),
		TrueBranch: DirectAction[Bool]{Target: container, Result: s.IsMultipleOfGrid()},
		FalseBranch: DirectAction[Bool]{Target: container, Result: BranchFactory{
			Condition: s.pred.CompareWithZero(),
			TrueBranch: DirectAction[Bool]{Target: &UniversalContainer[Bool]{}, Result: False{}},
			FalseBranch: DirectAction[Bool]{Target: &UniversalContainer[Bool]{}, Result: s.pred.EvaluateGridStep(currentStep.(Successor).pred)},
		}.Create()},
	}.Create()
	return container.Value
}

func (s Successor) Differentiate(other Number, accumulator Number) Number {
	container := &UniversalContainer[Number]{}
	BranchFactory{
		Condition: other.CompareWithZero(),
		TrueBranch: DirectAction[Number]{Target: container, Result: s},
		FalseBranch: DirectAction[Number]{Target: container, Result: s.pred.Differentiate(other.(SimpleSuccessorResolver).GetPred(), accumulator.Next())},
	}.Create().Select().Execute()
	return container.Value
}

type SimpleSuccessorResolver interface{ GetPred() Number }
func (s Successor) GetPred() Number { return s.pred }

func (s Successor) EvaluateWaveCenter(maxX Number, threshold Number) Bool {
	distance := maxX.Differentiate(s, Zero{})
	container := &UniversalContainer[Bool]{}
	BranchFactory{
		Condition: distance.Differentiate(threshold, Zero{}).CompareWithZero(),
		TrueBranch: DirectAction[Bool]{Target: container, Result: True{}},
		FalseBranch: DirectAction[Bool]{Target: container, Result: False{}},
	}.Create()
	return container.Value
}

// ============================================================================
// ПОЛИМОРФНЫЕ ВЕТВЛЕНИЯ СИСТЕМЫ
// ============================================================================

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

func (bf BranchFactory) Create() Bool {
	container := &UniversalContainer[Bool]{}
	TypeResolver{ClassName: bf.Condition.Class(), T: bf.TrueBranch, F: bf.FalseBranch, Target: container}.Resolve()
	bf.Condition.Select().Execute()
	return bf.Condition
}

type TypeResolver struct {
	ClassName string
	T, F      Action
	Target    *UniversalContainer[Bool]
}

func (tr TypeResolver) Resolve() {
	tr.Target.Value = True{TrueBranch: tr.T, FalseBranch: tr.F}
}
