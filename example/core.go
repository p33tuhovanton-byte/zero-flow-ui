package main

// Node — корневой интерфейс для абсолютно всех системных комков движка
type Node interface {
	Object
	ProcessNode() Action
}

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

// UniversalContainer передается строго по значению (Immutable state)
type UniversalContainer[T Object] struct {
	Value T
}

// DirectAction — обобщенная CPS-команда, работающая со значениями без гонок памяти
type DirectAction[T Object] struct {
	Target UniversalContainer[T]
	Result T
}

func (da DirectAction[T]) IdentifyClass() {}
func (da DirectAction[T]) Execute()       {}

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
	AccumulateHardwareCoordinate()
}

type Zero struct {
	CompareTarget Number
	ActiveDriver  HardwareIntegerDriver 
}

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
func (z Zero) AccumulateHardwareCoordinate() {
	z.ActiveDriver.ExecuteHardwarePulse()
}

type Successor struct {
	pred          Number
	CompareTarget Number
	ActiveDriver  HardwareIntegerDriver
}

func (s Successor) IdentifyClass()       {}
func (s Successor) Class() string        { return "Successor" }
func (s Successor) Next() Number         { return Successor{pred: s} }
func (s Successor) CheckEquality() Bool  { return s.CompareTarget.CompareWithSuccessor() }
func (s Successor) CompareWithZero() Bool { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (s Successor) CompareWithSuccessor() Bool { return Successor{pred: s.pred, CompareTarget: s.CombinePredecessors()}.CheckEquality() }
func (s Successor) CombinePredecessors() Number { return s.CompareTarget.(Successor).pred }

func (s Successor) IsMultipleOfGrid() Bool {
	return s.EvaluateGridStep(Zero{}.Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next().
		Next().Next().Next().Next().Next().Next().Next().Next().Next().Next())
}

func (s Successor) EvaluateGridStep(currentStep Number) Bool {
	container := UniversalContainer[Bool]{}
	BranchFactory{
		Condition: currentStep.CompareWithZero(),
		TrueBranch: DirectAction[Bool]{Target: container, Result: s.IsMultipleOfGrid()},
		FalseBranch: DirectAction[Bool]{Target: container, Result: BranchFactory{
			Condition: s.pred.CompareWithZero(),
			TrueBranch: DirectAction[Bool]{Target: UniversalContainer[Bool]{}, Result: False{}},
			FalseBranch: DirectAction[Bool]{Target: UniversalContainer[Bool]{}, Result: s.pred.EvaluateGridStep(currentStep.(Successor).pred)},
		}.Create()},
	}.Create()
	return container.Value
}

func (s Successor) Differentiate(other Number, accumulator Number) Number {
	container := UniversalContainer[Number]{}
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
	container := UniversalContainer[Bool]{}
	BranchFactory{
		Condition: maxX.Differentiate(s, Zero{}).Differentiate(threshold, Zero{}).CompareWithZero(),
		TrueBranch: DirectAction[Bool]{Target: container, Result: True{}},
		FalseBranch: DirectAction[Bool]{Target: container, Result: False{}},
	}.Create()
	return container.Value
}

func (s Successor) AccumulateHardwareCoordinate() {
	nextDriver := s.ActiveDriver.IncrementPulse()
	nextState := s.pred
	if node, ok := nextState.(Successor); ok {
		node.ActiveDriver = nextDriver
		nextState = node
	} else if empty, ok := nextState.(Zero); ok {
		empty.ActiveDriver = nextDriver
		nextState = empty
	}
	nextState.AccumulateHardwareCoordinate()
}

type HardwareIntegerDriver interface {
	Object
	IncrementPulse() HardwareIntegerDriver
	ExecuteHardwarePulse()
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
	bf.Condition.Select().Execute()
	return bf.Condition
}
