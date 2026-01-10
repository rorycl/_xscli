package main

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPagination(t *testing.T) {

	tests := []struct {
		name         string
		inputURL     string
		pageLen      int
		totalRecords int
		currentPage  int

		pagination  *Pagination
		nextURL     string
		previousURL string
		err         error
	}{
		{
			name:         "valid next and previous pages",
			inputURL:     "?status=ok&page=2&something=there",
			pageLen:      5,
			totalRecords: 13,
			currentPage:  2,
			pagination: &Pagination{
				PageNo:   2,
				Pages:    3,
				Next:     3,
				Previous: 1,
			},
			nextURL:     "?page=3&something=there&status=ok",
			previousURL: "?page=1&something=there&status=ok",
		},
		{
			name:         "same next and previous pages",
			inputURL:     "?status=ok&page=1&something=there",
			pageLen:      5,
			totalRecords: 5,
			currentPage:  1,
			pagination: &Pagination{
				PageNo:   1,
				Pages:    1,
				Next:     0,
				Previous: 0,
			},
			nextURL:     "",
			previousURL: "",
		},
		{
			name:         "invalid page length",
			inputURL:     "?status=ok&page=1&something=there",
			pageLen:      -5,
			totalRecords: 5,
			currentPage:  1,
			pagination: &Pagination{
				PageNo:   1,
				Pages:    5, // default pagelen 5
				Next:     2,
				Previous: 0,
			},
			nextURL:     "?page=2&something=there&status=ok",
			previousURL: "",
		},
		{
			name:         "invalid page number",
			inputURL:     "?status=ok&page=4&something=there",
			pageLen:      5,
			totalRecords: 14,
			currentPage:  4,
			err:          ErrInvalidPageNo{4, 3},
		},
	}

	for ii, tt := range tests {
		t.Run(fmt.Sprintf("%d_%s", ii, tt.name), func(t *testing.T) {

			parsedURL, err := url.Parse(tt.inputURL)
			if err != nil {
				t.Fatalf("could not parse inputURL: %v", err)
			}
			pg, err := NewPagination(tt.pageLen, tt.totalRecords, tt.currentPage, parsedURL.Query())
			if err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("could not parse inputURL: %v", err)
				}
				t.Log("err", err)
				return
			}
			if err == nil && tt.err != nil {
				t.Fatalf("expected error: %v", tt.err)
			}

			if diff := cmp.Diff(pg, tt.pagination, cmpopts.IgnoreUnexported(Pagination{})); diff != "" {
				t.Errorf("pagination struct error %v", diff)
			}

			if got, want := pg.NextURL(), tt.nextURL; got != want {
				t.Errorf("next url error:\ngot  %q\nwant %q", got, want)
			}
			if got, want := pg.PreviousURL(), tt.previousURL; got != want {
				t.Errorf("prev url error:\ngot  %q\nwant %q", got, want)
			}
		})
	}
}
