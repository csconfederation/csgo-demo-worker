package main

import (
	"errors"
	"strings"

	"github.com/csconfederation/demoScrape2/pkg/demoscrape2"
)

// classifyParseResult turns a ProcessDemo outcome into an HTTP status and body.
//
// The status is load-bearing for callers, not just for logging: CSC-Stats
// relays it and CSC-Core counts only >=500 toward its stats circuit breaker. A
// bad demo reported as 5xx therefore reads as "the parser is down" and
// suppresses unrelated uploads — which is how one unparseable recording dropped
// two of three maps from match 9077 (csconfederation/csc-issues#8). So an
// unprocessable demo must be a 4xx, and 5xx must mean we are genuinely broken.
func classifyParseResult(game *demoscrape2.Game, err error) (int, any) {
	if err == nil {
		return 200, game
	}

	switch {
	case strings.Contains(err.Error(), "ErrInvalidFileType"):
		// Not a demo at all — a malformed request rather than a bad recording.
		return 400, err.Error()

	case errors.Is(err, demoscrape2.ErrNoValidRounds):
		// Checked before ErrUnexpectedEndOfDemo: a truncated demo with no rounds
		// is both, and "nothing importable in here" is the actionable half.
		return 422, err.Error()

	case strings.Contains(err.Error(), "ErrUnexpectedEndOfDemo") && game != nil && game.Result == "Ended":
		// Demo stream was cut short but the match had already finished, so the
		// stats are complete. Not an error from the caller's point of view.
		return 200, game

	default:
		return 500, err.Error()
	}
}
