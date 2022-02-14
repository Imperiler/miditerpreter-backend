class Track:
    name = str
    time = int

    class InstrumentName:
        name = str
        time = int

    class TrackEnd:
        time = int

    class Note:
        channel = int
        note = int
        velocity = int
        time = int
        note_on = bool


class TimeSignature:
    numerator = int
    denomonator = int
    clocks_per_tick = int
    notated_32nd_notes_per_beat = 8
    time = 0


class KeySignature:
    key = str
    time = int


class Marker:
    text = str
    time = int


class SetTempo:
    tempo = int
    time = int
