// Package http contains HTTP client commands and helpers.
package http

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/smshahbaj/Xef-CLI/internal/core/interfaces"
	"github.com/smshahbaj/Xef-CLI/internal/core/logger"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/tui"
	"github.com/smshahbaj/Xef-CLI/internal/pkg/utils"
	"github.com/spf13/cobra"
)

// NewCommand creates the http command group.
func NewCommand(client interfaces.HTTPClient, log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP client tools",
		Long:  "GET, POST, download, and benchmark HTTP endpoints.",
	}

	cmd.AddCommand(newGetCmd(client, log))
	cmd.AddCommand(newDownloadCmd(client, log))
	cmd.AddCommand(newBenchmarkCmd(client, log))
	return cmd
}

func newGetCmd(client interfaces.HTTPClient, log logger.Logger) *cobra.Command {
	var headers []string
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:     "get [url]",
		Short:   "Perform HTTP GET request",
		Example: `  xef http get https://api.github.com/users/octocat`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h := parseHeaders(headers)
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()

			start := time.Now()
			resp, err := client.Get(ctx, args[0], h)
			if err != nil {
				return fmt.Errorf("request failed: %w", err)
			}

			fmt.Println(string(resp.Body))
			log.Info("request completed",
				logger.String("url", args[0]),
				logger.Int("status", resp.StatusCode),
				logger.String("duration", utils.FormatDuration(time.Since(start))),
			)
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&headers, "header", nil, "request headers (Key:Value)")
	cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "request timeout")
	return cmd
}

func newDownloadCmd(client interfaces.HTTPClient, log logger.Logger) *cobra.Command {
	var output string
	var headers []string

	cmd := &cobra.Command{
		Use:     "download [url]",
		Short:   "Download file from URL",
		Example: `  xef http download https://example.com/file.zip -o output.zip`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			if output == "" {
				parts := strings.Split(url, "/")
				output = parts[len(parts)-1]
				if output == "" {
					output = "download"
				}
			}

			h := parseHeaders(headers)
			ctx := cmd.Context()

			log.Info("downloading", logger.String("url", url), logger.String("output", output))
			if err := client.Download(ctx, url, output, h); err != nil {
				return fmt.Errorf("download failed: %w", err)
			}

			tui.PrintSuccess(fmt.Sprintf("Downloaded to %s", output))
			return nil
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "output file path")
	cmd.Flags().StringSliceVar(&headers, "header", nil, "request headers")
	return cmd
}

func newBenchmarkCmd(client interfaces.HTTPClient, _ logger.Logger) *cobra.Command {
	var requests int
	var concurrency int
	var timeout time.Duration

	cmd := &cobra.Command{
		Use:     "benchmark [url]",
		Short:   "Benchmark HTTP endpoint",
		Example: `  xef http benchmark https://api.example.com -n 1000 -c 50`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			fmt.Printf("Benchmarking %s\n", url)
			fmt.Printf("Requests: %d, Concurrency: %d\n\n", requests, concurrency)

			var wg sync.WaitGroup
			semaphore := make(chan struct{}, concurrency)

			type result struct {
				err    error
				status int
				dur    time.Duration
			}

			results := make(chan result, requests)
			start := time.Now()

			for i := 0; i < requests; i++ {
				wg.Add(1)
				semaphore <- struct{}{}
				go func() {
					defer wg.Done()
					defer func() { <-semaphore }()

					ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
					defer cancel()

					rStart := time.Now()
					resp, err := client.Get(ctx, url, nil)
					res := result{dur: time.Since(rStart)}
					if err != nil {
						res.err = err
					} else {
						res.status = resp.StatusCode
					}
					results <- res
				}()
			}

			wg.Wait()
			close(results)

			total := time.Since(start)
			var success, failed int
			var totalDur time.Duration
			var minDur, maxDur time.Duration

			for r := range results {
				if r.err != nil {
					failed++
					continue
				}
				success++
				totalDur += r.dur
				if minDur == 0 || r.dur < minDur {
					minDur = r.dur
				}
				if r.dur > maxDur {
					maxDur = r.dur
				}
			}

			avg := time.Duration(0)
			if success > 0 {
				avg = totalDur / time.Duration(success)
			}

			tui.PrintTitle("Benchmark Results")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Metric", "Value"})
			table.Append([]string{"Total Requests", fmt.Sprintf("%d", requests)})
			table.Append([]string{"Success", fmt.Sprintf("%d", success)})
			table.Append([]string{"Failed", fmt.Sprintf("%d", failed)})
			table.Append([]string{"Total Time", utils.FormatDuration(total)})
			table.Append([]string{"Average", utils.FormatDuration(avg)})
			table.Append([]string{"Min", utils.FormatDuration(minDur)})
			table.Append([]string{"Max", utils.FormatDuration(maxDur)})
			if total.Seconds() > 0 {
				table.Append([]string{"RPS", fmt.Sprintf("%.2f", float64(requests)/total.Seconds())})
			}
			table.Render()

			return nil
		},
	}

	cmd.Flags().IntVar(&requests, "requests", 100, "number of requests")
	cmd.Flags().IntVar(&concurrency, "concurrency", 10, "concurrent requests")
	cmd.Flags().DurationVar(&timeout, "timeout", 10*time.Second, "request timeout")
	return cmd
}

func parseHeaders(headers []string) map[string]string {
	h := make(map[string]string)
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			h[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return h
}
