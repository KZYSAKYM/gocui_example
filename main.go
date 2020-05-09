package main

import (
    "log"
    "github.com/jroimartin/gocui"
    "./env"
    "./textbox"
)

type Keybindings struct {
    /* if viewname is empty string, apply keybinding to all views */
    viewname    string
    /* key must be rune or gocui.Key */
    key         interface{}
    mod         gocui.Modifier
    handler     func(*gocui.Gui, *gocui.View) error
}

func main() {
    gui, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Panicln(err)
    }
    defer gui.Close()

    gui.SetManagerFunc(layout)

    keybindings := []Keybindings {
        {"", gocui.KeyCtrlC, gocui.ModNone, quit},
    }

    updateGlobalVar(gui)

    for _, v := range(keybindings) {
        err := gui.SetKeybinding(v.viewname, v.key, v.mod, v.handler)
        if err != nil {
            log.Panicln(err)
        }
        debug("Set Keybing: %s %v %v %v\n", v.viewname, v.key, v.mod, v.handler)
    }

    if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Panicln(err)
    }
}

func layout(gui *gocui.Gui) error {
    updateGlobalVar(gui)

    defaultViews := []*textbox.TextBox {
        textbox.New(gui, "hello0",
            textbox.WithText("hello0"),
            textbox.WithHAlign(textbox.TextHAlignCenter),
            textbox.WithVAlign(textbox.TextVAlignCenter),
            textbox.WithMarginX(3),
            textbox.WithMarginY(3),
            textbox.WithPos(0, 0),
        ),
        textbox.New(gui, "hello1",
            textbox.WithText("hello1"),
            textbox.WithHAlign(textbox.TextHAlignLeft),
            textbox.WithVAlign(textbox.TextVAlignTop),
            textbox.WithMarginX(2),
            textbox.WithMarginY(2),
            textbox.WithPos(env.CenterY, env.CenterX),
        ),
        textbox.New(gui, "hello2",
            textbox.WithText("hello2"),
            textbox.WithHAlign(textbox.TextHAlignRight),
            textbox.WithVAlign(textbox.TextVAlignBottom),
            textbox.WithMarginX(4),
            textbox.WithMarginY(4),
            textbox.WithPos(env.MaxY, env.MaxX),
        ),
    }

    for _, v := range(defaultViews) {
        if v == nil {
            continue
        }
        err := v.SetView()
        if err != nil {
            return err
        }
        debug("Set View: %s\n", v.Name())
    }
    return nil
}

func updateGlobalVar(gui *gocui.Gui) {
    env.MaxX, env.MaxY = gui.Size()
    env.CenterX = env.MaxX/2
    env.CenterY = env.MaxY/2
    debug("MAX: %dx%d", env.MaxX, env.MaxY)
    debug("Center: %dx%d", env.CenterX, env.CenterY)
}

func quit(gui *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func debug(format string, args ...interface{}) {
    if env.DEBUG {
        log.Printf(format, args...)
    }
}
