package main

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

// ============================================================================
// ЯДРО ДИСПЕТЧЕРИЗАЦИИ СИГНАЛОВ ZEROFLOWUI БЕЗ IF И OK
// ============================================================================

type ZeroFlowEngine struct{}

func (ZeroFlowEngine) DispatchSignal(a app.App, rawEvent interface{}, node MobileEventChain, runner *ApplicationRunner) {
	// Двойная диспетчеризация (Double Dispatch) через полиморфные интерфейсы x/mobile
	switch ev := rawEvent.(type) {
	case lifecycle.Event:
		ActiveTypeMatcher{}.MatchLifecycle(ev, node, a, runner)
	case size.Event:
		ActiveTypeMatcher{}.MatchSize(ev, node, a, runner)
	case touch.Event:
		ActiveTypeMatcher{}.MatchTouch(ev, node, a, runner)
	case paint.Event:
		ActiveTypeMatcher{}.MatchPaint(ev, node, a, runner)
	}
}

type ActiveTypeMatcher struct{}

func (ActiveTypeMatcher) MatchLifecycle(ev lifecycle.Event, node MobileEventChain, a app.App, runner *ApplicationRunner) {
	node.ProcessLifecycle(a, ev, runner)
}
func (ActiveTypeMatcher) MatchSize(ev size.Event, node MobileEventChain, a app.App, runner *ApplicationRunner) {
	node.ProcessSize(a, ev, runner)
}
func (ActiveTypeMatcher) MatchTouch(ev touch.Event, node MobileEventChain, a app.App, runner *ApplicationRunner) {
	node.ProcessTouch(a, ev, runner)
}
func (ActiveTypeMatcher) MatchPaint(ev paint.Event, node MobileEventChain, a app.App, runner *ApplicationRunner) {
	node.ProcessPaint(a, ev, runner)
}

// ============================================================================
// ГРАФИЧЕСКИЙ КОНТЕКСТ И АДАПТЕР ПОДЛОЖКИ
// ============================================================================

type UIContext struct {
	EdgeX            rune
	CurrentY         rune
	ScreenHeightByte rune
}

type OpenGLBackgroundAdapter struct{}

func (OpenGLBackgroundAdapter) ClearTargetScreen(glCtx gl.Context, colorValue rune) {
	glCtx.ClearColor(float32(colorValue)/255.0, float32(colorValue)/255.0, float32(colorValue)/255.0, 1.0)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}

// ============================================================================
// ПАТТЕРН "СОСТОЯНИЕ" (STATE) ДЛЯ СТРОК СТАТУСА UI
// ============================================================================

type RenderState interface {
	RenderGlyphs(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext)
}

type RightAnchoredButtonState struct{}

func (RightAnchoredButtonState) RenderGlyphs(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext) {
	atlas.W.DrawVector(glCtx, ctx.EdgeX-40, ctx.CurrentY, 0, 0, 0)
	atlas.O.DrawVector(glCtx, ctx.EdgeX-20, ctx.CurrentY, 0, 0, 0)
	atlas.K.DrawVector(glCtx, ctx.EdgeX-10, ctx.CurrentY, 0, 0, 0)
}

type InteractionState struct{}

func (InteractionState) RenderGlyphs(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext) {
	atlas.I.DrawVector(glCtx, ctx.EdgeX-40, ctx.CurrentY, 0, 255, 0)
	atlas.N.DrawVector(glCtx, ctx.EdgeX-20, ctx.CurrentY, 0, 255, 0)
}

type DefaultState struct{}

func (DefaultState) RenderGlyphs(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext) {
	atlas.L.DrawVector(glCtx, ctx.EdgeX-40, ctx.CurrentY, 0, 0, 0)
	atlas.Y.DrawVector(glCtx, ctx.EdgeX-20, ctx.CurrentY, 0, 0, 0)
}

// ============================================================================
// ПОТОКОВЫЙ ИТЕРАТОР СЕТКИ СТРОК UI
// ============================================================================

type ScreenStreamIterator interface {
	RenderNextRow(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext)
}

type EndOfScreenStream struct{}

func (EndOfScreenStream) RenderNextRow(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext) {
	glCtx.Flush()
}

type ActiveScreenRow struct {
	CurrentRowState RenderState
	NextRow         ScreenStreamIterator
}

func (row ActiveScreenRow) RenderNextRow(glCtx gl.Context, atlas StructuralAtlas, ctx UIContext) {
	row.CurrentRowState.RenderGlyphs(glCtx, atlas, ctx)
	row.NextRow.RenderNextRow(glCtx, atlas, UIContext{
		EdgeX:            ctx.EdgeX,
		CurrentY:         ctx.CurrentY - 24,
		ScreenHeightByte: ctx.ScreenHeightByte,
	})
}

// ============================================================================
// ЦЕПОЧКА ОБРАБОТКИ СОБЫТИЙ С АВТО-МАСШТАБИРОВАНИЕМ СПРАВА
// ============================================================================

type MobileEventChain interface {
	ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner)
	ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner)
	ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner)
	ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner)
}

type BaseEventChainNode struct {
	Next MobileEventChain
}

func (node BaseEventChainNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner) {
	node.Next.ProcessLifecycle(a, ev, runner)
}
func (node BaseEventChainNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner) {
	node.Next.ProcessSize(a, ev, runner)
}
func (node BaseEventChainNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner) {
	node.Next.ProcessTouch(a, ev, runner)
}
func (node BaseEventChainNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner) {
	node.Next.ProcessPaint(a, ev, runner)
}

