package cache_test

import (
	. "github.com/andreasf/spotify-weekly-releases/cache"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
)

var _ = Describe("DiskCache", func() {
	// ioutil.WriteFile
	// ioutil.ReadFile
	// os.MkDirAll
	// os.Stat
	// os.IsNotExists
	var tempDir string

	BeforeEach(func() {
		var tempErr error
		tempDir, tempErr = ioutil.TempDir("", "test")
		Expect(tempErr).To(BeNil())
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDir)
		Expect(err).To(BeNil())
	})

	It("Stores data in a file below the cache root directory", func() {
		fooSha256 := "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"

		cache := NewDiskCache(tempDir)
		expectedPath := path.Join(tempDir, fooSha256[0:4], fooSha256[4:])

		err := cache.Set("foo", []byte("bar"))

		Expect(err).To(BeNil())

		contents, err := ioutil.ReadFile(expectedPath)
		Expect(err).To(BeNil())
		Expect(contents).ToNot(BeNil())
		Expect(contents).To(Equal([]byte("bar")))
	})

	It("Can retrieve data again", func() {
		cache := NewDiskCache(tempDir)
		err := cache.Set("foo", []byte("bar"))
		Expect(err).To(BeNil())

		err = cache.Set("bar", []byte("baz"))
		Expect(err).To(BeNil())

		err = cache.Set("baz", []byte("foo"))
		Expect(err).To(BeNil())

		data, err := cache.Get("foo")
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte("bar")))

		data, err = cache.Get("bar")
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte("baz")))

		data, err = cache.Get("baz")
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte("foo")))
	})

	It("Can delete entries", func() {
		cache := NewDiskCache(tempDir)
		err := cache.Set("foo", []byte("bar"))
		Expect(err).To(BeNil())

		data, err := cache.Get("foo")
		Expect(err).To(BeNil())
		Expect(data).To(Equal([]byte("bar")))

		err = cache.Delete("foo")
		Expect(err).To(BeNil())

		data, err = cache.Get("foo")
		Expect(err).ToNot(BeNil())
		Expect(data).To(BeNil())
	})
})
