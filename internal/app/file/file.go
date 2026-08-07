// Package file provides file management commands.
package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/utils"
	"github.com/spf13/cobra"
)

// (removed unused hashSize)

// NewCommand creates the file command group.
func NewCommand(fs interfaces.FileSystem, log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "File management tools",
		Long:  "Organize, analyze, and manipulate files and directories.",
	}

	cmd.AddCommand(newOrganizeCmd(fs, log))
	cmd.AddCommand(newStatsCmd(fs, log))
	cmd.AddCommand(newDuplicatesCmd(fs, log))
	cmd.AddCommand(newCleanCmd(fs, log))
	return cmd
}

func newOrganizeCmd(fs interfaces.FileSystem, log logger.Logger) *cobra.Command {
	var dryRun bool
	var by string

	cmd := &cobra.Command{
		Use:   "organize [directory]",
		Short: "Organize files by date or extension",
		Example: `  xef file organize ~/Downloads --by extension
  xef file organize ~/Downloads --by date --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			dir := args[0]
			if !fs.Exists(dir) {
				return fmt.Errorf("directory not found: %s", dir)
			}
			if !fs.IsDir(dir) {
				return fmt.Errorf("path is not a directory: %s", dir)
			}

			log.Info("organizing files", logger.String("directory", dir), logger.String("by", by))

			entries, err := fs.ListDir(dir)
			if err != nil {
				return fmt.Errorf("failed to list directory: %w", err)
			}

			moved := 0
			for _, entry := range entries {
				if entry.IsDir {
					continue
				}

				var targetDir string
				switch by {
				case "extension":
					ext := filepath.Ext(entry.Name)
					if ext == "" {
						targetDir = filepath.Join(dir, "no-extension")
					} else {
						targetDir = filepath.Join(dir, strings.TrimPrefix(ext, "."))
					}
				case "date":
					t := time.Unix(entry.ModTime, 0)
					targetDir = filepath.Join(dir, t.Format("2006-01"))
				default:
					return fmt.Errorf("unknown organize method: %s (use extension or date)", by)
				}

				if dryRun {
					log.Info("would move", logger.String("from", entry.Path), logger.String("to", targetDir))
					continue
				}

				if err := fs.MkdirAll(targetDir, 0o755); err != nil {
					log.Warn("failed to create directory", logger.String("path", targetDir), logger.Error(err))
					continue
				}

				newPath := filepath.Join(targetDir, entry.Name)
				if err := os.Rename(entry.Path, newPath); err != nil {
					log.Warn("failed to move file", logger.String("file", entry.Name), logger.Error(err))
					continue
				}
				moved++
			}

			if dryRun {
				tui.PrintInfo(fmt.Sprintf("Dry run complete. Would organize %d files in %s", len(entries), dir))
			} else {
				tui.PrintSuccess(fmt.Sprintf("Organized %d files in %s", moved, dir))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be done without moving files")
	cmd.Flags().StringVar(&by, "by", "extension", "organize by: extension, date")
	return cmd
}

func newStatsCmd(fs interfaces.FileSystem, _ logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "stats [path]",
		Short:   "Show file and directory statistics",
		Example: `  xef file stats ./my-project`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			path := args[0]
			if !fs.Exists(path) {
				return fmt.Errorf("path not found: %s", path)
			}

			var (
				totalFiles int64
				totalDirs  int64
				totalSize  int64
				extensions = make(map[string]int64)
				mu         sync.Mutex
			)

			start := time.Now()
			err := fs.WalkDir(path, func(_ string, info interfaces.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				mu.Lock()
				defer mu.Unlock()
				if info.IsDir {
					totalDirs++
				} else {
					totalFiles++
					totalSize += info.Size
					ext := filepath.Ext(info.Name)
					if ext != "" {
						extensions[strings.ToLower(ext)]++
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to walk directory: %w", err)
			}

			type extCount struct {
				Ext   string
				Count int64
			}
			var extList []extCount
			for ext, count := range extensions {
				extList = append(extList, extCount{Ext: ext, Count: count})
			}
			sort.Slice(extList, func(i, j int) bool {
				return extList[i].Count > extList[j].Count
			})

			tui.PrintTitle("File Statistics")
			fmt.Printf("Path:        %s\n", path)
			fmt.Printf("Total Files: %d\n", totalFiles)
			fmt.Printf("Total Dirs:  %d\n", totalDirs)
			fmt.Printf("Total Size:  %s\n", utils.FormatBytes(uint64(totalSize)))
			fmt.Printf("Duration:    %s\n", utils.FormatDuration(time.Since(start)))
			fmt.Println()

			if len(extList) > 0 {
				tui.PrintTitle("Top Extensions")
				var rows [][]string
				for i, ec := range extList {
					if i >= 10 {
						break
					}
					rows = append(rows, []string{ec.Ext, fmt.Sprintf("%d", ec.Count)})
				}
				fmt.Println(tui.Table([]string{"Extension", "Count"}, rows))
			}

			return nil
		},
	}
}

func newDuplicatesCmd(fs interfaces.FileSystem, log logger.Logger) *cobra.Command {
	var minSize int64
	var workers int

	cmd := &cobra.Command{
		Use:     "duplicates [directory]",
		Short:   "Find duplicate files by content hash",
		Example: `  xef file duplicates ~/Downloads --min-size 1024`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			dir := args[0]
			if !fs.Exists(dir) {
				return fmt.Errorf("directory not found: %s", dir)
			}
			if !fs.IsDir(dir) {
				return fmt.Errorf("path is not a directory: %s", dir)
			}

			log.Info("scanning for duplicates", logger.String("directory", dir))

			type fileHash struct {
				path string
				size int64
			}
			hashes := make(map[string][]fileHash)
			var mu sync.Mutex
			var wg sync.WaitGroup
			semaphore := make(chan struct{}, workers)

			start := time.Now()
			err := fs.WalkDir(dir, func(p string, info interfaces.FileInfo, err error) error {
				if err != nil || info.IsDir || info.Size < minSize {
					return nil
				}

				wg.Add(1)
				semaphore <- struct{}{}
				go func(path string, size int64) {
					defer wg.Done()
					defer func() { <-semaphore }()

					hash, err := hashFile(fs, path)
					if err != nil {
						log.Warn("failed to hash file", logger.String("file", path), logger.Error(err))
						return
					}
					mu.Lock()
					hashes[hash] = append(hashes[hash], fileHash{path: path, size: size})
					mu.Unlock()
				}(p, info.Size)

				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to walk directory: %w", err)
			}

			wg.Wait()

			var duplicateCount int
			var dupSize int64
			for _, files := range hashes {
				if len(files) > 1 {
					duplicateCount += len(files) - 1
					for _, f := range files[1:] {
						dupSize += f.size
					}
				}
			}

			tui.PrintTitle("Duplicate Files")
			fmt.Printf("Scanned in:  %s\n", utils.FormatDuration(time.Since(start)))
			fmt.Printf("Duplicates:  %d files\n", duplicateCount)
			fmt.Printf("Wasted:      %s\n", utils.FormatBytes(uint64(dupSize)))
			fmt.Println()

			for hash, files := range hashes {
				if len(files) < 2 {
					continue
				}
				displayHash := hash
				if len(displayHash) > 16 {
					displayHash = displayHash[:16]
				}
				fmt.Printf("Hash: %s... (%d files)\n", displayHash, len(files))
				for _, f := range files {
					fmt.Printf("  %s (%s)\n", f.path, utils.FormatBytes(uint64(f.size)))
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().Int64Var(&minSize, "min-size", 0, "minimum file size in bytes")
	cmd.Flags().IntVar(&workers, "workers", 50, "number of concurrent workers")
	return cmd
}

// hashFile computes the SHA-256 hash of a file using streaming to handle large files.
func hashFile(_ interfaces.FileSystem, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing file %s: %v", path, err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func newCleanCmd(fs interfaces.FileSystem, log logger.Logger) *cobra.Command {
	var patterns []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:     "clean [directory]",
		Short:   "Clean files matching patterns",
		Example: `  xef file clean ~/Downloads --pattern "*.tmp" --pattern "*.log"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			if !fs.Exists(dir) {
				return fmt.Errorf("directory not found: %s", dir)
			}
			if !fs.IsDir(dir) {
				return fmt.Errorf("path is not a directory: %s", dir)
			}

			dryRunEnabled, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return fmt.Errorf("failed to read dry-run flag: %w", err)
			}

			activePatterns := patterns
			if len(activePatterns) == 0 {
				activePatterns = []string{"*.tmp", "*.log", "*.bak", ".DS_Store", "Thumbs.db"}
			}

			var removed int
			var freed int64

			err = fs.WalkDir(dir, func(p string, info interfaces.FileInfo, err error) error {
				if err != nil || info.IsDir {
					return nil
				}

				for _, pattern := range activePatterns {
					matched, _ := filepath.Match(pattern, info.Name)
					if matched {
						if dryRunEnabled {
							log.Info("would remove", logger.String("file", p))
						} else {
							target := filepath.Join(dir, info.Name)
							if err := fs.Remove(filepath.Clean(target)); err == nil {
								removed++
								freed += info.Size
								log.Debug("removed", logger.String("file", target))
							}
						}
						break
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to clean directory: %w", err)
			}

			if dryRunEnabled {
				tui.PrintInfo(fmt.Sprintf("Dry run: would clean files matching %v in %s", activePatterns, dir))
			} else {
				tui.PrintSuccess(fmt.Sprintf("Removed %d files, freed %s", removed, utils.FormatBytes(uint64(freed))))
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&patterns, "pattern", nil, "file patterns to clean")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be removed")
	return cmd
}
