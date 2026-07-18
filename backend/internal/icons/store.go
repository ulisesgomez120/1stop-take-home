// Package icons stores per-device uploaded icon files on disk, validating
// content by sniffing bytes rather than trusting client-supplied metadata.
package icons

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ErrUnsupportedType is returned when the sniffed content type isn't one of
// the allowed image formats.
var ErrUnsupportedType = errors.New("unsupported content type")

var extByContentType = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/webp": ".webp",
}

// Store writes icon uploads under a single directory.
type Store struct {
	dir string
}

// NewStore ensures dir exists and returns a Store rooted there.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create uploads dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Dir returns the directory icons are stored in, for mounting a static
// file handler.
func (s *Store) Dir() string {
	return s.dir
}

// Save sniffs r's content type, rejects anything outside the allowed image
// formats, and writes it under an opaque, server-generated filename. Any
// previously saved icon for deviceID is removed first so re-uploads
// overwrite rather than accumulate. It returns the URL path to serve the
// file from (e.g. "/icons/<file>"). Callers are responsible for bounding
// the size of r before calling Save.
func (s *Store) Save(deviceID string, r io.Reader) (string, error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read upload: %w", err)
	}
	head = head[:n]

	contentType := http.DetectContentType(head)
	ext, ok := extByContentType[contentType]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedType, contentType)
	}

	suffix := make([]byte, 8)
	if _, err := rand.Read(suffix); err != nil {
		return "", fmt.Errorf("generate filename: %w", err)
	}
	filename := fmt.Sprintf("%s-%s%s", deviceID, hex.EncodeToString(suffix), ext)

	if err := s.removeExisting(deviceID); err != nil {
		return "", fmt.Errorf("remove existing icon: %w", err)
	}

	dst, err := os.Create(filepath.Join(s.dir, filename))
	if err != nil {
		return "", fmt.Errorf("create icon file: %w", err)
	}
	defer dst.Close()

	if _, err := dst.Write(head); err != nil {
		return "", fmt.Errorf("write icon file: %w", err)
	}
	if _, err := io.Copy(dst, r); err != nil {
		return "", fmt.Errorf("write icon file: %w", err)
	}

	return "/icons/" + filename, nil
}

// removeExisting deletes any previously saved icon file for deviceID.
// deviceID is assumed to already be validated as a safe filename component
// by the caller.
func (s *Store) removeExisting(deviceID string) error {
	matches, err := filepath.Glob(filepath.Join(s.dir, deviceID+"-*"))
	if err != nil {
		return err
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
