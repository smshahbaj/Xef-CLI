// Package crypto provides cryptographic utility commands.
package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/xef/xefcli/internal/core/interfaces"
	"github.com/xef/xefcli/internal/core/logger"
	"github.com/xef/xefcli/internal/pkg/tui"
)

// charset constants for password generation.
const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// NewCommand creates the crypto command group.
func NewCommand(hasher interfaces.Hasher, log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crypto",
		Short: "Cryptographic tools",
		Long:  "Hashing, encoding, UUID generation, and password tools.",
	}

	cmd.AddCommand(newSHA256Cmd(hasher, log))
	cmd.AddCommand(newSHA512Cmd(hasher, log))
	cmd.AddCommand(newBCryptCmd(hasher, log))
	cmd.AddCommand(newUUIDCmd(log))
	cmd.AddCommand(newBase64Cmd(log))
	cmd.AddCommand(newPasswordCmd(log))
	return cmd
}

func newSHA256Cmd(hasher interfaces.Hasher, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "sha256 [file|string]",
		Short:   "Compute SHA-256 hash",
		Example: `  xef crypto sha256 document.pdf`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte
			var err error

			if _, statErr := os.Stat(args[0]); statErr == nil {
				data, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			} else {
				data = []byte(args[0])
			}

			hash := hasher.SHA256(data)
			fmt.Println(hash)
			return nil
		},
	}
}

func newSHA512Cmd(hasher interfaces.Hasher, log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:     "sha512 [file|string]",
		Short:   "Compute SHA-512 hash",
		Example: `  xef crypto sha512 "hello world"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte
			var err error

			if _, statErr := os.Stat(args[0]); statErr == nil {
				data, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			} else {
				data = []byte(args[0])
			}

			hash := hasher.SHA512(data)
			fmt.Println(hash)
			return nil
		},
	}
}

func newBCryptCmd(hasher interfaces.Hasher, log logger.Logger) *cobra.Command {
	var cost int
	var compare string

	cmd := &cobra.Command{
		Use:     "bcrypt [password]",
		Short:   "Hash or verify password with bcrypt",
		Example: `  xef crypto bcrypt mypassword --cost 12`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			password := args[0]

			if compare != "" {
				if err := hasher.CompareBCrypt(compare, password); err != nil {
					tui.PrintError("Password does not match")
					return fmt.Errorf("verification failed")
				}
				tui.PrintSuccess("Password matches")
				return nil
			}

			hash, err := hasher.BCrypt(password, cost)
			if err != nil {
				return fmt.Errorf("bcrypt failed: %w", err)
			}
			fmt.Println(hash)
			return nil
		},
	}

	cmd.Flags().IntVar(&cost, "cost", 10, "bcrypt cost factor (4-31)")
	cmd.Flags().StringVar(&compare, "compare", "", "hash to compare against")
	return cmd
}

func newUUIDCmd(log logger.Logger) *cobra.Command {
	var count int
	var upper bool

	cmd := &cobra.Command{
		Use:     "uuid",
		Short:   "Generate UUIDs",
		Example: `  xef crypto uuid --count 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for i := 0; i < count; i++ {
				id := uuid.New().String()
				if upper {
					id = strings.ToUpper(id)
				}
				fmt.Println(id)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&count, "count", 1, "number of UUIDs to generate")
	cmd.Flags().BoolVar(&upper, "upper", false, "uppercase output")
	return cmd
}

func newBase64Cmd(log logger.Logger) *cobra.Command {
	var decode bool

	cmd := &cobra.Command{
		Use:     "base64 [input]",
		Short:   "Base64 encode or decode",
		Example: `  xef crypto base64 "hello"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			if decode {
				data, err := base64.StdEncoding.DecodeString(input)
				if err != nil {
					return fmt.Errorf("decode failed: %w", err)
				}
				fmt.Println(string(data))
			} else {
				encoded := base64.StdEncoding.EncodeToString([]byte(input))
				fmt.Println(encoded)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&decode, "decode", false, "decode instead of encode")
	return cmd
}

func newPasswordCmd(log logger.Logger) *cobra.Command {
	var length int
	var noSpecial bool
	var noNumbers bool
	var noUpper bool

	cmd := &cobra.Command{
		Use:     "password",
		Short:   "Generate secure random password",
		Example: `  xef crypto password --length 32`,
		RunE: func(cmd *cobra.Command, args []string) error {
			password, err := generatePassword(length, !noUpper, !noNumbers, !noSpecial)
			if err != nil {
				return fmt.Errorf("failed to generate password: %w", err)
			}
			fmt.Println(password)
			return nil
		},
	}

	cmd.Flags().IntVar(&length, "length", 16, "password length (minimum 4)")
	cmd.Flags().BoolVar(&noSpecial, "no-special", false, "exclude special characters")
	cmd.Flags().BoolVar(&noNumbers, "no-numbers", false, "exclude numbers")
	cmd.Flags().BoolVar(&noUpper, "no-upper", false, "exclude uppercase letters")
	return cmd
}

// generatePassword creates a cryptographically secure random password.
func generatePassword(length int, upper, numbers, special bool) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	var chars string
	chars += lowerChars
	if upper {
		chars += upperChars
	}
	if numbers {
		chars += numberChars
	}
	if special {
		chars += specialChars
	}

	if upper == false && numbers == false && special == false {
		return "", fmt.Errorf("no character sets selected")
	}

	// Build password using crypto/rand for security
	result := make([]byte, length)
	charSetLen := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charSetLen)
		if err != nil {
			return "", fmt.Errorf("random generation failed: %w", err)
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}
