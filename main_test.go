package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_name(t *testing.T) {
	g := NewWithT(t)
	g.Expect(true).To(BeTrue())
}
