package main

type Object interface {
	IdentifyClass()
}

type ClassConsumer interface {
	AcceptClassName()
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
	EvaluateGridStep(currentStep Number) Bool
	// ДИФФЕРЕНЦИАЛЬНЫЙ КОНТРАКТ: Вычисляет абсолютную разность (производную) между числами
	Differentiate(other Number, accumulator Number) Number
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

func (z Zero) Differentiate(other Number, accumulator Number) Number {
	// Разность нуля и любого числа равна самому числу (накопленному в рекурсии)
	return accumulator
}

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
	return BranchFactory{
		Condition: currentStep.CompareWithZero(),
		TrueBranch: DirectBoolAction{Result: s.IsMultipleOfGrid()},
		FalseBranch: DirectBoolAction{Result: BranchFactory{
			Condition: s.pred.CompareWithZero(),
			TrueBranch: DirectBoolAction{Result: False{}},
			FalseBranch: DirectBoolAction{Result: s.pred.EvaluateGridStep(currentStep.(Successor).pred)},
		}.Create()},
	}.Create()
}

func (s Successor) Differentiate(other Number, accumulator Number) Number {
	// Рекурсивно уменьшаем оба числа Пеано, накапливая разность в CPS-стиле
	return BranchFactory{
		Condition: other.CompareWithZero(),
		TrueBranch: DirectNumberAction{Result: s},
		FalseBranch: DirectNumberAction{Result: s.pred.Differentiate(other.(Successor).pred, accumulator.Next())},
	}.Create().Select().(NumberAction).ResultNum
}

type DirectNumberAction struct{ Result Number }
func (dna DirectNumberAction) IdentifyClass() {}
func (dna DirectNumberAction) Execute()       {}

type NumberAction interface {
	Action
	GetNumber() Number
}

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
	container := &BoolResultContainer{}
	TypeResolver{ClassName: bf.Condition.Class(), T: bf.TrueBranch, F: bf.FalseBranch, Target: container}.Resolve()
	bf.Condition.Select().Execute()
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
