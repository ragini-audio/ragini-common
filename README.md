# ragini-common

Shared Go modules for the Ragini ecosystem.

| Package | Contents |
|---|---|
| `pkg/schema` | Core data types: `Track`, `Profile`, `PlayEvent`, `SwipeEvent`, `Session`, `TagWeight`, `TrackTaxonomy` |
| `pkg/enrichment` | Tag processing helpers: `BuildTagSentence` (for embedding), `BuildTagsText` (for FTS5 indexing) |

## Usage

```go
import (
    "github.com/ragini-audio/ragini-common/pkg/schema"
    "github.com/ragini-audio/ragini-common/pkg/enrichment"
)
```

## Track struct changes (schema_version 2, 2026-03-29)

The `schema.Track` struct was updated to support the expanded taxonomy schema:

| Old field | New field | Notes |
|-----------|-----------|-------|
| `Embedding []float32` | `AudioEmbedding []float32` | Renamed; 384-dim |
| *(new)* | `LyricEmbedding []float32` | 384-dim lyric embedding |
| `Tags []string` | `Moods []TagWeight` | Join table; multi-label |
| `Genre *string` | `Genres []TagWeight` | Join table; multi-label |
| `HasLyrics bool` + `LyricText string` | `Lyrics *string` | nil=not fetched; ""=fetched/not found |
| `Wishlist bool` | *(removed)* | |
| `MBID *string` | *(removed)* | Not in current schema |
| *(new)* | `Instruments []TagWeight` | Join table |
| *(new)* | `LyricTags []TagWeight` | Join table; lyric themes, situational tags |
| *(new)* | `Keywords []string` | Free-form extracted keywords |
| *(new)* | `Characters []string` | Named characters from lyrics |
| *(new)* | `Description string` | LLM-generated prose description |
| *(new)* | `AnalysisTier int` | 0=none 1=fingerprint 2=acoustic 3=lyrics 4=LLM |
| `AnalysisTier *int` | `AnalysisTier int` | Was pointer, now value |

New types added:
- `TagWeight{Tag string, Weight float64, Source string}` — pairs a label with confidence and provenance (`"acoustic"|"llm"|"id3"|"manual"`)
- `Session{ID, ProfileID, TokenHash, CreatedAt, ExpiresAt}` — authenticated session record
- `TrackTaxonomy{Moods, Genres, Instruments, LyricTags []TagWeight; Keywords, Characters []string}` — all six join-table dimensions in one struct; used when passing full taxonomy sets between layers

`SwipeEvent.Timestamp` renamed to `SwipeEvent.RatedAt`. `SwipeEvent.Source string` added.

## Track struct changes (schema_version 3, 2026-03-30)

| New field | Notes |
|-----------|-------|
| `TagsText string` | Denormalized space-separated string of all taxonomy tag values; updated atomically after every join-table write; indexed by FTS5 porter stemming in both Ragini and Sargam |

## pkg/enrichment helpers

`BuildTagSentence(moods, genres, lyricTags, keywords []string) string`
Concatenates the four LLM-facing tag dimensions into a single embeddable sentence.
Instruments and characters are excluded: instruments describe sound (already captured by
`audio_embedding`); character names are proper nouns that embed poorly.
Used by Tier-4 to produce the tag-sentence vector that replaces the raw lyric embedding.

`BuildTagsText(moods, genres, instruments, lyricTags, keywords, characters []string) string`
Produces a space-separated string of all six tag dimensions for FTS5 indexing.
Used wherever `tags_text` is rebuilt after a taxonomy upsert.

## Licence

Business Source Licence 1.1. Copyright © 2026 ArthIQ Labs LLC.
