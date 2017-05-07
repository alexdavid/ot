package ot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoOt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go OT Suite")
}
