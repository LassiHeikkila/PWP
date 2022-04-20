package json

import (
	stdjson "encoding/json"
	"fmt"
	"testing"
	"time"
)

type TestStruct struct {
	Dur Duration `json:"duration"`
}

func TestMarshalling(t *testing.T) {
	tests := []struct {
		input TestStruct
		want  string
	}{
		{input: TestStruct{Dur: Duration{time.Second}}, want: `{"duration":"1s"}`},
		{input: TestStruct{Dur: Duration{time.Minute}}, want: `{"duration":"1m0s"}`},
		{input: TestStruct{Dur: Duration{time.Hour}}, want: `{"duration":"1h0m0s"}`},
		{input: TestStruct{Dur: Duration{24 * time.Hour}}, want: `{"duration":"24h0m0s"}`},
		{input: TestStruct{Dur: Duration{7 * 24 * time.Hour}}, want: `{"duration":"168h0m0s"}`},
		{input: TestStruct{Dur: Duration{5 * time.Millisecond}}, want: `{"duration":"5ms"}`},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			gotB, err := stdjson.Marshal(&tc.input)
			if err != nil {
				t.Fatal("error marshalling:", err)
			}
			got := string(gotB)
			if tc.want != got {
				t.Fatalf("wanted:\n\t%s\ngot:\n\t%s\n", tc.want, got)
			}
		})
	}
}

func TestUnmarshalling(t *testing.T) {
	tests := []struct {
		input string
		want  TestStruct
	}{
		{want: TestStruct{Dur: Duration{time.Second}}, input: `{"duration":"1s"}`},
		{want: TestStruct{Dur: Duration{time.Minute}}, input: `{"duration":"1m0s"}`},
		{want: TestStruct{Dur: Duration{time.Hour}}, input: `{"duration":"1h0m0s"}`},
		{want: TestStruct{Dur: Duration{24 * time.Hour}}, input: `{"duration":"24h0m0s"}`},
		{want: TestStruct{Dur: Duration{7 * 24 * time.Hour}}, input: `{"duration":"168h0m0s"}`},
		{want: TestStruct{Dur: Duration{5 * time.Millisecond}}, input: `{"duration":"5ms"}`},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var got TestStruct
			err := stdjson.Unmarshal([]byte(tc.input), &got)
			if err != nil {
				t.Fatal("error marshalling:", err)
			}
			if tc.want != got {
				t.Fatalf("wanted:\n\t%v\ngot:\n\t%v\n", tc.want, got)
			}
		})
	}
}

func TestUnmarshallingBad(t *testing.T) {
	tests := []struct {
		input string
	}{
		{input: `{"duration":"7d"}`}, // bad unit
		{input: `{"duration":1000}`}, // not a string
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var got TestStruct
			err := stdjson.Unmarshal([]byte(tc.input), &got)
			if err == nil {
				t.Fatal("no error marshalling")
			}
		})
	}
}
