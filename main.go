package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func test_midi() {
	dir := os.TempDir()
	f := filepath.Join(dir, "smf-test.mid")

	defer os.Remove(f)

	var p printer

	err := writer.WriteSMF(f, 2, func(wr *writer.SMF) error {

		wr.SetChannel(11) // sets the channel for the next messages
		writer.NoteOn(wr, 120, 50)
		wr.SetDelta(120)
		writer.NoteOff(wr, 120)

		wr.SetDelta(240)
		writer.NoteOn(wr, 125, 50)
		wr.SetDelta(20)
		writer.NoteOff(wr, 125)
		writer.EndOfTrack(wr)

		wr.SetChannel(2)
		writer.NoteOn(wr, 120, 50)
		wr.SetDelta(60)
		writer.NoteOff(wr, 120)
		writer.EndOfTrack(wr)
		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", f)
		return
	}

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(reader.NoLogger(),
		// set the functions for the messages you are interested in
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)

	err = reader.ReadSMFFile(rd, f)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", f)
	}

	// Output: Track: 0 Pos: 0 NoteOn (ch 11: key 120 vel: 50)
	// Track: 0 Pos: 120 NoteOff (ch 11: key 120)
	// Track: 0 Pos: 360 NoteOn (ch 11: key 125 vel: 50)
	// Track: 0 Pos: 380 NoteOff (ch 11: key 125)
	// Track: 1 Pos: 0 NoteOn (ch 2: key 120 vel: 50)
	// Track: 1 Pos: 60 NoteOff (ch 2: key 120)
}

func convertMidi(midiFilePath string) {
	var p printer

	rd := reader.New(reader.NoteOn(p.noteOn),
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
		convertMidi(*midiFileInPtr)
	}

	fmt.Printf("kill youself:", *actionPtr)
	fmt.Printf("also kill youself:", *midiFileInPtr)

}
