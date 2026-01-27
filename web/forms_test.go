package web

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func newRequest(t *testing.T, urlString string) *http.Request {
	t.Helper()
	r, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// TestSearchForm tests the SearchForm behaviour
func TestSearchForm(t *testing.T) {

	defaultDateFrom, defaultDateTo := defaultDateToAndFrom()

	tests := []struct {
		name           string
		inputURL       string
		searchForm     *SearchForm
		err            error      // top level errors
		validationErrs *Validator // validation errors
	}{
		{
			name:     "default",
			inputURL: "http://127.0.0.1:8080/invoices/",
			searchForm: &SearchForm{
				ReconciliationStatus: "NotReconciled",
				DateFrom:             defaultDateFrom,
				DateTo:               defaultDateTo,
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{},
			},
		},
		{
			name:     "defaults with page 2",
			inputURL: "http://127.0.0.1:8080/invoices/?date-from=2025-06-01&date-to=2025-05-01&search=search string&page=2",
			searchForm: &SearchForm{
				ReconciliationStatus: "NotReconciled",
				DateFrom:             time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
				DateTo:               time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 2,
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{
					"date-to": "End date cannot be before the start date.",
				},
			},
		},
		{
			name:     "all fields specified",
			inputURL: "http://127.0.0.1:8080/invoices/?status=NotReconciled&date-from=2025-06-01&date-to=2025-07-01&search=search string",
			searchForm: &SearchForm{
				ReconciliationStatus: "NotReconciled",
				DateFrom:             time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
				DateTo:               time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{},
			},
		},
		{
			name:     "default status",
			inputURL: "http://127.0.0.1:8080/invoices/?date-from=2025-06-01&date-to=2025-05-01&search=search string",
			searchForm: &SearchForm{
				ReconciliationStatus: "NotReconciled",
				DateFrom:             time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
				DateTo:               time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{
					"date-to": "End date cannot be before the start date.",
				},
			},
		},
		{
			name:     "invalid dateto",
			inputURL: "http://127.0.0.1:8080/invoices/?status=Reconciled&date-from=2025-06-01&date-to=2025-05-01&search=search string",
			searchForm: &SearchForm{
				ReconciliationStatus: "Reconciled",
				DateFrom:             time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
				DateTo:               time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{
					"date-to": "End date cannot be before the start date.",
				},
			},
		},
		{
			name:     "invalid datefrom",
			inputURL: "http://127.0.0.1:8080/invoices/?status=Reconciled&date-from=INVALID-06-01&date-to=2026-05-01&search=search string",
			searchForm: &SearchForm{
				ReconciliationStatus: "Reconciled",
				DateFrom:             time.Time{},
				DateTo:               time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{
					"date-from": "From date must be provided.",
				},
			},
		},
		{
			name:     "invalid status",
			inputURL: "http://127.0.0.1:8080/invoices/?status=XXXXX&date-from=2025-06-01&date-to=2025-07-01&search=search string",
			searchForm: &SearchForm{
				ReconciliationStatus: "XXXXX",
				DateFrom:             time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
				DateTo:               time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
				SearchString:         "search string",
				Page:                 1, // 1-based pagination.
			},
			err: nil,
			validationErrs: &Validator{
				Errors: map[string]string{
					"status": "Invalid status value provided.",
				},
			},
		},
	}

	for ii, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", ii, tt.name), func(t *testing.T) {
			simulatedRequest := newRequest(t, tt.inputURL)
			form := NewSearchForm()
			if err := DecodeURLParams(simulatedRequest, form); err != nil {
				if tt.err != err {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			validator := NewValidator()
			form.Validate(validator)

			if diff := cmp.Diff(form, tt.searchForm); diff != "" {
				t.Errorf("unexpected searchform diff %s", diff)
			}

			if diff := cmp.Diff(validator, tt.validationErrs); diff != "" {
				t.Errorf("unexpected validation diff %s", diff)
			}
		})
	}
}
