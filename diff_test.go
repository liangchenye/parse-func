package main
 
import (
	"testing"
 
	"github.com/stretchr/testify/assert"
)
 
func TestDiffItemFragAll(t *testing.T) {
	testData := []byte(`
@@ -141,7 +141,7 @@ func Test
	func (gc *GoCrypto) GetPrivKey() ([]byte, error) {
		var a string
-		var b string
-       if tlsKeyByte == nil {
+       if len(tlsKeyByte) == 0 {
+       }
@@ -241,7 +241,7 @@ func Test2
`)
	data := ParseData(testData)
	frag, err := NewDiffItemFrag(data)
	assert.Nil(t, err, "new item should not get error")
	assert.Equal(t, "@@ -141,7 +141,7 @@", frag.LineInfo, "info not equal")
	assert.Equal(t, " func Test", frag.Function, "function not equal")
	assert.Equal(t, int64(2), frag.Related, "related not equal")
	assert.Equal(t, int64(1), frag.RelatedSeg, "related seg not equal")
	assert.Equal(t, int64(2), frag.Added, "added not equal")
	assert.Equal(t, int64(1), frag.AddedSeg, "added seg not equal")
	assert.Equal(t, int64(2), frag.Removed, "removed not equal")
	assert.Equal(t, int64(1), frag.RemovedSeg, "removed seg not equal")
 
	// to get function
	assert.Equal(t, "", frag.GetFunction(), "must be nil")
	frag.Function = "static void int demo()"
	assert.Equal(t, "demo", frag.GetFunction(), "must be nil")

	testData2 := []byte(`
	func (gc *GoCrypto) GetPrivKey() ([]byte, error) {
-       if tlsKeyByte == nil {
+       if len(tlsKeyByte) == 0 {
+       }
invalid
`)
	data = ParseData(testData2)
	_, err = NewDiffItemFrag(data)
	assert.NotNil(t, err, "new item should not get error")
 
	frag = DiffItemFrag{}
	err = frag.SetHead("invalid")
	assert.NotNil(t, err, "fail to set head")
}
 
func TestNewDiffItem(t *testing.T) {
	testData := []byte(`
diff --git a/controller/webv1/dockerhub.go b/controller/webv1/dockerhub.go
new file abcd
index 5979b477..e928601e 100644
--- a/controller/webv1/dockerhub.go
+++ b/controller/webv1/dockerhub.go
@@ -149,6 +149,10 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
-               }
@@ -157,6 +161,34 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
+func GetDockerhubDescription(namespace, repository string) (int, []byte) {
`)
 
	data := ParseData(testData)
	item, err := NewDiffItem(data)
	assert.Nil(t, err, "new diff item error")
	assert.Equal(t, 2, len(item.Frags), "get frags error")
	assert.Equal(t, LineNew, item.Label, "get label error")
	assert.Equal(t, " --git a/controller/webv1/dockerhub.go b/controller/webv1/dockerhub.go", item.File, "fail to get diff info")

	// get the real file for openssl.. 
	assert.Equal(t, "b/controller/webv1/dockerhub.go", item.GetFile(), "fail to get fail")

	summary := item.Summary()
	assert.Equal(t, int64(1), summary.Added, "get summary added error")
	assert.Equal(t, int64(1), summary.Removed, "get summary removed error")
 
	testData2 := []byte(`
index 5979b477..e928601e 100644
--- a/controller/webv1/dockerhub.go
+++ b/controller/webv1/dockerhub.go
@@ -149,6 +149,10 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
                }
        }
`)
 
	data = ParseData(testData2)
	_, err = NewDiffItem(data)
	assert.NotNil(t, err, "new diff item error")
}
 
func TestNewDiffItems(t *testing.T) {
	testData := []byte(`
diff --git a/controller/webv1/dockerhub.go b/controller/webv1/dockerhub.go
new file abcd
@@ -149,6 +149,10 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
                }
diff --git a/controller/webv1/dockerhub1.go b/controller/webv1/dockerhub1.go
new file abcd
@@ -157,6 +161,34 @@ func GetDockerhubRepos(filter map[string]string) (int, []byte) {
+       router := fmt.Sprintf("%s%s/%s/%s/", DockerHubHost, DockerHubRepositoriesRouter, namespac
`)
 
	data := ParseData(testData)
	items, err := NewDiffItems(data, SearchCondition{})
	assert.Nil(t, err, "new diff items error")
	assert.Equal(t, 2, len(items), "get items error")
}
