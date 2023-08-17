package code

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
)

func Test_slogtest(t *testing.T) {
	// range:slogtest
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)
	err := slogtest.TestHandler(
		h, func() []map[string]any {
			ss := strings.Split(buf.String(), "\n")
			ret := make([]map[string]any, 0, len(ss))

			for _, s := range ss {
				if s == "" {
					continue
				}

				result := map[string]any{}
				err := json.Unmarshal([]byte(s), &result)
				if err != nil {
					t.Fatal(err)
				}

				ret = append(ret, result)
			}
			return ret
		},
	)
	if err != nil {
		t.Error(err)
	}
	// range.end
}
