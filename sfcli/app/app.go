package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sfcli/app/salesforce"
	"sfcli/config"
	"sfcli/db"
)

// App is the central orchestrator for the application's business logic.
// It coordinates interactions between configuration, the Salesforce API client, and the database.
type App struct{}

// New creates and returns a new App instance.
func New() *App {
	return &App{}
}

// Login orchestrates the OAuth2 login flow.
// It loads configuration and initiates the interactive authentication process.
func (a *App) Login(ctx context.Context, cfgPath string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	return salesforce.InitiateLogin(ctx, cfg.Salesforce.OAuth2Config, cfg.TokenFilePath)
}

// Wipe removes local data for security and confidentiality.
// It deletes the OAuth2 token file and the DuckDB database file.
func (a *App) Wipe(ctx context.Context, cfgPath string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	log.Printf("Deleting token file at: %s", cfg.TokenFilePath)
	if err := os.Remove(cfg.TokenFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	log.Printf("Deleting database file at: %s", cfg.DatabasePath)
	if err := os.Remove(cfg.DatabasePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete database file: %w", err)
	}

	log.Println("Wipe complete.")
	return nil
}

// SyncOpportunities fetches Opportunity records from Salesforce and persists them to the database.
// It handles filtering by a date range and incremental updates.
func (a *App) SyncOpportunities(ctx context.Context, cfgPath string, fromDate, ifModifiedSince time.Time) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if fromDate.IsZero() {
		fromDate = cfg.DateRangeStart
		log.Printf("No --fromDate specified, using default from config: %s", fromDate.Format("2006-01-02"))
	}

	sfClient, err := salesforce.NewClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create Salesforce client: %w", err)
	}
	log.Println("Salesforce client authenticated successfully.")

	dbConn, err := db.New(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbConn.Close()

	if err := dbConn.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	log.Println("Fetching Opportunities from Salesforce...")
	records, err := sfClient.GetOpportunities(ctx, fromDate, ifModifiedSince)
	if err != nil {
		return err
	}
	log.Printf("Fetched %d opportunities.", len(records))

	if err := dbConn.UpsertOpportunities(records); err != nil {
		return fmt.Errorf("failed to upsert opportunities: %w", err)
	}
	log.Println("Successfully upserted opportunities to database.")
	return nil
}

// BatchUpdateOpportunityRefs updates the Payout_Reference__c field for multiple opportunities.
// It uses the efficient Salesforce Composite Batch API.
func (a *App) BatchUpdateOpportunityRefs(ctx context.Context, cfgPath, reference string, ids []string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	sfClient, err := salesforce.NewClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create Salesforce client: %w", err)
	}

	log.Printf("Batch updating %d opportunities with reference '%s'...", len(ids), reference)
	err = sfClient.BatchUpdateOpportunityRefs(ctx, reference, ids)
	if err != nil {
		return err
	}

	log.Println("Batch update successful.")
	// Note: For a fully robust system, one would re-fetch these updated records
	// to update the local database. For this CLI, we assume a subsequent sync
	// will catch the changes.
	return nil
}
