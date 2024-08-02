package extensions

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func isValidUUID(u string) bool {
	if len(u) != 36 {
		return false
	}
	_, err := uuid.Parse(u)
	return err == nil
}

func readSingleFromZip(zipPath string, cb func(*zip.File) error) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		err = cb(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(ctx context.Context, u, target string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	req.Header.Set("User-Agent", "Mozilla/5.0 (AnytypeExtensionDownloader/1.0)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status code: %d", resp.StatusCode)
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// https://gist.github.com/paulerickson/6d8650947ee4e3f3dbcc28fde10eaae7
/**
 * Extract a zip file named source to directory destination.  Handles cases where destination dir…
 *  - does not exist (creates it)
 *  - is empty
 *  - already has source archive extracted into it (files are overwritten)
 *  - has other files in it, not in source archive (not overwritten)
 * But is expected to fail if it…
 *  - is not writable
 *  - contains a non-empty directory with the same path as a file in source archive (that's not a simple overwrite)
 */
func unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
