// Copyright (c) 2026 ArthIQ Labs LLC. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

// Package schema defines the core SurrealDB data types for the Ragini ecosystem.
package schema

import "time"

// TrackSource describes the origin and confidence of a track record.
type TrackSource string

const (
	SourceLocalFLAC            TrackSource = "local_flac"
	SourceAirPlayKnown         TrackSource = "airplay_known"
	SourceAirPlayStream        TrackSource = "airplay_stream"
	SourceAirPlayFingerprinted TrackSource = "airplay_fingerprinted"
	SourceSargam               TrackSource = "sargam" // imported from a .sargam package
)

// SwipeSignal is one of the four swipe vocabulary signals.
type SwipeSignal string

const (
	SignalLove            SwipeSignal = "love"
	SignalDislike         SwipeSignal = "dislike"
	SignalSituationalSkip SwipeSignal = "situational_skip"
	SignalDeepen          SwipeSignal = "deepen"
)

// Track is a music track record in SurrealDB.
// ID is SHA-256(chromaprint_fingerprint) — stable across renames and re-encodes.
type Track struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Artist       string      `json:"artist"`
	Album        string      `json:"album"`
	Year         *int        `json:"year,omitempty"`
	DurationMS   int64       `json:"duration_ms"`
	FilePath     *string     `json:"file_path,omitempty"`
	Source       TrackSource `json:"source"`
	AnalysisTier *int        `json:"analysis_tier,omitempty"`
	Wishlist     bool        `json:"wishlist"`
	AnalysedAt   *time.Time  `json:"analysed_at,omitempty"`
	BPM          *float64    `json:"bpm,omitempty"` // beats per minute; set by Tier-2 analysis
	Key          *string     `json:"key,omitempty"` // musical key e.g. "C major"; set by Tier-2 analysis

	// Mood fields — populated by the MoodAnalyser during Tier-2 analysis.
	// Values are in the range [0.0, 1.0].
	// Valence:  0.0 = very negative/sad, 1.0 = very positive/happy (Russell Circumplex x-axis)
	// Arousal:  0.0 = very calm/sleepy,  1.0 = very excited/intense (Russell Circumplex y-axis)
	// Energy:   0.0 = very quiet/soft,   1.0 = very loud/powerful
	Valence *float32 `json:"valence,omitempty"`
	Arousal *float32 `json:"arousal,omitempty"`
	Energy  *float32 `json:"energy,omitempty"`

	// Tags are situational descriptors derived from mood values during Tier-2 analysis.
	// Examples: "energetic", "melancholy", "peaceful", "uplifting", "late-night", "focus".
	Tags []string `json:"tags,omitempty"`

	// Lyric fields — populated by Tier-3 analysis.
	HasLyrics bool   `json:"has_lyrics"`
	LyricText string `json:"lyric_text,omitempty"`
}

// Profile is a named user profile within a Ragini instance.
type Profile struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	AvatarColour string  `json:"avatar_colour"`
	PINHash      *string `json:"pin_hash,omitempty"` // bcrypt; never logged
	IsDefault    bool    `json:"is_default"`
}

// SwipeEvent records a single swipe gesture from a user.
type SwipeEvent struct {
	ID        string      `json:"id"`
	TrackID   string      `json:"track_id"`
	ProfileID string      `json:"profile_id"`
	Signal    SwipeSignal `json:"signal"`
	Timestamp time.Time   `json:"timestamp"`
}

// PlayEvent records a single playback session for a track.
// CompletedPct is 0.0 at start; updated toward 1.0 as playback progresses.
type PlayEvent struct {
	ID           string    `json:"id"`
	TrackID      string    `json:"track_id"`
	ProfileID    string    `json:"profile_id"`
	StartedAt    time.Time `json:"started_at"`
	CompletedPct float64   `json:"completed_pct"`
	SourcePath   string    `json:"source_path"`
}
