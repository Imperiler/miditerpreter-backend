import mido
import argparse


def read_file(midi_file_path):
    meta_dict = {}
    # note_dict = {}
    file_in = mido.MidiFile(midi_file_path)
    for i, track in enumerate(file_in.tracks):
        for msg in track:
            msg_string = str(msg)
            if "note" in msg_string:
                if "MetaMessage" not in msg_string:
                    note_dict = dict()
                    print(msg)


def main():
    parser = argparse.ArgumentParser()

    parser.add_argument('-f', '--file', help='full path to file on which to operate', required=False)
    parser.add_argument('-a', '--action', help='Action to take on specified file', required=True)

    args = parser.parse_args()

    if args.action == 'read':
        read_file(args.file)


if __name__ == "__main__":
    main()
