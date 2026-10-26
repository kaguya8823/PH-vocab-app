package usecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const maxAudioSize int64 = 10 << 20 // 10MB

// AudioStorage は音声ファイルの保存を担当
type AudioStorage struct {
	baseDir string
}

// NewAudioStorage は保存先を指定し、保存機能を作成
func NewAudioStorage(baseDir string) (*AudioStorage, error) {
	baseDir = strings.TrimSpace(baseDir)
	if baseDir == "" {
		return nil, fmt.Errorf("audio storage directory is required")
	}

	absDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("resolve audio storage directory: %w", err)
	}

	if err := os.MkdirAll(absDir, 0750); err != nil {
		return nil, fmt.Errorf("create audio storage directory: %w", err)
	}

	return &AudioStorage{baseDir: absDir}, nil
}

// Save は音声データを新しいファイルとして保存
// 拡張子とMIMEタイプは呼び出しもとで検証済みの値を渡す
func (s *AudioStorage) Save(
	src io.Reader,
	extension string,
) (string, error) {
	switch extension {
	case ".webm", ".ogg", ".mp4", "wav":
	default:
		return "", fmt.Errorf("unsupported audio extension")
	}

	file, err := os.CreateTemp(s.baseDir, "audio-*"+extension)
	if err != nil {
		return "", fmt.Errorf("create audio file: %w", err)
	}

	filePath := file.Name()
	ok := false
	defer func() {
		_ = file.Close()
		if !ok {
			_ = os.Remove(filePath)
		}
	}()

	written, err := io.Copy(file, io.LimitReader(src, maxAudioSize+1))
	if err != nil {
		return "", fmt.Errorf("write audio file: %w", err)
	}
	if written > maxAudioSize {
		return "", fmt.Errorf("audio file exceeds 10 MiB limit")
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close audio file: %w", err)
	}

	ok = true
	return filePath, nil
}

// 削除処理
func (s *AudioStorage) Remove(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("resolve audio file path: %w", err)
	}

	relPath, err := filepath.Rel(s.baseDir, absPath)
	if err != nil {
		return fmt.Errorf("check audio file path: %w", err)
	}

	// 保存先ディレクトリの外にあるファイルは削除しない
	if relPath == "." ||
		relPath == ".." ||
		strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("audio fils is outside storage directory")
	}

	if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove audio file: %w", err)
	}

	return nil
}

// 録音→保存→再生　を可能にするAPI
func (s *AudioStorage) Open(filePath string) (*os.File, os.FileInfo, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve audio file path: %w", err)
	}

	relPath, err := filepath.Rel(s.baseDir, absPath)
	if err != nil {
		return nil, nil, fmt.Errorf("check audio file path: %w", err)
	}

	// 保存先ディレクトリの外にあるファイルは開かない
	if relPath == "." ||
		relPath == ".." ||
		strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
		return nil, nil, fmt.Errorf("audio file is outside storage directory")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open audio file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, nil, fmt.Errorf("read audio file information: %w", err)
	}

	if !info.Mode().IsRegular() {
		_ = file.Close()
		return nil, nil, fmt.Errorf("audio path is not a regular file")
	}

	return file, info, nil
}
