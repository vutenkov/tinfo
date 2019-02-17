package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tinfo Suite")
}
