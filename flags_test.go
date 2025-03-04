package main_test

import (
	"testing"

	main "github.com/NatoBoram/Sc2Bot"
)

func TestParseArg_EmptyArgs(t *testing.T) {
	advance, name, value := main.ParseArg([]string{}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "" {
		t.Errorf("Expected name to be empty, got %q", name)
	}
	if value != "" {
		t.Errorf("Expected value to be empty, got %q", value)
	}
}

func TestParseArg_IndexOutOfBounds(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag"}, 1)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "" {
		t.Errorf("Expected name to be empty, got %q", name)
	}
	if value != "" {
		t.Errorf("Expected value to be empty, got %q", value)
	}
}

func TestParseArg_NotAFlag(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"notaflag"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be \"0\", got %d", advance)
	}
	if name != "" {
		t.Errorf("Expected name to be empty, got %q", name)
	}
	if value != "" {
		t.Errorf("Expected value to be empty, got %q", value)
	}
}

func TestParseArg_SimpleFlagWithDash(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "true" {
		t.Errorf("Expected value to be \"true\", got %q", value)
	}
}

func TestParseArg_SimpleFlagWithDoubleDash(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--flag"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "true" {
		t.Errorf("Expected value to be \"true\", got %q", value)
	}
}

func TestParseArg_FlagWithEquals(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag=value"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "value" {
		t.Errorf("Expected value to be \"value\", got %q", value)
	}
}

func TestParseArg_FlagWithDoubleDashAndEquals(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--flag=value"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "value" {
		t.Errorf("Expected value to be \"value\", got %q", value)
	}
}

func TestParseArg_FlagWithSeparateValue(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag", "value"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "value" {
		t.Errorf("Expected value to be \"value\", got %q", value)
	}
}

func TestParseArg_FlagWithDoubleDashAndSeparateValue(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--flag", "value"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "value" {
		t.Errorf("Expected value to be \"value\", got %q", value)
	}
}

func TestParseArg_FlagWithoutValueWhenAnotherFlagFollows(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag1", "-flag2"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag1" {
		t.Errorf("Expected name to be \"flag1\", got %q", name)
	}
	if value != "true" {
		t.Errorf("Expected value to be \"true\", got %q", value)
	}
}

func TestParseArg_BooleanTrueFlagFormat(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-realtime", "true"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "realtime" {
		t.Errorf("Expected name to be \"realtime\", got %q", name)
	}
	if value != "true" {
		t.Errorf("Expected value to be \"true\", got %q", value)
	}
}

func TestParseArg_BooleanEqualsFormat(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-realtime=true"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "realtime" {
		t.Errorf("Expected name to be \"realtime\", got %q", name)
	}
	if value != "true" {
		t.Errorf("Expected value to be \"true\", got %q", value)
	}
}

func TestParseArg_IntegerFlag(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-port", "8168"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "port" {
		t.Errorf("Expected name to be \"port\", got %q", name)
	}
	if value != "8168" {
		t.Errorf("Expected value to be \"8168\", got %q", value)
	}
}

func TestParseArg_IntegerFlagWithEquals(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-port=8168"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "port" {
		t.Errorf("Expected name to be \"port\", got %q", name)
	}
	if value != "8168" {
		t.Errorf("Expected value to be \"8168\", got %q", value)
	}
}

func TestParseArg_IntegerFlagWithDoubleDash(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--port", "8168"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "port" {
		t.Errorf("Expected name to be \"port\", got %q", name)
	}
	if value != "8168" {
		t.Errorf("Expected value to be \"8168\", got %q", value)
	}
}

func TestParseArg_StringFlag(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-listen", "127.0.0.1"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "listen" {
		t.Errorf("Expected name to be \"listen\", got %q", name)
	}
	if value != "127.0.0.1" {
		t.Errorf("Expected value to be \"127.0.0.1\", got %q", value)
	}
}

func TestParseArg_StringFlagWithEquals(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-listen=127.0.0.1"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "listen" {
		t.Errorf("Expected name to be \"listen\", got %q", name)
	}
	if value != "127.0.0.1" {
		t.Errorf("Expected value to be \"127.0.0.1\", got %q", value)
	}
}

func TestParseArg_StringFlagWithDoubleDash(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--listen", "127.0.0.1"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "listen" {
		t.Errorf("Expected name to be \"listen\", got %q", name)
	}
	if value != "127.0.0.1" {
		t.Errorf("Expected value to be \"127.0.0.1\", got %q", value)
	}
}

func TestParseArg_DurationFlag(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-timeout", "30s"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "timeout" {
		t.Errorf("Expected name to be \"timeout\", got %q", name)
	}
	if value != "30s" {
		t.Errorf("Expected value to be \"30s\", got %q", value)
	}
}

func TestParseArg_DurationFlagWithEquals(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-timeout=30s"}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "timeout" {
		t.Errorf("Expected name to be \"timeout\", got %q", name)
	}
	if value != "30s" {
		t.Errorf("Expected value to be \"30s\", got %q", value)
	}
}

func TestParseArg_DurationFlagWithDoubleDash(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"--timeout", "30s"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "timeout" {
		t.Errorf("Expected name to be \"timeout\", got %q", name)
	}
	if value != "30s" {
		t.Errorf("Expected value to be \"30s\", got %q", value)
	}
}

func TestParseArg_ComplexDurationFlag(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-timeout", "1h30m45s"}, 0)

	if advance != 1 {
		t.Errorf("Expected advance to be 1, got %d", advance)
	}
	if name != "timeout" {
		t.Errorf("Expected name to be \"timeout\", got %q", name)
	}
	if value != "1h30m45s" {
		t.Errorf("Expected value to be \"1h30m45s\", got %q", value)
	}
}

func TestParseArg_FlagWithEmptyValue(t *testing.T) {
	advance, name, value := main.ParseArg([]string{"-flag="}, 0)

	if advance != 0 {
		t.Errorf("Expected advance to be 0, got %d", advance)
	}
	if name != "flag" {
		t.Errorf("Expected name to be \"flag\", got %q", name)
	}
	if value != "" {
		t.Errorf("Expected value to be empty, got %q", value)
	}
}
