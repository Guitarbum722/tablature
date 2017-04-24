package tabio

import (
	"bufio"
	"github.com/Guitarbum722/go-tabs/instrument"
	"github.com/pkg/errors"
	"io"
)

// TablatureWriter embeds a buffered writer
// The wrapPosition can be used to choose how far the tablature section will reach before wrapping to the next.
// 20 is the default.
type TablatureWriter struct {
	*bufio.Writer
	WrapPosition int
	totalLength  int
	tb           tablatureBuilder
}

type tablatureBuilder struct {
	builder map[byte][]byte
}

// NewTablatureWriter creates a buffered writer to be used for staging tablature
func NewTablatureWriter(w io.Writer, pos int) *TablatureWriter {

	return &TablatureWriter{
		bufio.NewWriter(w),
		pos,
		0,
		tablatureBuilder{
			make(map[byte][]byte),
		},
	}
}

// StageTablature writes the Instrument's current tablature to w for buffering.
// The purpose is only to stage or buffer the current tablature but it does not
// write the tablature to a file.
func StageTablature(i instrument.Instrument, w *TablatureWriter) {

	for k, v := range i.Fretboard() {

		w.tb.builder[k] = append(w.tb.builder[k], []byte(v)...)
		w.totalLength = len(w.tb.builder[k])
	}

	return
}

// ExportTablature will flush the bufferred writer to the io.Writer of which it was initialized
func ExportTablature(i instrument.Instrument, w *TablatureWriter) error {

	var done int

	for ; done < w.totalLength; done += w.WrapPosition {

		for _, v := range i.Order() {

			w.Write([]byte{v, ':', ' '})

			// use the TablatureWriter's wrap position to determine when to continue a new
			// tablature section in the output
			if w.totalLength < w.WrapPosition {
				if _, err := w.Write(w.tb.builder[v][:w.totalLength]); err != nil {
					return errors.Wrap(err, "write to buffer failed")
				}
			} else if (done + w.WrapPosition) < w.totalLength {
				if _, err := w.Write(w.tb.builder[v][done:(done + w.WrapPosition)]); err != nil {
					return errors.Wrap(err, "write to buffer failed")
				}
			} else {
				if _, err := w.Write(w.tb.builder[v][done:(w.totalLength)]); err != nil {
					return errors.Wrap(err, "write to buffer failed")
				}
			}
			w.WriteByte('\n')
		}
		w.WriteByte('\n')
	}

	w.Flush()

	return nil
}
