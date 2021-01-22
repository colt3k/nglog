/*
//http://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#8-colors
	//Global Format
		\u001b[+COLOR+m		Non BRIGHT
		\u001b[+COLOR+;1m	BRIGHT

	// FORMATTING
		\u001b[1m		BOLD
		\u001b[4m		Underline
		\u001b[7m		Reversed

	// Formats Foreground
		8 bit		\u001b[30m
		16 bit		\u001b[30;1m
		256 bit 	\u001b[38;5;${ID}m

	//Colors
		CLR_BLK       = "\x1b[30m"   // black
		CLR_BLKBRIGHT = "\x1b[30;1m" // black bright
		CLR_RED       = "\x1b[31m"   // red
		CLR_REDBRT    = "\x1b[31;1m" // red
		CLR_GRN       = "\x1b[32m"   // green
		CLR_YLLW      = "\x1b[33m"   // yellow
		CLR_BLU       = "\x1b[34;1m" // blue
		CLR_MAG       = "\x1b[35;1m" // magenta
		CLR_CYAN      = "\x1b[36;1m" // cyan
		CLR_WHT   = "\x1b[37;1m" // white

	CLR_RESET = "\x1b[0m" // reset to default

	//ASCII 256
		CLR_DEFAULT = "\x1b[38;5;10m"
		CLR_ERR     = "\x1b[38;5;196m" //red
		CLR_WARN    = "\x1b[38;5;11m"  //yellow
		CLR_INFO    = "\x1b[38;5;174m" // light red/brown

	//Formats Background Color
	bright versions of the background colors do not change the background,
		but rather make the foreground text brighter
	8 bit 	\u001b[40m
	16 bit 	\u001b[40;1m
	256 bit	\u001b[48;5;
	CLRBG_BLK        = "\x1b[40m"
	CLRBG_BLKBRIGHTb = "\x1b[40;1m"
*/
package ng

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	escapeSequence = "\x1b"
)

var (
	// if NOT a terminal then turn off COLOR
	DumbTerm = os.Getenv("TERM") == "dumb"
	Term = isatty.IsTerminal(os.Stdout.Fd())
	CygwinTerm = isatty.IsCygwinTerminal(os.Stdout.Fd())
	NotTerminal = os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	colorsCache   = make(map[ColorAttr]*Clr)
	colorsCacheMu sync.Mutex // protects colorsCache
)


type ColorAttr int

// Base Color Attributes
const (
	Reset	ColorAttr = iota
	BOLD
	FAINT		// does nothing
	ITALIC
	UNDERLINE
	BlinkSlow	// does nothing
	BlinkRapid	// does nothing
	Reversed
	NonDisplayed
)


// FG Text Colors
const (
	FgBlack ColorAttr = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite

	CLRRESET = "\u001b[0m" // reset to default
	ClearTerminalSequence = "\033[2J"
)

// FG Hi-Intensity Text
const (
	HiBlack ColorAttr = iota + 90
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

// BG text colors
const (
	BgBlack ColorAttr = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack ColorAttr = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

var logcolor = map[int]string{
	30:"Black", 31:"Red", 32:"Green", 33:"Yellow", 34:"Blue",35:"Magenta",36:"Cyan",37:"White",
}

func (l ColorAttr) String() string {
	val := logcolor[int(l)]
	return val
}
func (l ColorAttr) Val() int {
	return int(l-1)
}
func (l ColorAttr) ValAsString() string {
	return strconv.Itoa(int(l-1))
}

/*
Types pulls full list as []string
*/
func (l ColorAttr) Types() map[int]string {
	return logcolor
}

type Clr struct {
	params []ColorAttr
	noColor *bool
}

func New(value ...ColorAttr) *Clr {
	t := new(Clr)
	t.params = make([]ColorAttr, 0)
	t.Add(value...)

	return t
}

func (c *Clr) Add(value ...ColorAttr) *Clr {
	c.params = append(c.params, value...)
	return c
}
func (c *Clr) format() string {
	str := fmt.Sprintf("%s[%sm", escapeSequence, c.sequence())
	return str
}
func (c *Clr) unformat() string {
	return fmt.Sprintf("%s[%dm", escapeSequence, Reset)
}

func CachedColor(p ColorAttr) *Clr {
	colorsCacheMu.Lock()
	defer colorsCacheMu.Unlock()

	c, ok := colorsCache[p]
	if !ok {
		c = New(p)
		colorsCache[p] = c
	}

	return c
}

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Clr) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

func (c *Clr) DisableColor() {
	c.noColor = boolPtr(true)
}
func boolPtr(v bool) *bool {
	return &v
}

func (c *Clr) Print(format string) string {
	s := fmt.Sprint(c.Set(), format, c.unset())
	return s
}
func (c *Clr) Printf(format string, a ...interface{}) string {

	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return fmt.Sprint(c.Set(), format, c.unset())
}
func Bool2Str(val bool) string {
	return strconv.FormatBool(val)
}
func (c *Clr) Set() string {
	//os.Stdout.Write([]byte("set not terminal, dumb? "+Bool2Str(DumbTerm)+" terminal? "+
	//	Bool2Str(Term)+" cygwinterminal? "+Bool2Str(CygwinTerm)+"\n"))
	if NotTerminal {
		return ""
	}

	return c.format()
}
func (c *Clr) unset() string {
	//os.Stdout.Write([]byte("unset not terminal, dumb? "+Bool2Str(DumbTerm)+" terminal? "+
	//	Bool2Str(Term)+" cygwinterminal? "+Bool2Str(CygwinTerm)+"\n"))
	if NotTerminal {
		return ""
	}
	return fmt.Sprintf( "%s[%dm", escapeSequence, Reset)
}


func Black(format string, a ...interface{}) string {
	return CachedColor(FgBlack).Print(format)
}
func BrightBlack(format string, a ...interface{}) string {
	return CachedColor(HiBlack).Print(format)
}
func Red(format string, a ...interface{}) string {
	return CachedColor(FgRed).Printf(format, a...)
}
func Green(format string, a ...interface{}) string {
	return CachedColor(FgGreen).Print(format)
}
func Yellow(format string, a ...interface{}) string {
	return CachedColor(FgYellow).Print(format)
}
func Blue(format string, a ...interface{}) string {
	return CachedColor(FgBlue).Print(format)
}
func Magenta(format string, a ...interface{}) string {
	return CachedColor(FgMagenta).Print(format)
}
func Cyan(format string, a ...interface{}) string {
	return CachedColor(FgCyan).Print(format)
}
func White(format string, a ...interface{}) string {
	return CachedColor(FgWhite).Print(format)
}
