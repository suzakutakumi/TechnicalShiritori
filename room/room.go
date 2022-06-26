package room

type Room struct {
	Number int
	Users  []string
	Words  []Word
}
type Word struct {
	Text string
	Lead string
	Tail string
}

const (
	WAIT = iota
	RUN
	FINISH
)

func (w *Word) SetWord(t string) {
	w.Text = t
	runes := w.Charcters()
	w.Lead = string([]rune{runes[0]})
	w.Lead = string([]rune{runes[len(runes)-1]})
}

func (w Word) Charcters() []rune {
	return []rune(w.Text)
}
