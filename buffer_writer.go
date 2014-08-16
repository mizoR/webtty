package webtty

import "bytes"

type BufferWriter struct {
}

func NewBufferWriter() *BufferWriter {
	writer := &BufferWriter{}
	return writer
}

func (writer BufferWriter) write(buf *bytes.Buffer, r rune) {
	switch r {
	case 34: // `"`
		buf.WriteString("&quot;")
	case 38: // `&`
		buf.WriteString("&amp;")
	case 39: // `'`
		buf.WriteString("&#039;")
	case 60: // `<`
		buf.WriteString("&lt;")
	case 62: // `>`
		buf.WriteString("&gt;")
	default:
		buf.WriteRune(r)
	}
}

func (writer BufferWriter) writeLF(buf *bytes.Buffer) {
	buf.WriteRune(10) // LF
}

func (writer BufferWriter) writeCursor(buf *bytes.Buffer) {
	buf.WriteString("<div class='cursor'></div>")
}
