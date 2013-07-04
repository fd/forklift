package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fd/forklift/root"
	"github.com/fd/go-cli/cli"
)

func init() {
	cli.Register(Update{})
}

type Update struct {
	root.Root
	cli.Arg0 `name:"update"`

	cli.Manual `
    Usage:   forklift update
    Summary: Update forklift.
  `
}

func (cmd *Update) Main() error {
	bin, err := find_bin(string(cmd.Root.Arg0))
	if err != nil {
		return err
	}

	fmt.Println("Looking for a new release:")
	release, err := get_latest_release(true)
	if err != nil {
		return err
	}

	asset, err := release.get_asset()
	if err != nil {
		return err
	}
	fmt.Printf(" - %s (%s)\n", release.Name, release.TagName)

	fmt.Printf(" - downloading ...")
	r, err := asset.load_bin()
	if err != nil {
		return err
	}
	fmt.Printf(" done\n")

	fmt.Printf(" - installing ...")
	f, err := os.OpenFile(bin, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	f.Close()

	os.Chtimes(bin, time.Now(), release.PublishedAt)
	fmt.Printf(" done\n")

	return nil
}

func find_bin(arg0 string) (string, error) {
	path, err := filepath.Abs(arg0)
	if err == nil {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	path, err = exec.LookPath(arg0)
	if err == nil {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Unable to determine location of the forklift binary")
}

type release_t struct {
	Draft       bool
	Prerelease  bool
	Name        string
	TagName     string    `json:"tag_name"`
	AssetsURL   string    `json:"assets_url"`
	PublishedAt time.Time `json:"published_at"`
}

type asset_t struct {
	Id      int
	Name    string
	release *release_t
}

func get_latest_release(prerelease bool) (*release_t, error) {
	var (
		releases []*release_t
		latest   *release_t
	)

	resp, err := http.Get("https://api.github.com/repos/fd/forklift/releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.Draft {
			continue
		}

		if release.Prerelease != prerelease {
			continue
		}

		if latest == nil {
			latest = release
			continue
		}

		if latest.PublishedAt.Before(release.PublishedAt) {
			latest = release
			continue
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("No releases were found")
	}

	return latest, nil
}

func (r *release_t) get_asset() (*asset_t, error) {
	var (
		assets   []*asset_t
		targeted *asset_t
		name     string
	)

	resp, err := http.Get(r.AssetsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	name = fmt.Sprintf("forklift-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	for _, asset := range assets {
		if asset.Name == name {
			targeted = asset
		}
	}

	if targeted == nil {
		return nil, fmt.Errorf("No releases found for %s %s", runtime.GOOS, runtime.GOARCH)
	}

	targeted.release = r
	return targeted, nil
}

func (a *asset_t) load_bin() (io.Reader, error) {
	var (
		buf  bytes.Buffer
		url  string
		name string
	)

	url = fmt.Sprintf("https://github.com/fd/forklift/releases/%s/%d/%s", a.release.TagName, a.Id, a.Name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	gzipr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	name = fmt.Sprintf("forklift-%s-%s/bin/forklift", runtime.GOOS, runtime.GOARCH)

	tr := tar.NewReader(gzipr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}

		if hdr.Name != name {
			continue
		}

		_, err = io.Copy(&buf, tr)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(buf.Bytes()), nil
	}

	return nil, fmt.Errorf("missing binary in archive: %s", name)
}
