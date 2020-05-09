package textbox

import (
    "fmt"
    "log"
    "github.com/jroimartin/gocui"
    "../env"
)

/*
 * Use Applicable Functional Option Pattern.
 * Reference:
 *   https://ww24.jp/2019/07/go-option-pattern
 *   https://play.golang.org/p/rFoYqiYicB9
 */

/*= TextBox ================================================================*/

type TextBox struct {
    gui     *gocui.Gui
    name    string
    text    string
    halign  TextHAlign
    valign  TextVAlign
    marginX int
    marginY int
    top     int
    left    int
    width   int
    height  int
    /* private */
    bottom  int
    right   int
    repr    string
}

type Option interface {
    Apply(*TextBox)
}

/*= TextHAlign =============================================================*/

/*-- Type --*/
type TextHAlign uint
const (
    TextHAlignCenter = iota + 1
    TextHAlignLeft
    TextHAlignRight
)
/*-- Util --*/
func validHAlign(v TextHAlign) error {
    switch v {
    case TextHAlignCenter:
        fallthrough
    case TextHAlignLeft:
        fallthrough
    case TextHAlignRight:
        return nil
    default:
        return fmt.Errorf("Found invalid value: %v", v)
    }
}
/*-- Stringer --*/
func (v TextHAlign) String() string {
    switch v {
    case TextHAlignCenter:
        return "TextHAlignCenter"
    case TextHAlignLeft:
        return "TextHAlignLeft"
    case TextHAlignRight:
        return "TextHAlignRight"
    default:
        return "Unknown"
    }
}
/*-- Optional Arguments --*/
func (v TextHAlign) Apply(t *TextBox) {
    t.SetHAlign(v)
}
func WithHAlign(v TextHAlign) TextHAlign {
    return TextHAlign(v)
}
/*-- Setter/Getter --*/
func (self *TextBox) SetHAlign(v TextHAlign) {
    if err := validHAlign(v); err != nil {
        log.Fatal(err)
    }
    self.halign = v
}
func (self *TextBox) HAlign() TextHAlign {
    return self.halign
}

/*= TextVAlign =============================================================*/

/*-- Type --*/
type TextVAlign uint
const (
    TextVAlignCenter = iota + 1
    TextVAlignTop
    TextVAlignBottom
)
/*-- Util --*/
func validVAlign(v TextVAlign) error {
    switch v {
    case TextVAlignCenter:
        fallthrough
    case TextVAlignTop:
        fallthrough
    case TextVAlignBottom:
        return nil
    default:
        return fmt.Errorf("Found invalid value: %v", v)
    }
}
/*-- Stringer --*/
func (v TextVAlign) String() string {
    switch v {
    case TextVAlignCenter:
        return "TextVAlignCenter"
    case TextVAlignTop:
        return "TextVAlignTop"
    case TextVAlignBottom:
        return "TextVAlignBottom"
    default:
        return "Unknown"
    }
}
/*-- Optional Arguments --*/
func (v TextVAlign) Apply(t *TextBox) {
    t.SetVAlign(v)
}
func WithVAlign(v TextVAlign) TextVAlign {
    return TextVAlign(v)
}
/*-- Setter/Getter --*/
func (self *TextBox) SetVAlign(v TextVAlign) {
    if err := validVAlign(v); err != nil {
        log.Fatal(err)
    }
    self.valign = v
}
func (self *TextBox) VAlign() TextVAlign {
    return self.valign
}

/*= Name ===================================================================*/
/*-- Setter/Getter --*/
func (self *TextBox) Name() string {
    return self.name
}

/*= Text ===================================================================*/

/*-- Type --*/
type Text string
/*-- Optional Arguments --*/
func (v Text) Apply(t *TextBox) {
    t.SetText(string(v))
}
func WithText(v string) Text {
    return Text(v)
}
/*-- Setter/Getter --*/
func (self *TextBox) SetText(text string) {
    self.text = text
}
func (self *TextBox) Text() string {
    return self.text
}
func (self *TextBox) Repr() string {
    return self.repr
}

/*= Margin =================================================================*/

