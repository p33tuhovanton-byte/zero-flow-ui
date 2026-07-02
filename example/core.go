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

// ИСПРАВЛЕНО: Полное избавление от TypeResolver и метода Class()
func (bf BranchFactory) Create() Bool {
	// Объект Bool сам выбирает и возвращает нужный Action, заряженный фабрикой
	bf.Condition.Select().Execute()
	return True{TrueBranch: bf.TrueBranch, FalseBranch: bf.FalseBranch}
}
