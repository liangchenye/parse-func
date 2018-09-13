package main
 
import (
	"errors"
	"fmt"
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
