package main

import "golang.org/x/mobile/gl"

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

// Каждая структура цвета теперь хранит низкоуровневый контекст и координаты для отрисовки точки
type SolidWhiteColor struct {
	GL gl.Context
	U  int
	V  int
}
func (swc SolidWhiteColor) IdentifyClass() {}
func (swc SolidWhiteColor) PaintHardwarePixel() {
	// Белый фон кадра (уже залит gl.Clear, метод остается чистым Null-Object-ом)
}

type GridLineColor struct {
	GL gl.Context
	U  int
	V  int
}
func (glc GridLineColor) IdentifyClass() {}
func (glc GridLineColor) PaintHardwarePixel() {
	// ФИЗИЧЕСКАЯ ОТРИСОВКА СЕТКИ: Ограничиваем область пикселя и заливаем серым цветом
	glc.GL.Scissor(int32(glc.U), int32(glc.V), 1, 1)
	glc.GL.ClearColor(0.8, 0.8, 0.8, 1.0)
	glc.GL.Clear(gl.COLOR_BUFFER_BIT)
}

type Object3DColor struct {
	GL gl.Context
	U  int
	V  int
}
func (o3c Object3DColor) IdentifyClass() {}
func (o3c Object3DColor) PaintHardwarePixel() {
	// ФИЗИЧЕСКАЯ ОТРИСОВКА КУБА: Заливаем точку контура куба ярко-красным цветом
	o3c.GL.Scissor(int32(o3c.U), int32(o3c.V), 2, 2)
	o3c.GL.ClearColor(1.0, 0.0, 0.0, 1.0)
	o3c.GL.Clear(gl.COLOR_BUFFER_BIT)
}

type TransparentColor struct{}
func (tc TransparentColor) IdentifyClass() {}
func (tc TransparentColor) PaintHardwarePixel() {}

type ColorAcceptor interface{ AcceptColor() }

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
func (c AppLifecycleLoop) Update()        {}
func (c Canvase) Scene()         {}
func (c Canvase) CreateScene()   {}
func (c Canvase) Frame()         {}
func (c Canvase) FrameScene()    {}
