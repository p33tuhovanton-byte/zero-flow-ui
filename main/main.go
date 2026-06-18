package main

import (
        "zeroflowui"
)

func main() {
        textSignal := zeroflowui.TextSignal{
                Type:    zeroflowui.TextType,
                Payload: "Сигнальное сообщение успешно передано через функциональный поток.\n",
        }

        uiTimeline := zeroflowui.EndOfUI()
        uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "MainWindow", "Rendered")
        uiTimeline = zeroflowui.LogUIEvent(uiTimeline, true,  zeroflowui.EventLifecycle, "RefreshButton", "Initialized")
        uiTimeline = zeroflowui.LogUIEvent(uiTimeline, true,  zeroflowui.EventInteraction, "RefreshButton", "Clicked")
        uiTimeline = zeroflowui.LogUIEvent(uiTimeline, false, zeroflowui.EventLifecycle, "MainWindow", "Crash")
        uiTimeline = zeroflowui.LogUIEvent(uiTimeline, true,  zeroflowui.EventInteraction, "RetryButton", "Clicked")

        pipeline := zeroflowui.SystemPipelineDecorator{
                Next: zeroflowui.TerminalProcessor{},
        }

        pipeline.Process(zeroflowui.NewTextFlow(textSignal), uiTimeline)
}