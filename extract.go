package jzip

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Extract ZIPファイルからCSVファイルを抽出し、各行をチャンネルで非同期に返す
func Extract(ctx context.Context, in string) (chan []string, error) {
	rd, err := zip.OpenReader(in)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var cf *zip.File
	for _, f := range rd.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			cf = f
			break
		}
	}
	if cf == nil {
		rd.Close()
		return nil, errors.New("not found csv in archive")
	}
	fp, err := cf.Open()
	if err != nil {
		rd.Close()
		return nil, errors.WithStack(err)
	}
	ch := make(chan []string, 1000)
	go func() {
		defer func() {
			fp.Close()
			rd.Close()
			close(ch)
		}()
		sc := csv.NewReader(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
		sc.LazyQuotes = true
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			r, err := sc.Read()
			if err == io.EOF {
				return
			} else if err != nil {
				return
			}
			ch <- r
		}
	}()
	return ch, nil
}
