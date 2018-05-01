package backup

import "testing"

var classificationCases = []struct {
	ext                string
	wantClassification fileClassification
}{
	{".jpg", imageFile},
	{".JpG", imageFile},
	{".mp4", videoFile},
	{".txt", unclassifiedFile},
	{"", unclassifiedFile},
}

func TestClassifyExt(t *testing.T) {
	for _, c := range classificationCases {
		actual := classifyExt(c.ext)
		if actual != c.wantClassification {
			t.Errorf("[%s] Expected %d, got %d", c.ext, c.wantClassification, actual)
		}
	}
}
