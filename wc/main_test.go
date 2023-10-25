package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4\n")

	expected := 4

	got := count(b, false, false)

	if got != expected {
		t.Errorf("Expected %d, Got %d instead\n", expected, got)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\nline2\nline3 word1")

	expected := 3

	got := count(b, true, false)

	if got != expected {
		t.Errorf("Expected %d, Got %d instead\n", expected, got)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("word1\n word2 word3\n")
	expected := 15
	got := count(b, false, true)

	if got != expected {
		t.Errorf("Expected %d, Got %d instead\n", expected, got)

	}
}
