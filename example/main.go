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
// НАЧАЛО ИНТЕГРИРОВАННОГО ЯДРА БИБЛИОТЕКИ ZEROFLOWUI
// ============================================================================

type UIStateDescriptor struct {
	EventType rune
	Message   rune
}

const (
	EventLifecycle   rune = 'L'
	EventInteraction rune = 'I'
)

type ZeroFlowEngine struct{}

func (ZeroFlowEngine) DispatchSignal(a app.App, rawEvent interface{}, node MobileEventChain, runner *ApplicationRunner, atlas StructuralAtlas) {
	if ev, ok := rawEvent.(lifecycle.Event); ok {
		ActiveTypeMatcher{}.MatchLifecycle(ev, node, a, runner, atlas)
	}
	if ev, ok := rawEvent.(size.Event); ok {
		ActiveTypeMatcher{}.MatchSize(ev, node, a, runner, atlas)
	}
	if ev, ok := rawEvent.(touch.Event); ok {
		ActiveTypeMatcher{}.MatchTouch(ev, node, a, runner, atlas)
	}
	if ev, ok := rawEvent.(paint.Event); ok {
		ActiveTypeMatcher{}.MatchPaint(ev, node, a, runner, atlas)
	}
}

type ActiveTypeMatcher struct{}

func (ActiveTypeMatcher) MatchLifecycle(ev lifecycle.Event, node MobileEventChain, a app.App, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.ProcessLifecycle(a, ev, runner, atlas)
}
func (ActiveTypeMatcher) MatchSize(ev size.Event, node MobileEventChain, a app.App, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.ProcessSize(a, ev, runner, atlas)
}
func (ActiveTypeMatcher) MatchTouch(ev touch.Event, node MobileEventChain, a app.App, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.ProcessTouch(a, ev, runner, atlas)
}
func (ActiveTypeMatcher) MatchPaint(ev paint.Event, node MobileEventChain, a app.App, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.ProcessPaint(a, ev, runner, atlas)
}

// ============================================================================
// СЛОЙ ГРАФИЧЕСКИХ КОНТЕКСТОВ И АДАПТЕРОВ (OPENGL СВЯЗЬ)
// ============================================================================

type GraphicContext struct {
	GL gl.Context
}

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
// ПАТТЕРН "СОСТОЯНИЕ" (STATE) ДЛЯ ИСКЛЮЧЕНИЯ IF/ELSE ПРИ ОТРИСОВКЕ
// ============================================================================

type RenderState interface {
	RenderGlyphs(glCtx gl.Context, chain GlyphDecorator, ctx UIContext)
}

type RightAnchoredButtonState struct{}

func (RightAnchoredButtonState) RenderGlyphs(glCtx gl.Context, chain GlyphDecorator, ctx UIContext) {
	// ИСПРАВЛЕНО: Увеличен шаг по оси X между 'W' и 'O' (с 10 до 20 пикселей), чтобы символы не сливались
	chain.RenderGlyph(glCtx, 'W', ctx.EdgeX-40, ctx.CurrentY, 1, 0, 0, 0)
	chain.RenderGlyph(glCtx, 'O', ctx.EdgeX-20, ctx.CurrentY, 1, 0, 0, 0)
}

type InteractionState struct{}

func (InteractionState) RenderGlyphs(glCtx gl.Context, chain GlyphDecorator, ctx UIContext) {
	// ИСПРАВЛЕНО: Шаг каретки по горизонтали увеличен до 16 пикселей для читаемости строки 'In'
	chain.RenderGlyph(glCtx, 'I', ctx.EdgeX-40, ctx.CurrentY, 1, 0, 255, 0) // Зеленый маркер
	chain.RenderGlyph(glCtx, 'n', ctx.EdgeX-20, ctx.CurrentY, 1, 0, 255, 0)
}

type DefaultState struct{}

func (DefaultState) RenderGlyphs(glCtx gl.Context, chain GlyphDecorator, ctx UIContext) {
	// ИСПРАВЛЕНО: Расстояние между 'L' и 'y' скорректировано под общую сетку
	chain.RenderGlyph(glCtx, 'L', ctx.EdgeX-40, ctx.CurrentY, 1, 0, 0, 0)
	chain.RenderGlyph(glCtx, 'y', ctx.EdgeX-20, ctx.CurrentY, 1, 0, 0, 0)
}

// ============================================================================
// КОНВЕЙЕР РЕНДЕРИНГА ЭКРАНА С СИСТЕМОЙ ШАГОВ (БЕЗ РЕКУРСИИ С RETURN)
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
	row.CurrentRowState.RenderGlyphs(glCtx, atlas.Chain, ctx)
	// ИСПРАВЛЕНО: Шаг между строками увеличен с 12 до 24 пикселей.
	// Это разнесет строки по вертикали и уберет наложение пикселей.
	row.NextRow.RenderNextRow(glCtx, atlas, UIContext{
		EdgeX:            ctx.EdgeX,
		CurrentY:         ctx.CurrentY - 24,
		ScreenHeightByte: ctx.ScreenHeightByte,
	})
}

