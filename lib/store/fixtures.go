package store

import (
	"io/ioutil"
	"os"

	"code.uber.internal/infra/kraken/utils/testutil"

	"github.com/uber-go/tally"
)

// ConfigFixture returns a Config with initialized temporary directories.
func ConfigFixture() (Config, func()) {
	cleanup := &testutil.Cleanup{}
	defer cleanup.Recover()

	upload, err := ioutil.TempDir("/tmp", "upload")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(upload) })

	download, err := ioutil.TempDir("/tmp", "download")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(download) })

	cache, err := ioutil.TempDir("/tmp", "cache")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(cache) })

	config := Config{
		UploadDir:   upload,
		DownloadDir: download,
		CacheDir:    cache,
	}.applyDefaults()

	return config, cleanup.Run
}

// LocalFileStoreFixture returns a LocalFileStore using temp directories.
func LocalFileStoreFixture() (*LocalFileStore, func()) {
	var cleanup testutil.Cleanup
	defer cleanup.Recover()

	config, c := ConfigFixture()
	cleanup.Add(c)

	fs, err := NewLocalFileStore(config, tally.NewTestScope("", nil))
	if err != nil {
		panic(err)
	}
	cleanup.Add(fs.Close)

	return fs, cleanup.Run
}

// CAStoreConfigFixture returns config for CAStore for testing purposes.
func CAStoreConfigFixture() (CAStoreConfig, func()) {
	var cleanup testutil.Cleanup
	defer cleanup.Recover()

	upload, err := ioutil.TempDir("/tmp", "upload")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(upload) })

	cache, err := ioutil.TempDir("/tmp", "cache")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(cache) })

	return CAStoreConfig{
		UploadDir: upload,
		CacheDir:  cache,
	}, cleanup.Run
}

// CAStoreFixture returns a CAStore for testing purposes.
func CAStoreFixture() (*CAStore, func()) {
	var cleanup testutil.Cleanup
	defer cleanup.Recover()

	config, c := CAStoreConfigFixture()
	cleanup.Add(c)

	s, err := NewCAStore(config, tally.NoopScope)
	if err != nil {
		panic(err)
	}
	return s, cleanup.Run
}

// TorrentStoreFixture returns a TorrentStore for testing purposes.
func TorrentStoreFixture() (*TorrentStore, func()) {
	var cleanup testutil.Cleanup
	defer cleanup.Recover()

	download, err := ioutil.TempDir("/tmp", "download")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(download) })

	cache, err := ioutil.TempDir("/tmp", "cache")
	if err != nil {
		panic(err)
	}
	cleanup.Add(func() { os.RemoveAll(cache) })

	config := TorrentStoreConfig{
		DownloadDir: download,
		CacheDir:    cache,
	}
	s, err := NewTorrentStore(config, tally.NoopScope)
	if err != nil {
		panic(err)
	}
	return s, cleanup.Run
}
