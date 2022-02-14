package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type fileReader struct{}

// Note individual noteon and noteoff messages
type Note struct {
	position reader.Position
	channel  uint8
	key      uint8
	vel      uint8
	noteOn   bool
}

// Notes a collection of notes
type Notes struct {
	msg []Note
}

type Track struct {
	name     string
	position reader.Position
	notes    Notes
}

type Tracks struct {
	track []Track
}

type score struct {
	Items []Note
}

var curScore = Notes{}
var gScore = Tracks{}

// AddNote adds item to current score array
func (nt *Notes) AddNote(item Note) {
	nt.msg = append(nt.msg, item)
}

func (pr fileReader) noteOn(p *reader.Position, channel, key, vel uint8) {
	//fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
	nt := Note{}
	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	nt.noteOn = true
	curScore.AddNote(nt)
}

func (pr fileReader) noteOff(p *reader.Position, channel, key, vel uint8) {
	nt := Note{}
	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	nt.noteOn = false
	curScore.AddNote(nt)
}

func (pr fileReader) instrument(p reader.Position, name string) {
	tk := Track{}
	tk.name = name
	tk.position = p
}

func readMidiFile(midiFilePath string) {
	//var n note
	var p fileReader
	rd := reader.New(reader.NoLogger(),
		//reader.TempoBPM(),
		reader.Instrument(p.instrument),
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
		//reader.EndOfTrack(),

	)
	err := reader.ReadSMFFile(rd, midiFilePath)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", midiFilePath)
	}
}

func copyMidiFile(midiFilePath string, midiFileOut string) {
	readMidiFile(midiFilePath)
	err := writer.WriteSMF(midiFileOut, 1, func(wr *writer.SMF) error {
		for _, note := range curScore.msg {
			wr.SetChannel(note.channel)
			wr.SetDelta(note.position.DeltaTicks)
			if note.noteOn {
				writer.NoteOn(wr, note.key, note.vel)

			}
			if !note.noteOn {
				writer.NoteOff(wr, note.key)
			}
			//wr.SetDelta()
			//wr.Position()
			//writer.EndOfTrack(wr)
		}
		writer.EndOfTrack(wr)
		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", midiFileOut)
	}
}

func main() {
	actionPtr := flag.String("action", "nil", "What operation to perform. Choose from 'convert', 'read', 'copy' ")
	midiFileInPtr := flag.String("file_in", "nil", "Midi file to operate on")
	midiFileOutPtr := flag.String("file_out", "nil", "path and name of midi file to save")

	flag.Parse()
	if *actionPtr == "read" {
		readMidiFile(*midiFileInPtr)
	}
	if *actionPtr == "copy" {
		copyMidiFile(*midiFileInPtr, *midiFileOutPtr)
	}
	fmt.Println(len(curScore.msg))

}