// ============================================================================
// ПАТТЕРН "ЦЕПОЧКА ОБЯЗАННОСТЕЙ" ДЛЯ ОБРАБОТКИ СИСТЕМНЫХ ПРЕРЫВАНИЙ
// ============================================================================

type MobileEventChain interface {
	ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner, atlas StructuralAtlas)
	ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner, atlas StructuralAtlas)
	ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner, atlas StructuralAtlas)
	ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner, atlas StructuralAtlas)
}

type BaseEventChainNode struct {
	Next MobileEventChain
}

func (node BaseEventChainNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessLifecycle(a, ev, runner, atlas)
}
func (node BaseEventChainNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessSize(a, ev, runner, atlas)
}
func (node BaseEventChainNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessTouch(a, ev, runner, atlas)
}
func (node BaseEventChainNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessPaint(a, ev, runner, atlas)
}

type LifecycleNode struct {
	BaseEventChainNode
}

func (node LifecycleNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	if glCtx, ok := ev.DrawContext.(gl.Context); ok {
		runner.GL = glCtx
		a.Send(paint.Event{})
	}
	node.Next.ProcessLifecycle(a, ev, runner, atlas)
}

type SizeNode struct {
	BaseEventChainNode
}

func (node SizeNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessSize(a, ev, runner, atlas)
}

type TouchNode struct {
	BaseEventChainNode
}

func (node TouchNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	node.Next.ProcessTouch(a, ev, runner, atlas)
}

type PaintNode struct {
	BaseEventChainNode
}

func (node PaintNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner, atlas StructuralAtlas) {
	if runner.GL != nil {
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
		}.RenderNextRow(runner.GL, atlas, runner.InitialContext)
		
		a.Publish()
	}
	node.Next.ProcessPaint(a, ev, runner, atlas)
}

type TerminalEventNode struct{}

func (TerminalEventNode) ProcessLifecycle(a app.App, ev lifecycle.Event, runner *ApplicationRunner, atlas StructuralAtlas) {}
func (TerminalEventNode) ProcessSize(a app.App, ev size.Event, runner *ApplicationRunner, atlas StructuralAtlas)           {}
func (TerminalEventNode) ProcessTouch(a app.App, ev touch.Event, runner *ApplicationRunner, atlas StructuralAtlas)          {}
func (TerminalEventNode) ProcessPaint(a app.App, ev paint.Event, runner *ApplicationRunner, atlas StructuralAtlas)          {}

// ============================================================================
// ИСТИННЫЙ МАТРИЧНЫЙ РЕНДЕРЕР
// ============================================================================

type GlyphDecorator interface {
	RenderGlyph(glCtx gl.Context, charCode rune, x, y, scale, r, g, b rune)
}

type ActiveHardwareGlyphRenderer struct{}

func (ActiveHardwareGlyphRenderer) RenderGlyph(glCtx gl.Context, charCode rune, x, y, scale, r, g, b rune) {
	glCtx.Enable(gl.SCISSOR_TEST)
	// ИСПРАВЛЕНО: Увеличена базовая пиксельная сетка вырезания (Scissor) до 32x32, 
	// чтобы маркеры букв стали крупнее, четче и превратились в полноценные элементы UI.
	glCtx.Scissor(int32(x)*4, int32(y)*4, int32(scale)*32, int32(scale)*32)
	glCtx.ClearColor(float32(r)/255.0, float32(g)/255.0, float32(b)/255.0, 1.0)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
	glCtx.Disable(gl.SCISSOR_TEST)
}

type StructuralAtlas struct {
	Chain GlyphDecorator
}

type UIElementContainer interface {
	DispatchTouch(pipe MobileEventChain, tx, ty rune)
}

// ============================================================================
// ТОЧКА СБОРКИ И ЗАПУСКА: APPLICATIONRUNNER (ОБЪЕКТНЫЙ IOC КОНТЕЙНЕР)
// ============================================================================

type ApplicationRunner struct {
	Atlas          StructuralAtlas
	InitialContext UIContext
	EventPipeline  MobileEventChain
	Engine         ZeroFlowEngine
	GL             gl.Context
}

func (runner ApplicationRunner) Start(a app.App) {
	for e := range a.Events() {
		runner.Engine.DispatchSignal(a, e, runner.EventPipeline, &runner, runner.Atlas)
	}
}

func main()         {app.Main(ApplicationRunner{
 Atlas: StructuralAtlas{
  Chain: ActiveHardwareGlyphRenderer{},},
  InitialContext: UIContext{
    EdgeX: 180,
    CurrentY: 160, 
    ScreenHeightByte: 240,},
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
