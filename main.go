package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type noteReader struct{}

type noteOn struct {
	position reader.Position
	channel  uint8
	key      uint8
	vel      uint8
	noteOn   bool
}

type noteOff struct {
	position reader.Position
	channel  uint8
	key      uint8
	vel      uint8
	noteOn   bool
}

type Notes struct {
	notesOn  []noteOn
	notesOff []noteOff
}

type score struct {
	Items []noteOn
}

var curScore = Notes{}

// AddNoteOn adds item to current score array
func (nt *Notes) AddNoteOn(item noteOn) {
	nt.notesOn = append(nt.notesOn, item)
}

func (nt *Notes) AddNoteOff(item noteOff) {
	nt.notesOff = append(nt.notesOff, item)
}

func (pr noteReader) noteOn(p *reader.Position, channel, key, vel uint8) {
	//fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
	nt := noteOn{}
	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	nt.noteOn = true
	curScore.AddNoteOn(nt)
}

func (pr noteReader) noteOff(p *reader.Position, channel, key, vel uint8) {
	nt := noteOff{}
	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	nt.noteOn = true
	curScore.AddNoteOff(nt)

	//fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func readMidiFile(midiFilePath string) {
	//var n note
	var p noteReader
	rd := reader.New(reader.NoLogger(),
		//reader.EndOfTrack(),
		//reader.TempoBPM(),

		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)
	err := reader.ReadSMFFile(rd, midiFilePath)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", midiFilePath)
	}
}

func copyMidiFile(midiFilePath string, midiFileOut string) {
	readMidiFile(midiFilePath)
	err := writer.WriteSMF(midiFileOut, 1, func(wr *writer.SMF) error {
		for _, note := range curScore.notesOn {
			wr.SetChannel(note.channel)
			writer.NoteOn(wr, note.key, note.vel)
			//wr.SetDelta()
			//wr.Position()
			//writer.EndOfTrack(wr)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", midiFileOut)
	}
}

func main() {
	actionPtr := flag.String("action", "nil", "What operation to perform. Choose from 'convert', 'read', 'copy' ")
	midiFileInPtr := flag.String("file_in", "nil", "Midi file to operate on")
	//midiFileOutPtr := flag.String("file_out", "nil", "path and name of midi file to save")

	flag.Parse()
	if *actionPtr == "read" {
		readMidiFile(*midiFileInPtr)
	}
	//if *actionPtr == "copy" {
	//	copyMidiFile(*midiFileInPtr, *midiFileOutPtr)
	//}
	fmt.Println(len(curScore.notesOn))

}
