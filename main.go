package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/reader"
)

type noteReader struct{}

type note struct {
	track    int16
	position uint64
	channel  uint8
	key      uint8
	vel      uint8
}

type score struct {
	Items []note
}

var curScore = score{}

// AddItem adds item to current score array
func (nt *score) AddItem(item note) {
	nt.Items = append(nt.Items, item)
}

func (pr noteReader) noteOn(p *reader.Position, channel, key, vel uint8) {
	//fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
	nt := note{track: p.Track}
	nt.position = p.AbsoluteTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	curScore.AddItem(nt)
}

func (pr noteReader) noteOff(p *reader.Position, channel, key, vel uint8) {
	nt := note{track: p.Track}
	nt.position = p.AbsoluteTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	curScore.AddItem(nt)

	//fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func readMidiFile(midiFilePath string) {
	//var n note
	var p noteReader
	rd := reader.New(reader.NoLogger(),
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)
	err := reader.ReadSMFFile(rd, midiFilePath)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", midiFilePath)
	}
}

func main() {
	actionPtr := flag.String("action", "nil", "What operation to perform. Choose from 'convert', 'shift' ")
	midiFileInPtr := flag.String("file_in", "nil", "Midi file to operate on")

	flag.Parse()
	if *actionPtr == "convert" {
		readMidiFile(*midiFileInPtr)
	}
	fmt.Println(len(curScore.Items))

}
