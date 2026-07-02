package main

type Point3D struct {
	X Number
	Y Number
	Z Number
}

type ProjectionStrategy interface {
	Object
	Project() Action
	NextOrientation() ProjectionStrategy
}

type Vector2D struct {
	U Number
	V Number
}

type ProjectionAcceptor interface {
	AcceptProjection(vector Vector2D)
}

type TopViewProjection struct {
	Vertex       Point3D
	Continuation ProjectionAcceptor
}

func (tvp TopViewProjection) IdentifyClass(consumer ClassConsumer) {}
func (tvp TopViewProjection) Project() Action {
	tvp.Continuation.AcceptProjection(Vector2D{U: tvp.Vertex.X, V: tvp.Vertex.Z})
	return EmptyAction{}
}
func (tvp TopViewProjection) NextOrientation() ProjectionStrategy { return SideViewProjection{} }

type SideViewProjection struct {
	Vertex       Point3D
	Continuation ProjectionAcceptor
}

func (svp SideViewProjection) IdentifyClass(consumer ClassConsumer) {}
func (svp SideViewProjection) Project() Action {
	svp.Continuation.AcceptProjection(Vector2D{U: svp.Vertex.Y, V: svp.Vertex.Z})
	return EmptyAction{}
}
func (svp SideViewProjection) NextOrientation() ProjectionStrategy { return TopViewProjection{} }

type GameColor interface{ Object }
type SolidWhiteColor struct{}
func (swc SolidWhiteColor) IdentifyClass(consumer ClassConsumer) {}
type GridLineColor struct{}
func (glc GridLineColor) IdentifyClass(consumer ClassConsumer) {}
type Object3DColor struct{}
func (o3c Object3DColor) IdentifyClass(consumer ClassConsumer) {}
type TransparentColor struct{}
func (tc TransparentColor) IdentifyClass(consumer ClassConsumer) {}

type ColorAcceptor interface {
	AcceptColor(color GameColor)
}

type SceneLayer interface {
	Object
	RenderPixel() Action
}

type WhiteBackgroundLayer struct{ Output ColorAcceptor }
func (wbl WhiteBackgroundLayer) IdentifyClass(consumer ClassConsumer) {}
func (wbl WhiteBackgroundLayer) RenderPixel() Action {
	wbl.Output.AcceptColor(SolidWhiteColor{})
	return EmptyAction{}
}

type CoordinateGridLayer struct {
	CurrentStep Vector
	Output      ColorAcceptor
}
func (cgl CoordinateGridLayer) IdentifyClass(consumer ClassConsumer) {}
func (cgl CoordinateGridLayer) RenderPixel() Action {
	cgl.CurrentStep.IsGridIntersection().Select().Execute()
	return EmptyAction{}
}

type ThreeDimensionalObjectLayer struct {
	CurrentStep Vector
	Output      ColorAcceptor
}
func (tdol ThreeDimensionalObjectLayer) IdentifyClass(consumer ClassConsumer) {}
func (tdol ThreeDimensionalObjectLayer) RenderPixel() Action {
	tdol.CurrentStep.IsIntersecting3D().Select().Execute()
	return EmptyAction{}
}

type Composited3DScene struct {
	Background SceneLayer
	Grid       SceneLayer
	Object3D   SceneLayer
	FinalOutput ColorAcceptor
}
func (c3ds Composited3DScene) IdentifyClass(consumer ClassConsumer) {}
func (c3ds Composited3DScene) RenderPixel() Action {
	c3ds.Object3D.RenderPixel()
	return EmptyAction{}
}

type LayerAction struct{ Layer SceneLayer }
func (la LayerAction) IdentifyClass(consumer ClassConsumer) {}
func (la LayerAction) Execute()                             { la.Layer.RenderPixel() }
