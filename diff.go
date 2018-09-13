package main
 
import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)
 
//@@ -0,0 +1,179 @@
//LineInfo: @@ -0,0 +1,179 @@
//Function: null
 
//@@ -32,11 +32,11 @@ func Test
//LineInfo:  @@ -32,11 +32,11 @@
//Function:  func Test
type DiffItemFrag struct {
	LineInfo string
	Function string
	Related  int64
	Added    int64
	Removed  int64
 
	// success related lines with be merged into one
	RelatedSeg int64
	// success added lines with be merged into one
	AddedSeg int64
	// success removed lines with be merged into one
	RemovedSeg int64
 
	lastType LineType
}
 
func NewDiffItemFrag(lines []Line) (DiffItemFrag, error) {
	var frag DiffItemFrag
	frag.lastType = LineNull
	for _, l := range lines {
		if l.Type == LineFragHead {
			if err := frag.SetHead(l.Data); err != nil {
				break
			}
		} else if frag.HasHead() {
			frag.SetDetail(l.Type)
		}
	}
 
	if !frag.HasHead() {
		return frag, errors.New("Fail to parse frag file")
	}
 
	return frag, nil
}
 
func (f *DiffItemFrag) HasHead() bool {
	return f.LineInfo != ""
}
 
func (f *DiffItemFrag) SetHead(head string) error {
	if f.HasHead() {
		return errors.New("Already has header")
	}
 
	heads := LineRegexp[LineFragHead].FindStringSubmatch(head)
	if len(heads) != 3 {
		return fmt.Errorf("Invalid head: %s", head)
	}
 
	f.LineInfo = heads[1]
	f.Function = heads[2]
 
	return nil
}
 
func (f *DiffItemFrag) SetDetail(t LineType) error {
	switch t {
	case LineAdded:
		f.Added++
		if f.lastType != LineAdded {
			f.AddedSeg++
		}
	case LineRemoved:
		f.Removed++
		if f.lastType != LineRemoved {
			f.RemovedSeg++
		}
	case LineOther:
		f.Related++
		if f.lastType != LineOther {
			f.RelatedSeg++
		}
	default:
		return fmt.Errorf("Invalid diff item frag type: %d", t)
	}
	f.lastType = t
	return nil
}

type DiffSummary struct {
	Added   int64
	Removed int64
}
 
/*
diff --git a/controller/webv1/dockerhub.go b/controller/webv1/dockerhub.go
index 5979b477..e928601e 100644
--- a/controller/webv1/dockerhub.go
+++ b/controller/webv1/dockerhub.go
@@ -149,6 +149,10 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
	                }
		}
*/
type DiffItem struct {
	File  string
	Label LineType
	Frags []DiffItemFrag
}
 
func NewDiffItem(lines []Line) (DiffItem, error) {
	var item DiffItem
	lastLine := -1
	begin := false
	eof := false
	item.Label = LineNull
	for i, l := range lines {
		switch l.Type {
		case LineDiff:
			if begin {
				eof = true
				break
			}
			begin = true
			item.File = l.Data
		case LineNew:
			item.Label = l.Type
		case LineDeleted:
			item.Label = l.Type
		case LineFragHead:
			//Not start to parse yet
			if !begin {
				break
			}
			if lastLine < 0 {
				lastLine = i
				break
			}
			frag, err := NewDiffItemFrag(lines[lastLine : i+1])
			lastLine = i
			if err != nil {
				log.Debugf("Error in create diff frag: %s", err.Error())
			} else {
				item.Frags = append(item.Frags, frag)
			}
		}
		if eof {
			break
		}
	}
	if !begin || lastLine < 0 {
		return item, errors.New("Fail to create diff item")
	}
	frag, err := NewDiffItemFrag(lines[lastLine:])
	if err != nil {
		log.Debugf("Error in create diff frag: %s", err.Error())
	} else {
		item.Frags = append(item.Frags, frag)
	}
 
	return item, nil
}
 
func NewDiffItems(lines []Line, search SearchCondition) ([]DiffItem, error) {
	var items []DiffItem
	lastLine := -1
	begin := false
	for i, l := range lines {
		switch l.Type {
		case LineDiff:
			if !begin {
				begin = true
				lastLine = i
				continue
			}
			if item, err := NewDiffItem(lines[lastLine : i+1]); err != nil {
				log.Debugf("Error in create diff item: %s", err.Error())
			} else if item.Match(search) {
				items = append(items, item)
			}
			lastLine = i
		}
	}
	if !begin || lastLine < 0 {
		return nil, errors.New("Fail to create diff item")
	}
	if item, err := NewDiffItem(lines[lastLine:]); err != nil {
		log.Debugf("Error in create diff item: %s", err.Error())
	} else {
		items = append(items, item)
	}
 
	return items, nil
}
 
func (di *DiffItem) Match(search SearchCondition) bool {
	for _, s := range search.Skip {
		switch s.SearchType {
		case LineDiff:
			// match all
			if strings.Contains(di.File, s.SearchContent) {
				return false
			}
		}
	}
 
	for _, s := range search.Contain {
		switch s.SearchType {
		case LineDiff:
			// match all
			if !strings.Contains(di.File, s.SearchContent) {
				return false
			}
		}
	}
 
	return true
}
 
func (di *DiffItem) Summary() DiffSummary {
	var ds DiffSummary
	for _, f := range di.Frags {
		ds.Added += f.Added
		ds.Removed += f.Removed
	}
 
	return ds
}
 
func GetDiffFiles(items []DiffItem) map[string]int64 {
	df := make(map[string]int64)
	for _, item := range items {
		val, ok := df[item.File]
		if !ok {
			val = 0
		}
		is := item.Summary()
		val += is.Added
		val -= is.Removed
 
		df[item.File] = val
	}
 
	return df
}
 
