package main_test

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/format"

	. "github.com/vutenkov/tinfo"
)

func mustReadFile(path string) string {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(content)
}

func removeFilepath(json string) string {
	re := regexp.MustCompile("\\s*\"File\".+?,")
	return re.ReplaceAllString(json, "")
}

var _ = BeforeSuite(func() {
	TruncatedDiff = false
})

var _ = Describe("Tinfo", func() {
	var (
		app      TInfo
		basename string
		jsonMode bool
		expected string
		result   string
	)

	JustBeforeEach(func() {
		app = TInfo{
			Path:     fmt.Sprintf("./testdata/%s.torrent", basename),
			JSONMode: jsonMode,
		}

		result, _ = app.Run()

		suffix := map[bool]string{
			true:  "json",
			false: "txt",
		}

		expected = mustReadFile(
			fmt.Sprintf("./testdata/%s/%s.expected.%s", runtime.GOOS, basename, suffix[jsonMode]),
		)
	})

	Context("when decoding a single-file torrent", func() {
		BeforeEach(func() {
			basename = "single-file-single-tracker"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})

	Context("when decoding a multi-file torrent", func() {
		BeforeEach(func() {
			basename = "multi-files-single-trackers"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})

	Context("when decoding a multi-tracker torrent", func() {
		BeforeEach(func() {
			basename = "single-file-multi-tracker"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})
})
