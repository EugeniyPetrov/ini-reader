package ini_reader

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func mustOpenFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return f
}

func newFileReader(filename string) *Reader {
	return NewReader(mustOpenFile(filename))
}

func mustParseFile(t *testing.T, filename string) []*Section {
	res, err := newFileReader(filename).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func parseString(s string) ([]*Section, error) {
	return NewReader(strings.NewReader(s)).ReadAll()
}

func mustParseString(t *testing.T, s string) []*Section {
	res, err := parseString(s)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func assertINI(t *testing.T, actual []*Section, expected []*Section) bool {
	if len(expected) != len(actual) {
		t.Errorf("expected %d sections, got %d", len(expected), len(actual))
		return false
	}

	for i := range expected {
		if expected[i].Name != actual[i].Name {
			t.Errorf("expected section %d to have name %s, got %s", i, expected[i].Name, actual[i].Name)
			return false
		}

		if len(expected[i].Properties) != len(actual[i].Properties) {
			t.Errorf(
				"expected section %q to have %d properties, got %d",
				expected[i].Name,
				len(expected[i].Properties),
				len(actual[i].Properties),
			)
			return false
		}

		for k, v := range expected[i].Properties {
			if _, ok := actual[i].Properties[k]; !ok {
				t.Errorf(
					"expected section %q to have property %s=%v, got none",
					expected[i].Name,
					k, v,
				)
				return false
			}

			if actual[i].Properties[k] != v {
				t.Errorf(
					"expected section %q to have property %s=%v, got %s=%v",
					expected[i].Name,
					k, v,
					k, actual[i].Properties[k],
				)
				return false
			}
		}
	}

	return true
}

func TestReaderBasic(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/basic.ini"),
		[]*Section{
			{
				Name: "General",
				Properties: map[string]any{
					"Compiler": "FreePascal",
				},
			},
		})
}

func TestReaderExampleSettings(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/example.settings.ini"),
		[]*Section{
			{
				Name: "",
				Properties: map[string]any{
					"domain": false,
				},
			},
			{
				Name: "sql",
				Properties: map[string]any{
					"server":   "localhost",
					"database": "emoncms",
					"username": "emoncms",
					"password": "password",
					"dbtest":   true,
				},
			},
			{
				Name: "redis",
				Properties: map[string]any{
					"enabled": true,
				},
			},
			{
				Name: "mqtt",
				Properties: map[string]any{
					"enabled":  false,
					"user":     "username",
					"password": "password",
				},
			},
			{
				Name: "feed",
				Properties: map[string]any{
					"engines_hidden":         "[0,10]",
					"redisbuffer[enabled]":   false,
					"phpfina[datadir]":       "/var/opt/emoncms/phpfina/",
					"phptimeseries[datadir]": "/var/opt/emoncms/phptimeseries/",
				},
			},
			{
				Name: "interface",
				Properties: map[string]any{
					"feedviewpath": "graph/",
				},
			},
			{
				Name:       "public_profile",
				Properties: map[string]any{},
			},
			{
				Name:       "smtp",
				Properties: map[string]any{},
			},
			{
				Name: "log",
				Properties: map[string]any{
					"level": int64(2),
				},
			},
		},
	)
}

