package markup

import "testing"
import "time"

func TestConvertToJSON(t *testing.T) {
	t.Log(convertToJSON(42))
}

func TestFormatTime(t *testing.T) {
	t.Log(formatTime(time.Now(), "2006"))
}
