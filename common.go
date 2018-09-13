package utils
 
import (
	"regexp"
	"strings"
)
 
// match compiles the string to a regular expression.
var match = regexp.MustCompile
 
var (
	// git@code.huawei.com:paas/paas-dockyard.git
	// http://code.huawei.com/paas/paas-dockyard.git
	GitURLSsh  = match("^git@(.*):(.*)/(.*).git$")
	GitURLHttp = match("^http://(.*)/(.*)/(.*).git$")
	GitURL     = either(GitURLSsh, GitURLHttp)
)
 
type LineType int
 
const (
	LineAdded LineType = iota
	LineRelated
	LineRemoved
	LineFragHead
	LineDiff
	LineCommit
	LineAuthor
	LineDate
	LineMerge
	LineOther
	LineDeleted
	LineNew
	LineUnknown
	LineNull
)
 
type Line struct {
	Data string
	Type LineType
}
 
// We sorted them by usage, 'added' is more frequent.
// TODO: +++, --- is not used, but should better be excluded
var LineRegexp = map[LineType]*regexp.Regexp{
	LineAdded: match("^[+]+.*"),
	//	LineRelated:  match("^\\s.*"),
	LineRemoved:  match("^[-]+.*"),
	LineFragHead: match("^(@@.*@@)(.*)"),
	LineOther:    match("^\\s(.*)"),
	LineDiff:     match("^diff(.*)"),
	LineCommit:   match("^commit\\s(.*)"),
	LineAuthor:   match("^Author:\\s(.*)"),
	LineDate:     match("^Date:\\s(.*)"),
	LineMerge:    match("^Merge:\\s(.*)"),
	LineDeleted:  match("^deleted\\sfile\\s(.*)"),
	LineNew:      match("^new\\sfile\\s(.*)"),
}
 
func ParseLine(line string) Line {
	var d string
	for t, r := range LineRegexp {
		if r.MatchString(line) {
			switch t {
			case LineOther:
				fallthrough
			case LineDiff:
				fallthrough
			case LineCommit:
				fallthrough
			case LineAuthor:
				fallthrough
			case LineDate:
				fallthrough
			case LineMerge:
				fallthrough
			case LineDeleted:
				fallthrough
			case LineNew:
				d = LineRegexp[t].FindStringSubmatch(line)[1]
			default:
				d = line
			}
			return Line{Data: d, Type: t}
		}
	}
 
	return Line{Data: line, Type: LineUnknown}
}
 
func ParseData(data []byte) (diffLines []Line) {
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		diffLines = append(diffLines, ParseLine(l))
	}
 
	return
}
 
// either wraps the 'a or b' operation
func either(res ...*regexp.Regexp) *regexp.Regexp {
	var s string
	for i, re := range res {
		if i > 0 {
			s += `|`
		}
		s += re.String()
	}
 
	return match(`(?:` + s + `)`)
}
 
type SearchItem struct {
	SearchType    LineType
	SearchContent string
}
type SearchCondition struct {
	Skip    []SearchItem
	Contain []SearchItem
}
