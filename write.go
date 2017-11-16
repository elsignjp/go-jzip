package jzip

import "io"
import "bufio"
import "fmt"

// WriterType Writerのタイプ
type WriterType int

const (
	// JSONWriter JSON形式
	JSONWriter WriterType = iota
	// SQLWriter SQL形式
	SQLWriter
)

// WriterConfig Writerの設定項目
type WriterConfig struct {
	TableName        string
	LocalCodeField   string
	ZipCodeField     string
	PrefecturesField string
	CityField        string
	TownField        string
}

// NewWriter Writerを初期化する関数
func NewWriter(format WriterType, out io.Writer, cfg WriterConfig) Writer {
	if format == SQLWriter {
		return &sqlWriter{out: bufio.NewWriter(out), cfg: cfg, fst: true}
	}
	return &jsonWriter{out: bufio.NewWriter(out), cfg: cfg, fst: true}
}

// Writer Writerの構造体
type Writer interface {
	Write(record []string) error
	Close()
}

type jsonWriter struct {
	out *bufio.Writer
	cfg WriterConfig
	fst bool
}

func (w *jsonWriter) Write(record []string) error {
	if w.fst {
		w.out.Write([]byte{'[', '\n'})
		w.fst = false
	} else {
		w.out.Write([]byte{',', '\n'})
	}
	_, err := w.out.WriteString(fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":"%s","%s":"%s","%s":"%s"}`,
		w.cfg.LocalCodeField, record[0],
		w.cfg.ZipCodeField, record[2],
		w.cfg.PrefecturesField, record[6],
		w.cfg.CityField, record[7],
		w.cfg.TownField, record[8]))
	return err
}

func (w *jsonWriter) Close() {
	w.out.Write([]byte{'\n', ']'})
	w.out.Flush()
}

type sqlWriter struct {
	out *bufio.Writer
	cfg WriterConfig
	fst bool
}

func (w *sqlWriter) Write(record []string) error {
	if w.fst {
		w.out.WriteString("BEGIN;\n")
		w.fst = false
	}
	_, err := w.out.WriteString(fmt.Sprintf("INSERT INTO `%s` (`%s`,`%s`,`%s`,`%s`,`%s`) VALUES (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\");\n",
		w.cfg.TableName,
		w.cfg.LocalCodeField, w.cfg.ZipCodeField, w.cfg.PrefecturesField, w.cfg.CityField, w.cfg.TownField,
		record[0], record[2], record[6], record[7], record[8]))
	return err
}

func (w *sqlWriter) Close() {
	w.out.WriteString("COMMIT;\n")
	w.out.Flush()
}
