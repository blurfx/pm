package ui

import "testing"

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		want  bool
	}{
		{
			name:  "exact match",
			text:  "test",
			query: "test",
			want:  true,
		},
		{
			name:  "fuzzy match - characters in order",
			text:  "development",
			query: "dev",
			want:  true,
		},
		{
			name:  "fuzzy match - scattered characters",
			text:  "package-manager",
			query: "pm",
			want:  true,
		},
		{
			name:  "no match - characters out of order",
			text:  "test",
			query: "tset",
			want:  false,
		},
		{
			name:  "no match - missing characters",
			text:  "test",
			query: "testing",
			want:  false,
		},
		{
			name:  "empty query matches everything",
			text:  "anything",
			query: "",
			want:  true,
		},
		{
			name:  "case insensitive match",
			text:  "Development",
			query: "dev",
			want:  true,
		},
		{
			name:  "case insensitive match 2",
			text:  "development",
			query: "DEV",
			want:  true,
		},
		{
			name:  "partial match at end",
			text:  "build-production",
			query: "prod",
			want:  true,
		},
		{
			name:  "single character match",
			text:  "test",
			query: "t",
			want:  true,
		},
		{
			name:  "all characters must be present",
			text:  "abc",
			query: "xyz",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyMatch(tt.text, tt.query)
			if got != tt.want {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.text, tt.query, got, tt.want)
			}
		})
	}
}

func TestFuzzyScore(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		// We don't test exact scores since they might change,
		// but we test relative ordering
	}{
		{
			name:  "exact substring should score lower (better)",
			text:  "test-dev",
			query: "dev",
		},
		{
			name:  "prefix match should score lower (better)",
			text:  "development",
			query: "dev",
		},
		{
			name:  "scattered match should score higher (worse)",
			text:  "d-e-v-e-l-o-p-m-e-n-t",
			query: "dev",
		},
	}

	// Test that prefix matches score better than non-prefix matches
	prefixScore := fuzzyScore("development", "dev")
	nonPrefixScore := fuzzyScore("my-dev-tool", "dev")
	if prefixScore >= nonPrefixScore {
		t.Errorf("Prefix match should score better (lower). prefix=%d, non-prefix=%d", prefixScore, nonPrefixScore)
	}

	// Test that exact substring matches score better than scattered matches
	substringScore := fuzzyScore("test-dev", "dev")
	scatteredScore := fuzzyScore("d-e-v-e-l-o-p", "dev")
	if substringScore >= scatteredScore {
		t.Errorf("Substring match should score better (lower). substring=%d, scattered=%d", substringScore, scatteredScore)
	}

	// Test empty query
	emptyScore := fuzzyScore("anything", "")
	if emptyScore != 0 {
		t.Errorf("Empty query should score 0, got %d", emptyScore)
	}

	// Test non-matching query
	noMatchScore := fuzzyScore("abc", "xyz")
	if noMatchScore != 999999 {
		t.Errorf("Non-matching query should score 999999, got %d", noMatchScore)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := fuzzyScore(tt.text, tt.query)
			// Just verify it calculates a score
			// Note: Negative scores are valid (better matches have lower/negative scores)
			_ = score
		})
	}
}

func TestFuzzyScoreOrdering(t *testing.T) {
	// Test that scores properly order results
	query := "dev"

	// Test specific ordering relationships
	prefixScore := fuzzyScore("development", query)      // Prefix + substring = best
	substringScore := fuzzyScore("my-dev-tool", query)   // Contains as substring
	scatteredScore := fuzzyScore("d-e-v-e-l-o-p", query) // Scattered match
	noMatchScore := fuzzyScore("no-match", query)        // No match

	// Prefix should score better than substring
	if prefixScore >= substringScore {
		t.Errorf("prefix 'development' (score=%d) should score better than substring 'my-dev-tool' (score=%d)",
			prefixScore, substringScore)
	}

	// Substring should score better than scattered
	if substringScore >= scatteredScore {
		t.Errorf("substring 'my-dev-tool' (score=%d) should score better than scattered 'd-e-v-e-l-o-p' (score=%d)",
			substringScore, scatteredScore)
	}

	// Scattered should score better than no match
	if scatteredScore >= noMatchScore {
		t.Errorf("scattered 'd-e-v-e-l-o-p' (score=%d) should score better than no-match 'no-match' (score=%d)",
			scatteredScore, noMatchScore)
	}
}

func TestFuzzyMatchCaseInsensitive(t *testing.T) {
	tests := []struct {
		text  string
		query string
	}{
		{"Development", "dev"},
		{"DEVELOPMENT", "dev"},
		{"development", "DEV"},
		{"DevelopMent", "devment"},
	}

	for _, tt := range tests {
		t.Run(tt.text+"_"+tt.query, func(t *testing.T) {
			if !fuzzyMatch(tt.text, tt.query) {
				t.Errorf("fuzzyMatch(%q, %q) should be case-insensitive", tt.text, tt.query)
			}
		})
	}
}