type LifecycleNode struct {
	BaseEventChainNode
}

func (node LifecycleNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner) {
	switch glCtx := ev.DrawContext.(type) {
	case gl.Context:
		runner.GL = glCtx
		a.Send(paint.Event{})
	}
	node.Next.ProcessLifecycle(a, ev, runner)
}

type SizeNode struct {
	BaseEventChainNode
}

func (node SizeNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner) {
	runner.InitialContext.EdgeX = rune(ev.WidthPx / 4)
	a.Send(paint.Event{})
	node.Next.ProcessSize(a, ev, runner)
}

type TouchNode struct {
	BaseEventChainNode
}

func (node TouchNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner) {
	node.Next.ProcessTouch(a, ev, runner)
}

type PaintNode struct {
	BaseEventChainNode
}

func (node PaintNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner) {
	switch {
	case runner.GL != nil:
		OpenGLBackgroundAdapter{}.ClearTargetScreen(runner.GL, 255)
		
		ActiveScreenRow{
			CurrentRowState: RightAnchoredButtonState{},
			NextRow: ActiveScreenRow{
				CurrentRowState: InteractionState{},
				NextRow: ActiveScreenRow{
					CurrentRowState: DefaultState{},
					NextRow:         EndOfScreenStream{},
				},
			},
		}.RenderNextRow(runner.GL, runner.Atlas, runner.InitialContext)
		
		a.Publish()
	}
	node.Next.ProcessPaint(a, ev, runner)
}

type TerminalEventNode struct{}

func (TerminalEventNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner) {}
func (TerminalEventNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner)           {}
func (TerminalEventNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner)          {}
func (TerminalEventNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner)          {}

// ============================================================================
// ПАТТЕРН "ПОЛИМОРФНЫЙ ЗНАКОГЕНЕРАТОР" (АБСОЛЮТНО БЕЗ IF И SWITCH)
// ============================================================================

type GlyphVector interface {
	DrawVector(glCtx gl.Context, x, y, r, g, b rune)
}

type VectorW struct{}

func (VectorW) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Scissor(int32(x)*4, int32(y)*4, 4, 24)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Scissor(int32(x+2)*4, int32(y)*4, 4, 12)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Scissor(int32(x+4)*4, int32(y)*4, 4, 24)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorO struct{}

func (VectorO) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Scissor(int32(x)*4, int32(y)*4, 16, 24)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.ClearColor(255.0/255.0, 255.0/255.0, 255.0/255.0, 1.0) // Вырезаем внутреннюю часть белым
	glCtx.Scissor(int32(x+1)*4, int32(y+1)*4, 8, 16)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorK struct{}

func (VectorK) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Scissor(int32(x)*4, int32(y)*4, 4, 24)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Scissor(int32(x+2)*4, int32(y+2)*4, 8, 4)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorI struct{}

func (VectorI) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Scissor(int32(x+1)*4, int32(y)*4, 4, 24)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorN struct{}

func (VectorN) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Scissor(int32(x)*4, int32(y)*4, 4, 16)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Scissor(int32(x)*4, int32(y+3)*4, 12, 4)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Scissor(int32(x+2)*4, int32(y)*4, 4, 16)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorL struct{}

func (VectorL) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
 glCtx.Scissor(int32(x)*4, int32(y)*4, 4, 24)
	
 glCtx.Clear(gl.COLOR_BUFFER_BIT)
 glCtx.Scissor(int32(x)*4, int32(y)*4, 16, 4)   
 glCtx.Clear(gl.COLOR_BUFFER_BIT)
 glCtx.Disable(gl.SCISSOR_TEST)
}

type VectorY struct{}

func (VectorY) DrawVector(glCtx gl.Context, x, y, r, g, b rune) {
 glCtx.Enable(gl.SCISSOR_TEST)
 glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
 glCtx.Scissor(int32(x)*4, int32(y)*4, 12, 4)
 glCtx.Clear(gl.COLOR_BUFFER_BIT)
 glCtx.Scissor(int32(x+2)*4, int32(y)*4, 4, 20)
 glCtx.Clear(gl.COLOR_BUFFER_BIT)
 glCtx.Disable(gl.SCISSOR_TEST)
}

type StructuralAtlas struct {
 W, O, K, I, N, L, Y GlyphVector
}

type ApplicationRunner struct {
  Atlas   StructuralAtlasInitialContext
  UIContextEventPipeline MobileEventChainEngine
  ZeroFlowEngineGL  gl.Context
}

func (runner ApplicationRunner) Start(a app.App) {
 for e := range a.Events() {
  runner.Engine.DispatchSignal(a, e, runner.EventPipeline, &runner)
 }
}

func main()  {
app.Main(ApplicationRunner{
   Atlas: StructuralAtlas{
    W: VectorW{},
    O: VectorO{},
    K: VectorK{},
    I: VectorI{},
    N: VectorN{},
    L: VectorL{},
    Y: VectorY{},
   },
   InitialContext: UIContext{
    EdgeX:       180,
    CurrentY:    160,
    ScreenHeightByte: 240,
   },
   Engine: ZeroFlowEngine{},
   EventPipeline: LifecycleNode{
    BaseEventChainNode: BaseEventChainNode{
     Next: SizeNode{
      BaseEventChainNode: BaseEventChainNode{
       Next: TouchNode{
        BaseEventChainNode: BaseEventChainNode{
         Next: PaintNode{
          BaseEventChainNode: BaseEventChainNode{
           Next: TerminalEventNode{},
         },
        },
       },
      },
     },
    },
   },
  },
 }.Start)
}