/*-- Type --*/
type MarginX int
type MarginY int
/*-- Optional Arguments --*/
func (v MarginX) Apply(t *TextBox) {
    t.SetMarginX(int(v))
}
func WithMarginX(v int) MarginX {
    return MarginX(v)
}
func (v MarginY) Apply(t *TextBox) {
    t.SetMarginY(int(v))
}
func WithMarginY(v int) MarginY {
    return MarginY(v)
}
/*-- Setter/Getter --*/
func (self *TextBox) SetMarginX(marginX int) {
    self.marginX = marginX
}
func (self *TextBox) MarginX() int {
    return self.marginX
}
func (self *TextBox) SetMarginY(marginY int) {
    self.marginY = marginY
}
func (self *TextBox) MarginY() int {
    return self.marginY
}

/*= Position ================================================================*/

/*-- Type --*/
type Pos struct {
    top  int
    left int
}
/*-- Util --*/
func validTopLeft(t *TextBox, top, left int) (int, int) {
    var x, y int
    if top < 0 {
        y = 0
    } else if top > env.MaxY {
        y = env.MaxY - 2 - (t.marginY * 2) - (env.MarginY * 2)
    } else {
        y = top
    }
    if left < 0 {
        x = 0
    } else if left > env.MaxX {
        x = env.MaxX - len(t.text) - (t.marginX * 2) - (env.MarginX * 2)
    }
    return y, x
}
/*-- Optional Arguments --*/
func (v Pos) Apply(t *TextBox) {
    t.top, t.left = validTopLeft(t, v.top, v.left)
}
func WithPos(top, left int) Option {
    return Pos{top, left}
}
/*-- Setter/Getter --*/
func (self *TextBox) SetPos(top int, left int) {
    self.top, self.left = validTopLeft(self, top, left)
}
func (self *TextBox) Pos() (int, int, int, int) {
    return self.top, self.left, self.right, self.bottom
}

/*= Private Methods =========================================================*/

func (self *TextBox) init(name string, gui *gocui.Gui) {
    self.gui  = gui
    self.name = name
    self.text = ""
    self.halign = TextHAlignCenter
    self.valign = TextVAlignCenter
    self.top = 0
    self.left = 0
    self.marginX = 1
    self.marginY = 1
}

func (self *TextBox) align() {
    var text string

    self.width = len(self.text) + self.marginX*2 + env.MarginX*2
    /* character's height is 2 */
    self.height = 2 + self.marginY*2 + env.MarginY*2

    self.right = self.left + self.width
    self.bottom = self.top + self.height

    if self.right > env.MaxX {
        self.right = env.MaxX
    }
    if self.bottom > env.MaxY {
        self.bottom = env.MaxY
    }

    if self.left + self.width > env.MaxX {
        self.left = env.MaxX - self.width
        self.right = env.MaxX
    }
    if self.top + self.height > env.MaxY {
        self.top = env.MaxY - self.height - 1
        self.bottom = env.MaxY - 1
    }

    /* text margin top */
    var cond int
    switch self.valign {
    case TextVAlignCenter:
        cond = self.marginY
    case TextVAlignTop:
        cond = 0
    case TextVAlignBottom:
        cond = self.marginY*2
    default:
        panic("Found invalid value in self.valign")
    }
    for i := 0; i < cond; i++ {
        text += "\n"
    }
    /* text margin left */
    switch self.halign {
    case TextHAlignCenter:
        cond = self.marginX
    case TextHAlignLeft:
        cond = 0
    case TextHAlignRight:
        cond = self.marginX*2
    default:
        panic("Found invalid value in self.halign")
    }
    for i := 0; i < cond; i++ {
        text += " "
    }
    text += self.text
    self.repr = text
}

/*= Public Methods ==========================================================*/

func (self *TextBox) SetView() error {
    self.align()
    view, err := self.gui.SetView(
        self.name,
        self.left,
        self.top,
        self.right,
        self.bottom,
    )
    if err != gocui.ErrUnknownView {
        return err
    }
    fmt.Fprintln(view, self.Repr())
    return nil
}

/*= Public Functions ========================================================*/

func New(gui *gocui.Gui,
         name string,
         options ...Option) *TextBox {
    obj := &TextBox{}
    obj.init(name, gui)
    for _, o := range options {
        o.Apply(obj)
    }
    return obj
}
