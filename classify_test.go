package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/csconfederation/demoScrape2/pkg/demoscrape2"
)

func TestClassifyParseResult(t *testing.T) {
	game := &demoscrape2.Game{}
	endedGame := &demoscrape2.Game{Result: "Ended"}

	tests := []struct {
		name       string
		game       *demoscrape2.Game
		err        error
		wantStatus int
	}{
		{
			name:       "success",
			game:       game,
			err:        nil,
			wantStatus: 200,
		},
		{
			// The match 9077 case: an aborted force-start recorded no rounds.
			// Must not be a 5xx — CSC-Core reads that as an outage and stops
			// uploading the rest of the series.
			name:       "no valid rounds is unprocessable",
			game:       game,
			err:        demoscrape2.ErrNoValidRounds,
			wantStatus: 422,
		},
		{
			// ProcessDemo joins ParseToEnd's error with ErrNoValidRounds, so the
			// sentinel arrives wrapped. errors.Is must still see through it.
			name:       "no valid rounds wrapped in a join still classifies",
			game:       game,
			err:        errors.Join(errors.New("demo stream ended unexpectedly (ErrUnexpectedEndOfDemo)"), demoscrape2.ErrNoValidRounds),
			wantStatus: 422,
		},
		{
			name:       "not a demo file",
			game:       game,
			err:        errors.New("invalid File-Type; expecting HL2DEMO in the first 8 bytes (ErrInvalidFileType)"),
			wantStatus: 400,
		},
		{
			// Truncated stream but the match had finished: stats are complete.
			name:       "truncated demo of a finished match still succeeds",
			game:       endedGame,
			err:        errors.New("demo stream ended unexpectedly (ErrUnexpectedEndOfDemo)"),
			wantStatus: 200,
		},
		{
			name:       "truncated demo of an unfinished match is a server error",
			game:       game,
			err:        errors.New("demo stream ended unexpectedly (ErrUnexpectedEndOfDemo)"),
			wantStatus: 500,
		},
		{
			name:       "unrecognised failure stays a server error",
			game:       game,
			err:        fmt.Errorf("something went wrong"),
			wantStatus: 500,
		},
		{
			// ProcessDemo is contracted to return a non-nil Game, but this path
			// dereferences it, so don't panic if that ever stops being true.
			name:       "nil game does not panic",
			game:       nil,
			err:        errors.New("demo stream ended unexpectedly (ErrUnexpectedEndOfDemo)"),
			wantStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := classifyParseResult(tt.game, tt.err)
			if status != tt.wantStatus {
				t.Errorf("classifyParseResult() status = %d, want %d", status, tt.wantStatus)
			}
		})
	}
}

// ErrNoValidRounds must never be classified as 5xx: that is exactly the
// misclassification that made a single bad demo look like a stats-API outage.
func TestNoValidRoundsIsNever5xx(t *testing.T) {
	for _, err := range []error{
		demoscrape2.ErrNoValidRounds,
		errors.Join(errors.New("demo stream ended unexpectedly (ErrUnexpectedEndOfDemo)"), demoscrape2.ErrNoValidRounds),
		fmt.Errorf("wrapped: %w", demoscrape2.ErrNoValidRounds),
	} {
		status, _ := classifyParseResult(&demoscrape2.Game{}, err)
		if status >= 500 {
			t.Errorf("classifyParseResult(%v) = %d, must not be 5xx", err, status)
		}
	}
}
