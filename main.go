package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/midimessage/meta"
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

type Tempo struct {
	position reader.Position
	bpm      float64
}

type Track struct {
	name   string
	number int16
	notes  Notes
}

type Tracks struct {
	track []Track
}

//var curScore = Notes{}
var gTracks = Tracks{}

// AddNote adds item to current score array
func (nt *Notes) AddNote(item Note) {
	nt.msg = append(nt.msg, item)
}

func (at *Tracks) AddTrack(item Track) {
	at.track = append(at.track, item)
}

func (nt *Track) AddNotes(item Notes) {
	nt.notes = item
}

func createTrackIfNotExists(trackNum int16) Track {
	existingTracks := gTracks.track
	for _, track := range existingTracks {
		//println("existing tracks are:")
		//println(track.number)
		if track.number == trackNum {
			//println("brak")
			return track
		}
	}
	tk := Track{}
	tk.number = trackNum
	gTracks.AddTrack(tk)
	//println("adding track!")
	//fmt.Println(trackNum)
	return tk
}

func getTrackIndex(trackNum int16) int {
	for index, track := range gTracks.track {
		if track.number == trackNum {
			return index
		}
	}
	return 1
}

func (pr fileReader) noteOn(p *reader.Position, channel, key, vel uint8) {
	//fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
	tk := createTrackIfNotExists(p.Track)
	nt := Note{}

	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	//tk.notes.AddNote()
	nt.noteOn = true
	tkIndex := getTrackIndex(tk.number)
	gTracks.track[tkIndex].notes.AddNote(nt)
	//tk.notes.AddNote(nt)
}

func (pr fileReader) noteOff(p *reader.Position, channel, key, vel uint8) {
	tk := createTrackIfNotExists(p.Track)
	nt := Note{}

	tk.number = p.Track
	nt.position.AbsoluteTicks = p.AbsoluteTicks
	nt.position.Track = p.Track
	nt.position.DeltaTicks = p.DeltaTicks
	nt.channel = channel
	nt.key = key
	nt.vel = vel
	//tk.notes.AddNote()
	nt.noteOn = false
	tkIndex := getTrackIndex(tk.number)
	gTracks.track[tkIndex].notes.AddNote(nt)

}

func (pr fileReader) instrument(p reader.Position, name string) {
	tk := Track{}
	tk.name = name
	tk.number = p.Track
	//fmt.Println(tk.name)
	gTracks.AddTrack(tk)
}

func (pr fileReader) tempobpm(p reader.Position, bpm float64) {
	tp := Tempo{}
	tp.bpm = bpm
	tp.position = p
}

func readMidiFile(midiFilePath string) Tracks {
	//var n note
	var p fileReader
	rd := reader.New(reader.NoLogger(),
		reader.TempoBPM(p.tempobpm),
		reader.Instrument(p.instrument),
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
		//reader.EndOfTrack(),
	)

	err := reader.ReadSMFFile(rd, midiFilePath)
	//fmt.Println(gScore.track)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", midiFilePath)
	}
	//for _, val := range gTracks.track {
	//	fmt.Println(val)
	//}
	return gTracks
}

func copyMidiFile(midiFilePath string, midiFileOut string) {
	readMidiFile(midiFilePath)
	// count tracks
	numTracks := uint16(len(gTracks.track))

	err := writer.WriteSMF(midiFileOut, numTracks, func(wr *writer.SMF) error {
		for _, track := range gTracks.track {
			if len(track.name) > 0 {
				wr.Write(meta.Instrument(track.name))
			}

			for _, note := range track.notes.msg {
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
	midiFileOutPtr := flag.String("file_out", "nil", "path and name of midi file to save")

	flag.Parse()
	if *actionPtr == "read" {
		readMidiFile(*midiFileInPtr)
		//for _, val := range gTracks.track {
		//	fmt.Println(val.name)
		//	for _, nt := range val.notes.msg {
		//		fmt.Println("note is")
		//		fmt.Println(nt)
		//	}
		//	//
		//}

	}
	if *actionPtr == "copy" {
		copyMidiFile(*midiFileInPtr, *midiFileOutPtr)
	}
	//fmt.Println(len(gScore.track))

}
