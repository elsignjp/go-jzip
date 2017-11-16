package jzip

import (
	"strings"
)

// NewNormalize Normalizeを初期化する関数
func NewNormalize() *Normalize {
	return &Normalize{}
}

// Normalize Normalizeの構造体
type Normalize struct {
	join string
}

// Add CSVの1行を受け取り、加工を行い、スライスで返す
func (n *Normalize) Add(record []string) [][]string {
	if n.join != "" {
		record[8] = n.join + record[8]
		n.join = ""
	}
	if strings.ContainsRune(record[8], '（') && !strings.ContainsRune(record[8], '）') {
		n.join = record[8]
		return nil
	}
	main, args := splitMainAndArgs(record[8])
	rs := make([][]string, 0)
	for _, m := range main {
		if args == nil {
			rc := make([]string, len(record))
			copy(rc, record)
			rc[8] = m
			rs = append(rs, rc)
		} else {
			for _, a := range args {
				rc := make([]string, len(record))
				copy(rc, record)
				rc[8] = m + a
				rs = append(rs, rc)
			}
		}
	}
	return rs
}

func splitMainAndArgs(t string) ([]string, []string) {
	var main string
	var arg string
	if i1, i2, ok := indexOfRunePair(t, '（', '）'); ok {
		main = t[:i1]
		arg = t[i1+3 : i2]
	} else {
		main = t
	}
	return parseMainString(main), parseArgString(arg)
}

func parseMainString(m string) []string {
	ms := strings.Split(m, "、")
	rs := make([]string, 0, len(ms))
	for _, s := range ms {
		s = transformMainString(s)
		if !containsString(rs, s) {
			rs = append(rs, s)
		}
	}
	return rs
}

func parseArgString(a string) []string {
	if a == "" {
		return nil
	}
	args := transformArgsString(a)
	if args == nil || len(args) == 0 {
		return nil
	}
	rs := make([]string, 0)
	for _, s := range args {
		s = transformArgString(s)
		if !containsString(rs, s) {
			rs = append(rs, s)
		}
	}
	return rs
}

func transformMainString(m string) string {
	if m == "以下に掲載がない場合" {
		return ""
	}
	if strings.HasSuffix(m, "の次に番地がくる場合") {
		return ""
	}
	if m != "一円" && strings.HasSuffix(m, "一円") {
		return ""
	}
	if strings.Contains(m, "地割") {
		return getPrefix(m)
	}
	return m
}

func transformArgsString(a string) []string {
	if a == "その他" ||
		a == "無番地" ||
		a == "無番地を除く" ||
		a == "次のビルを除く" ||
		a == "全域" ||
		a == "堤防用地" ||
		a == "地階・階層不明" ||
		a == "大字" ||
		a == "丁目" ||
		a == "番地" ||
		strings.HasSuffix(a, "地区") ||
		strings.HasSuffix(a, "国有林") ||
		strings.HasSuffix(a, "空港内") ||
		strings.Contains(a, "自衛隊") ||
		strings.Contains(a, "空港関係") {
		return nil
	}
	for {
		if i1, i2, ok := indexOfRunePair(a, '「', '」'); ok {
			a = a[:i1] + a[i2+3:]
			continue
		}
		break
	}
	return strings.Split(a, "、")
}

func transformArgString(a string) string {
	if a == "大字" ||
		a == "丁目" ||
		a == "番地" {
		return ""
	}
	if strings.ContainsRune(a, '～') {
		return ""
	}
	if strings.HasSuffix(a, "丁目") ||
		strings.HasSuffix(a, "番地") ||
		strings.HasSuffix(a, "を除く") ||
		strings.HasSuffix(a, "を含む") ||
		strings.HasSuffix(a, "以降") ||
		strings.HasSuffix(a, "以内") ||
		strings.HasSuffix(a, "以上") ||
		strings.HasSuffix(a, "以下") ||
		strings.HasSuffix(a, "以外") {
		return ""
	}
	if isNumOnly(a) {
		return ""
	}
	if strings.HasSuffix(a, "線") && isNumOnly(a[:len(a)-3]) {
		return ""
	}
	if strings.HasSuffix(a, "その他") && len(a) > 9 {
		return a[:len(a)-9]
	}
	return a
}

func isNumOnly(a string) bool {
	num := "１２３４５６７８９０－～丁目番地の・"
	for _, r := range a {
		if !strings.ContainsRune(num, r) {
			return false
		}
	}
	return true
}

func getPrefix(m string) string {
	num := "１２３４５６７８９０－～"
	p := make([]rune, 0, len(m))
	for _, r := range m {
		if strings.ContainsRune(num, r) {
			break
		}
		p = append(p, r)
	}
	return string(p)
}

func indexOfRunePair(s string, r1 rune, r2 rune) (int, int, bool) {
	i1 := strings.IndexRune(s, r1)
	i2 := strings.IndexRune(s, r2)
	return i1, i2, i1 != -1 && i2 != -1
}
