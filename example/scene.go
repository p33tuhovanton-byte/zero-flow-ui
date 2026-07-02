package main

type Point3D struct{ X, Y, Z Number }

type ProjectionStrategy interface {
	Object
	Project() Action
	NextOrientation() ProjectionStrategy
	InjectContinuation()
}

type Vector2D struct{ U, V Number }
type ProjectionAcceptor interface{ AcceptProjection() }

type TopViewProjection struct {
	Vertex       Point3D
	Continuation ProjectionAcceptor
}

func (tvp TopViewProjection) IdentifyClass() {}
func (tvp TopViewProjection) Project() Action {
	tvp.Continuation.AcceptProjection()
	return EmptyAction{}
}
func (tvp TopViewProjection) NextOrientation() ProjectionStrategy { return SideViewProjection{} }
func (tvp TopViewProjection) InjectContinuation()                 {}

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
func (svp SideViewProjection) InjectContinuation()                 {}

type GameColor interface {
	Object
	PaintHardwarePixel()
}

type SolidWhiteColor struct{}
func (swc SolidWhiteColor) IdentifyClass() {}
func (swc SolidWhiteColor) PaintHardwarePixel() {}

type GridLineColor struct {
	DriverX HardwareIntegerDriver
	DriverY HardwareIntegerDriver
}
func (glc GridLineColor) IdentifyClass() {}
func (glc GridLineColor) PaintHardwarePixel() {
	glc.DriverX.ExecuteHardwarePulse()
}

type Object3DColor struct {
	DriverX HardwareIntegerDriver
	DriverY HardwareIntegerDriver
}
func (o3c Object3DColor) IdentifyClass() {}
func (o3c Object3DColor) PaintHardwarePixel() {
	o3c.DriverX.ExecuteHardwarePulse()
}

type TransparentColor struct{}
func (tc TransparentColor) IdentifyClass() {}
func (tc TransparentColor) PaintHardwarePixel() {}

type ColorAcceptor interface {
	Action
	AcceptColor()
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

// ============================================================================
// ИЕРАРХИЯ КОМКОВ ДВИЖКА (Node Tree)
// ============================================================================

type CameraNode struct {
	Projection ProjectionStrategy
	ChildNode  Node
}
func (cn CameraNode) IdentifyClass() {}
func (cn CameraNode) ProcessNode() Action {
	cn.ChildNode.ProcessNode().Execute()
	return EmptyAction{}
}

type SceneNode struct {
	Background  SceneLayer
	Grid        SceneLayer
	Object3D    SceneLayer
	FinalOutput ColorAcceptor
}
func (sn SceneNode) IdentifyClass() {}
func (sn SceneNode) ProcessNode() Action {
	sn.Object3D.RenderPixel()
	return EmptyAction{}
}
func (sn SceneNode) RenderPixel() Action {
	sn.Object3D.RenderPixel()
	return EmptyAction{}
}

type CharacterTouchControllerNode struct {
	CharacterPositionX Number
	CharacterPositionY Number
	NextAsset          Node
}
func (ctcn CharacterTouchControllerNode) IdentifyClass() {}
func (ctcn CharacterTouchControllerNode) ProcessNode() Action {
	ctcn.NextAsset.ProcessNode().Execute()
	return EmptyAction{}
}

type Composited3DScene struct {
	Background, Grid, Object3D SceneLayer
	FinalOutput                ColorAcceptor
}
func (c3ds Composited3DScene) IdentifyClass() {}
func (c3ds Composited3DScene) RenderPixel() Action {
	c3ds.Object3D.RenderPixel()
	return EmptyAction{}
}

type LayerAction struct{ Layer SceneLayer }
func (la LayerAction) IdentifyClass() {}
func (la LayerAction) Execute()       { la.Layer.RenderPixel() }

type Render interface {
	Object
	Update()
	Scene()
	CreateScene()
	Frame()
	FrameScene()
}

type AndroidFrame struct {
	ActiveCanvas Canvase
	NextFrame    Action
}
func (fa AndroidFrame) IdentifyClass() {}
func (fa AndroidFrame) Update()        { fa.Scene() }
func (fa AndroidFrame) Scene()         { fa.CreateScene() }
func (fa AndroidFrame) CreateScene()   { fa.Frame() }
func (fa AndroidFrame) Frame() {
	fa.ActiveCanvas.ScanTarget.Scan()
	fa.FrameScene()
}
func (fa AndroidFrame) FrameScene() { fa.NextFrame.Execute() }

type Canvase struct{ ScanTarget CanvasScanner }
func (c Canvase) IdentifyClass() {}
func (c Canvase) Update()        {}
func (c Canvase) Scene()         {}
func (c Canvase) CreateScene()   {}
func (c Canvase) Frame()         {}
func (c Canvase) FrameScene()    {}
