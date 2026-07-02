package core

type Object interface {
	IdentifyClass(consumer ClassConsumer)
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

func (ea EmptyAction) IdentifyClass(consumer ClassConsumer) {}
func (ea EmptyAction) Execute()                             {}

type Number interface {
	Object
	Next() Number
	CheckEquality() Bool
	CompareWithZero() Bool
	CompareWithSuccessor() Bool
	IsMultipleOfGrid() Bool
}

type Zero struct{ CompareTarget Number }

func (z Zero) IdentifyClass(consumer ClassConsumer) {}
func (z Zero) Class() string                        { return "Zero" }
func (z Zero) Next() Number                         { return Successor{pred: z} }
func (z Zero) CheckEquality() Bool                  { return z.CompareTarget.CompareWithZero() }
func (z Zero) CompareWithZero() Bool                { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) CompareWithSuccessor() Bool           { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (z Zero) IsMultipleOfGrid() Bool               { return True{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }

type Successor struct {
	pred          Number
	CompareTarget Number
}

func (s Successor) IdentifyClass(consumer ClassConsumer) {}
func (s Successor) Class() string                        { return "Successor" }
func (s Successor) Next() Number                         { return Successor{pred: s} }
func (s Successor) CheckEquality() Bool                  { return s.CompareTarget.CompareWithSuccessor() }
func (s Successor) CompareWithZero() Bool                { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }
func (s Successor) CompareWithSuccessor() Bool           { return Successor{pred: s.pred, CompareTarget: s.CombinePredecessors()}.CheckEquality() }
func (s Successor) CombinePredecessors() Number          { return s.CompareTarget.(Successor).pred }
func (s Successor) IsMultipleOfGrid() Bool               { return False{TrueBranch: EmptyAction{}, FalseBranch: EmptyAction{}} }

type True struct{ TrueBranch, FalseBranch Action }

func (t True) IdentifyClass(consumer ClassConsumer) {}
func (t True) Class() string                        { return "True" }
func (t True) Select() Action                       { return t.TrueBranch }

type False struct{ TrueBranch, FalseBranch Action }

func (f False) IdentifyClass(consumer ClassConsumer) {}
func (f False) Class() string                        { return "False" }
func (f False) Select() Action                       { return f.FalseBranch }

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
