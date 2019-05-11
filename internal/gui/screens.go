package gui

// Screen represents the the different screens of Tetris AI that the user sees.
type Screen int

// The possible screens are ScreenWelcome (intial) and ScreenPlay.
const (
	ScreenWelcome Screen = iota
	ScreenPlay
)
