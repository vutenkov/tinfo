package main_test

import (
	"io/ioutil"
	"regexp"

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
		path     string
		jsonMode bool
		expected string
		result   string
	)

	JustBeforeEach(func() {
		app = TInfo{
			Path:     path,
			JSONMode: jsonMode,
		}

		result, _ = app.Run()
	})

	Context("when decoding a single-file torrent", func() {
		BeforeEach(func() {
			path = "./testdata/single-file-single-tracker.torrent"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
				expected = mustReadFile("./testdata/single-file-single-tracker.expected.txt")
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
				expected = mustReadFile("./testdata/single-file-single-tracker.expected.json")
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})

	Context("when decoding a multi-file torrent", func() {
		BeforeEach(func() {
			path = "./testdata/multi-files-single-trackers.torrent"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
				expected = mustReadFile("./testdata/multi-files-single-trackers.expected.txt")
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
				expected = mustReadFile("./testdata/multi-files-single-trackers.expected.json")
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})

	Context("when decoding a multi-tracker torrent", func() {
		BeforeEach(func() {
			path = "./testdata/single-file-multi-tracker.torrent"
		})

		Context("and text mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = false
				expected = mustReadFile("./testdata/single-file-multi-tracker.expected.txt")
			})

			It("should present the torrent as text", func() {
				Expect(result).To(ContainSubstring(expected))
			})
		})

		Context("and json mode is enabled", func() {
			BeforeEach(func() {
				jsonMode = true
				expected = mustReadFile("./testdata/single-file-multi-tracker.expected.json")
			})

			It("should present the torrent as JSON", func() {
				Expect(removeFilepath(result)).To(ContainSubstring(expected))
			})
		})
	})
})