func TestReaderExample1(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/example_1.ini"),
		[]*Section{
			{
				Name: "DEVICE1",
				Properties: map[string]any{
					"UseTheseDomainSizes": int64(1),
					"UseCounts":           int64(0),
					"OneCoilWrite":        int64(0),
					"OneRegWrite":         int64(0),
					"ConservesConn":       int64(1),
					"ConnSecondary":       int64(0),
					"COILS":               int64(65535),
					"DISC INPUTS":         int64(65535),
					"INPUT REG.":          int64(65535),
					"HOLDING REG.":        int64(65535),
					"GEN REF FILE1":       int64(0),
					"GEN REF FILE2":       int64(0),
					"GEN REF FILE3":       int64(0),
					"GEN REF FILE4":       int64(0),
					"GEN REF FILE5":       int64(0),
					"GEN REF FILE6":       int64(0),
					"GEN REF FILE7":       int64(0),
					"GEN REF FILE8":       int64(0),
					"GEN REF FILE9":       int64(0),
					"GEN REF FILE10":      int64(0),
					"DP_INPUT REG.":       int64(0),
					"DP_HOLDING REG.":     int64(0),
				},
			},
			{
				Name: "DEVICE2",
				Properties: map[string]any{
					"UseTheseDomainSizes": int64(1),
					"UseCounts":           int64(0),
					"OneCoilWrite":        int64(0),
					"OneRegWrite":         int64(0),
					"ConservesConn":       int64(1),
					"ConnSecondary":       int64(0),
					"COILS":               int64(65535),
					"DISC INPUTS":         int64(65535),
					"INPUT REG.":          int64(65535),
					"HOLDING REG.":        int64(65535),
					"GEN REF FILE1":       int64(0),
					"GEN REF FILE2":       int64(0),
					"GEN REF FILE3":       int64(0),
					"GEN REF FILE4":       int64(0),
					"GEN REF FILE5":       int64(0),
					"GEN REF FILE6":       int64(0),
					"GEN REF FILE7":       int64(0),
					"GEN REF FILE8":       int64(0),
					"GEN REF FILE9":       int64(0),
					"GEN REF FILE10":      int64(0),
					"DP_INPUT REG.":       int64(0),
					"DP_HOLDING REG.":     int64(0),
				},
			},
		},
	)
}

func TestReaderExample2(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/example_2.ini"),
		[]*Section{
			{
				Name: "owner",
				Properties: map[string]any{
					"name":         "John Doe",
					"organization": "Acme Widgets Inc.",
				},
			},
			{
				Name: "database",
				Properties: map[string]any{
					"server": "192.0.2.62",
					"port":   int64(143),
					"file":   "payroll.dat",
				},
			},
		},
	)
}

func TestReaderSyntax(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/syntax.ini"),
		[]*Section{
			{
				Name: "EmLibraryInterface",
				Properties: map[string]any{
					"DefaultName":   "d:\\vamgr\\manager\\vavm020304_dev.dat",
					"ServerAddress": "192.168.1.101",
					"OpenReadOnly":  false,
				},
			},
		},
	)
}

func TestReaderTest1(t *testing.T) {
	assertINI(
		t,
		mustParseString(t, `😀   =   😃
        ; comment
prop=`),
		[]*Section{
			{
				Name: "",
				Properties: map[string]any{
					"😀":    "😃",
					"prop": nil,
				},
			},
		},
	)
}

func TestReaderTest2(t *testing.T) {
	assertINI(
		t,
		mustParseString(t, `prop`),
		[]*Section{
			{
				Name:       "",
				Properties: map[string]any{},
			},
		},
	)
}

func TestReaderTest3(t *testing.T) {
	_, err := parseString(`prop="hello`)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrUnexpectedEOF) {
		t.Fatalf("expected ErrUnexpectedEOF, got %v", err)
	}
}

func TestReaderLitePHPBrowscap(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/lite_php_browscap_part.ini"),
		[]*Section{
			{
				Name: "GJK_Browscap_Version",
				Properties: map[string]any{
					"Version":  int64(6001007),
					"Released": "Sun, 16 Jun 2024 11:25:25 +0000",
					"Format":   "php",
					"Type":     "LITE",
				},
			},
			{
				Name: "DefaultProperties",
				Properties: map[string]any{
					"Comment":        "DefaultProperties",
					"Browser":        "DefaultProperties",
					"Version":        "0.0",
					"Platform":       "unknown",
					"isMobileDevice": "false",
					"isTablet":       "false",
					"Device_Type":    "unknown",
				},
			},
			{
				Name: "Chromium 135.0",
				Properties: map[string]any{
					"Parent":      "DefaultProperties",
					"Comment":     "Chromium 135.0",
					"Browser":     "Chromium",
					"Version":     "135.0",
					"Platform":    "Linux",
					"Device_Type": "Desktop",
				},
			},
		},
	)
}

func TestReaderLuaHttpClient(t *testing.T) {
	assertINI(
		t,
		mustParseFile(t, "fixtures/lua_http_client.ini"),
		[]*Section{
			{
				Name: "lua-resty-http/0.07 *",
				Properties: map[string]any{
					"Parent":   "lua http client",
					"Version":  "0.07",
					"MinorVer": "07",
				},
			},
		},
	)
}
