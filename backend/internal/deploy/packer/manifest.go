package packer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	ProjectName string            `json:"project_name"`
	Files       []ManifestFile    `json:"files"`
	HashTree    string            `json:"hash_tree"`
}

type ManifestFile struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

func (m *Manifest) AddFile(path string, size int64, checksum string) {
	m.Files = append(m.Files, ManifestFile{
		Path:     path,
		Size:     size,
		Checksum: checksum,
	})
}

func (m *Manifest) CalculateHashTree() {
	hashes := make([]string, len(m.Files))
	for i, f := range m.Files {
		h := sha256.Sum256([]byte(f.Path + f.Checksum))
		hashes[i] = hex.EncodeToString(h[:])
	}

	h := sha256.New()
	for _, hash := range hashes {
		h.Write([]byte(hash))
	}
	m.HashTree = hex.EncodeToString(h.Sum(nil))
}

func (m *Manifest) ToJSON() ([]byte, error) {
	m.CalculateHashTree()
	return json.MarshalIndent(m, "", "  ")
}

func (m *Manifest) Save(path string) error {
	data, err := m.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func CalculateFileChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func CalculateDirHash(sourceDir string) (string, error) {
	files := []string{}

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.Contains(path, ".git") {
			relPath, _ := filepath.Rel(sourceDir, path)
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	h := sha256.New()
	for _, f := range files {
		h.Write([]byte(f))
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}