package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

// Applicator defines the interface for the core application logic.
// This allows the CLI to be tested independently of the main app implementation.
type Applicator interface {
	Login(ctx context.Context, cfgPath string) error
	Wipe(ctx context.Context, cfgPath string) error
	SyncOpportunities(ctx context.Context, cfgPath string, fromDate, ifModifiedSince time.Time) error
	BatchUpdateOpportunityRefs(ctx context.Context, cfgPath, reference string, ids []string) error
}

// BuildCLI creates the full CLI command structure for the application.
// It injects the core application logic (the Applicator) into the command actions.
func BuildCLI(app Applicator) *cli.Command {
	// Define flags that are common across multiple commands.
	configFlag := &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.yaml",
		Usage:   "path to the configuration file",
	}

	agoFlag := &cli.StringFlag{
		Name:    "ago",
		Usage:   "only refresh records updated within this duration (e.g., '2h', '15m')",
		Aliases: []string{"a"},
	}

	sinceFlag := &cli.StringFlag{
		Name:    "since",
		Usage:   "only refresh records updated since this timestamp (format: '2006-01-02T15:04:05Z')",
		Aliases: []string{"s"},
	}

	fromDateFlag := &cli.StringFlag{
		Name:    "fromDate",
		Usage:   "start date for the date range to sync (format: '2006-01-02')",
		Aliases: []string{"f"},
	}

	// Define all application commands.
	loginCmd := &cli.Command{
		Name:  "login",
		Usage: "Authorize the application with your Salesforce account",
		Flags: []cli.Flag{configFlag},
		Action: func(ctx context.Context, c *cli.Command) error {
			return app.Login(ctx, c.String("config"))
		},
	}

	wipeCmd := &cli.Command{
		Name:  "wipe",
		Usage: "Delete the local token and database files for security",
		Flags: []cli.Flag{configFlag},
		Action: func(ctx context.Context, c *cli.Command) error {
			return app.Wipe(ctx, c.String("config"))
		},
	}

	opportunitiesCmd := &cli.Command{
		Name:    "opportunities",
		Usage:   "Fetch and save opportunities from Salesforce",
		Aliases: []string{"opps"},
		Flags:   []cli.Flag{configFlag, agoFlag, sinceFlag, fromDateFlag},
		Action: func(ctx context.Context, c *cli.Command) error {
			fromDate, ifModifiedSince, err := parseDateFlags(c.String("fromDate"), c.String("since"), c.String("ago"))
			if err != nil {
				return err
			}
			return app.SyncOpportunities(ctx, c.String("config"), fromDate, ifModifiedSince)
		},
	}

	updateRefsCmd := &cli.Command{
		Name:    "opportunities-ref",
		Usage:   "Batch update the Payout Reference for multiple opportunities",
		Aliases: []string{"oppsref"},
		Flags: []cli.Flag{
			configFlag,
			&cli.StringFlag{Name: "ref", Usage: "the new reference value to set", Required: true},
			&cli.StringFlag{Name: "ids", Usage: "a comma-separated list of Opportunity IDs to update", Required: true},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			ids := strings.Split(c.String("ids"), ",")
			return app.BatchUpdateOpportunityRefs(ctx, c.String("config"), c.String("ref"), ids)
		},
	}

	// Assemble the root command.
	rootCmd := &cli.Command{
		Name:     "sfcli",
		Usage:    "A CLI tool for interacting with the Salesforce API",
		Commands: []*cli.Command{loginCmd, wipeCmd, opportunitiesCmd, updateRefsCmd},
	}

	return rootCmd
}

// parseDateFlags processes the date-related flags and returns parsed time values.
// It enforces mutual exclusivity between --since and --ago.
func parseDateFlags(fromDateStr, sinceStr, agoStr string) (time.Time, time.Time, error) {
	var fromDate, ifModifiedSince time.Time
	var err error

	if fromDateStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid --fromDate format: %w", err)
		}
	}

	if sinceStr != "" && agoStr != "" {
		return time.Time{}, time.Time{}, fmt.Errorf("--since and --ago flags are mutually exclusive")
	}

	if sinceStr != "" {
		ifModifiedSince, err = time.Parse("2006-01-02T15:04:05Z", sinceStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid --since format: %w", err)
		}
	}

	if agoStr != "" {
		duration, err := time.ParseDuration(agoStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid --ago duration format: %w", err)
		}
		ifModifiedSince = time.Now().Add(-duration)
	}

	return fromDate, ifModifiedSince, nil
}
