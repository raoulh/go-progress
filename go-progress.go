package progress

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/sys/unix"
	"modernc.org/libc/unistd"
)

type Format struct {
	Fill     []string
	Head     string
	Empty    string
	LeftEnd  string
	RightEnd string
}

var (
	ProgressFormats = []Format{

		// █▓░░░░░░░░░░ 9%
		{
			Fill:  []string{"▓", "█"},
			Empty: "░",
		},

		// ⬤◯◯◯◯◯◯◯◯◯ 9%
		{
			Fill:  []string{"⬤"},
			Empty: "◯",
		},

		// ■□□□□□□□□□□□ 9%
		{
			Fill:  []string{"■"},
			Empty: "□",
		},

		// ⚫⚫⚫⚫⚪⚪⚪⚪⚪⚪ 41%
		{
			Fill:  []string{"⚫"},
			Empty: "⚪",
		},

		// ▰▰▰▰▱▱▱▱▱▱ 41%
		{
			Fill:  []string{"▰"},
			Empty: "▱",
		},

		// ⬛⬛⬛⬛⬜⬜⬜⬜⬜⬜ 41%
		{
			Fill:  []string{"⬛"},
			Empty: "⬜",
		},

		// ⣿⣿⣿⣿⡟⣀⣀⣀⣀⣀⣀ 41%
		{
			Fill:  []string{"⡀", "⡄", "⡆", "⡇", "⡏", "⡟", "⡿", "⣿"},
			Empty: "⣀",
		},

		// [======>             ]
		{
			Fill:     []string{"="},
			Head:     ">",
			LeftEnd:  "[",
			RightEnd: "]",
			Empty:    " ",
		},

		// ▉▉▋            41%
		{
			Fill:  []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉"},
			Empty: " ",
		},
	}
)

type ProgressBar struct {
	Format

	Out   io.Writer //where to write
	Width int       //width of the bar

	ShowPercent    bool //show percent
	ShowNumeric    bool //show current/total
	ShowTextSuffix bool //show suffix text

	total      int
	progress   int
	lastWidth  int
	textSuffix string
}

func New(total int) *ProgressBar {
	return &ProgressBar{
		Format:         ProgressFormats[0],
		total:          total,
		Out:            os.Stdout,
		Width:          40,
		ShowNumeric:    true,
		ShowPercent:    true,
		ShowTextSuffix: false,
	}
}

func getTerminalWidth() int {
	ws, err := unix.IoctlGetWinsize(unistd.STDIN_FILENO, unix.TIOCGWINSZ)
	if err != nil {
		return 0
	}

	return int(ws.Col)
}

func (p *ProgressBar) SetTextSuffix(suffix string) {
	p.textSuffix = suffix
}

func (p *ProgressBar) Set(to int) bool {
	if to < 0 {
		return false
	} else if to > p.total {
		to = p.total
	}

	if to == p.total {
		p.clear()
		return false
	}

	p.progress = to
	p.paint()

	return true
}

func (p *ProgressBar) Inc() bool {
	return p.Set(p.progress + 1)
}

func (p *ProgressBar) clear() {
	s := "\r"

	w := getTerminalWidth()
	if w == 0 {
		w = p.lastWidth
	}

	for i := 0; i < w; i++ {
		s += " "
	}
	s += "\r"
	io.WriteString(p.Out, s)
}

func (p *ProgressBar) paint() {
	s := "\r"
	width := p.Width

	percent := (float64(p.progress) / float64(p.total)) * 100.0
	totalNumVal := p.Width * len(p.Fill) * int(percent) / 100

	s += p.LeftEnd

	for i := 0; i <= totalNumVal/len(p.Fill); i++ {
		fs := p.Fill[len(p.Fill)-1]
		width -= utf8.RuneCountInString(fs)
		s += fs
	}
	s += p.Head
	width -= utf8.RuneCountInString(p.Head)
	if totalNumVal%len(p.Fill) > 0 {
		fs := p.Fill[totalNumVal%len(p.Fill)]
		width -= utf8.RuneCountInString(fs)
		s += fs
	}

	for width > 0 {
		s += p.Empty
		width--
	}

	suffix := ""

	if p.ShowNumeric {
		s += fmt.Sprintf(" %d/%d", p.progress, p.total)
	}
	if p.ShowPercent {
		s += fmt.Sprintf(" [%d%%]", int(percent))
	}

	s += p.RightEnd
	s += suffix

	if p.ShowTextSuffix {
		s += " "
		s += p.textSuffix
	}

	for utf8.RuneCountInString(s) < p.lastWidth {
		s += " "
	}

	p.lastWidth = utf8.RuneCountInString(s)
	io.WriteString(p.Out, s)
}
