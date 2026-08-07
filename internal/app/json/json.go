// Package json provides JSON processing commands.
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/spf13/cobra"
)

// NewCommand creates the json command group.
func NewCommand(log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "JSON processing tools",
		Long:  "Format, validate, diff, and transform JSON data.",
	}

	cmd.AddCommand(newFormatCmd(log))
	cmd.AddCommand(newValidateCmd(log))
	cmd.AddCommand(newDiffCmd(log))
	return cmd
}

func newFormatCmd(_ logger.Logger) *cobra.Command {
	var indent string
	var compact bool

	cmd := &cobra.Command{
		Use:     "format [file]",
		Short:   "Format JSON file or stdin",
		Example: `  xef json format data.json --indent "  "`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var data []byte
			var err error

			if len(args) == 0 {
				data, err = readStdin()
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
			} else {
				data, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			}

			var obj interface{}
			if err = json.Unmarshal(data, &obj); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}

			var out []byte
			if compact {
				out, err = json.Marshal(obj)
			} else {
				out, err = json.MarshalIndent(obj, "", indent)
			}
			if err != nil {
				return fmt.Errorf("failed to format JSON: %w", err)
			}

			fmt.Println(string(out))
			return nil
		},
	}

	cmd.Flags().StringVar(&indent, "indent", "  ", "indentation string")
	cmd.Flags().BoolVar(&compact, "compact", false, "compact output")
	return cmd
}

func newValidateCmd(_ logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "validate [file]",
		Short:   "Validate JSON file or stdin",
		Example: `  xef json validate data.json`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var data []byte
			var err error

			if len(args) == 0 {
				data, err = readStdin()
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
			} else {
				data, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			}

			var obj interface{}
			if err := json.Unmarshal(data, &obj); err != nil {
				tui.PrintError(fmt.Sprintf("Invalid JSON: %v", err))
				return fmt.Errorf("validation failed")
			}

			tui.PrintSuccess("Valid JSON")
			return nil
		},
	}
}

func newDiffCmd(_ logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "diff [file1] [file2]",
		Short:   "Show differences between two JSON files",
		Example: `  xef json diff config.prod.json config.dev.json`,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			data1, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", args[0], err)
			}
			data2, err := os.ReadFile(args[1])
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", args[1], err)
			}

			var obj1, obj2 interface{}
			if err := json.Unmarshal(data1, &obj1); err != nil {
				return fmt.Errorf("invalid JSON in %s: %w", args[0], err)
			}
			if err := json.Unmarshal(data2, &obj2); err != nil {
				return fmt.Errorf("invalid JSON in %s: %w", args[1], err)
			}

			diffs := diffValue("", obj1, obj2)
			if len(diffs) == 0 {
				tui.PrintSuccess("Files are identical")
				return nil
			}

			tui.PrintTitle("Differences")
			for _, d := range diffs {
				fmt.Println(d)
			}
			return nil
		},
	}
}

// diffValue recursively compares two JSON values and returns differences.
func diffValue(path string, v1, v2 interface{}) []string {
	var diffs []string

	switch val1 := v1.(type) {
	case map[string]interface{}:
		val2, ok := v2.(map[string]interface{})
		if !ok {
			return append(diffs, fmt.Sprintf("~ %s: type mismatch (object != %T)", path, v2))
		}
		diffs = append(diffs, diffMaps(path, val1, val2)...)
	case []interface{}:
		val2, ok := v2.([]interface{})
		if !ok {
			return append(diffs, fmt.Sprintf("~ %s: type mismatch (array != %T)", path, v2))
		}
		diffs = append(diffs, diffArrays(path, val1, val2)...)
	default:
		if !bytes.Equal(mustJSON(v1), mustJSON(v2)) {
			diffs = append(diffs, fmt.Sprintf("~ %s: %v != %v", path, v1, v2))
		}
	}
	return diffs
}

func diffMaps(prefix string, m1, m2 map[string]interface{}) []string {
	var diffs []string
	allKeys := make(map[string]bool, len(m1)+len(m2))
	for k := range m1 {
		allKeys[k] = true
	}
	for k := range m2 {
		allKeys[k] = true
	}

	for k := range allKeys {
		path := prefix + "." + k
		if prefix == "" {
			path = k
		}

		v1, ok1 := m1[k]
		v2, ok2 := m2[k]

		switch {
		case !ok1:
			diffs = append(diffs, fmt.Sprintf("+ %s: %v", path, v2))
		case !ok2:
			diffs = append(diffs, fmt.Sprintf("- %s: %v", path, v1))
		default:
			diffs = append(diffs, diffValue(path, v1, v2)...)
		}
	}
	return diffs
}

func diffArrays(prefix string, a1, a2 []interface{}) []string {
	var diffs []string
	maxLen := len(a1)
	if len(a2) > maxLen {
		maxLen = len(a2)
	}

	for i := 0; i < maxLen; i++ {
		path := fmt.Sprintf("%s[%d]", prefix, i)
		switch {
		case i >= len(a1):
			diffs = append(diffs, fmt.Sprintf("+ %s: %v", path, a2[i]))
		case i >= len(a2):
			diffs = append(diffs, fmt.Sprintf("- %s: %v", path, a1[i]))
		default:
			diffs = append(diffs, diffValue(path, a1[i], a2[i])...)
		}
	}
	return diffs
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func readStdin() ([]byte, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 && stat.Mode()&os.ModeNamedPipe == 0 {
		return nil, fmt.Errorf("no data on stdin")
	}
	return os.ReadFile("/dev/stdin")
}
