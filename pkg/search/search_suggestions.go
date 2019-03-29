package search

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/neelance/parallel"
	"github.com/pkg/errors"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/types"
	"github.com/sourcegraph/sourcegraph/pkg/api"
	"github.com/sourcegraph/sourcegraph/pkg/errcode"
	"github.com/sourcegraph/sourcegraph/pkg/search/query"
	"github.com/sourcegraph/sourcegraph/pkg/symbols"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

const (
	maxSearchSuggestions = 100
)

type SearchSuggestionsArgs struct {
	First *int32
}

func (a *SearchSuggestionsArgs) applyDefaultsAndConstraints() {
	if a.First == nil || *a.First < 0 || *a.First > maxSearchSuggestions {
		n := int32(maxSearchSuggestions)
		a.First = &n
	}
}

type SearchSuggestion struct {
	Repo      *types.Repo
	FileMatch *FileMatch
	Symbol    *symbols.Symbol

	// Score defines how well this item matches the query for sorting purposes
	Score int
	// Length holds the length of the item name as a second sorting criterium
	Length int
	// Label to sort alphabetically by when all else is equal.
	Label string
}

func (r *SearchSuggestion) ToFile() (*gitTreeEntryResolver, bool) {
	res, ok := r.result.(*gitTreeEntryResolver)
	return res, ok
}

func (r *SearchSuggestion) ToGitBlob() (*gitTreeEntryResolver, bool) {
	res, ok := r.result.(*gitTreeEntryResolver)
	return res, ok && res.stat.Mode().IsRegular()
}

func (r *SearchSuggestion) ToGitTree() (*gitTreeEntryResolver, bool) {
	res, ok := r.result.(*gitTreeEntryResolver)
	return res, ok && res.stat.Mode().IsDir()
}

func (r *SearchSuggestion) ToSymbol() (*symbolResolver, bool) {
	res, ok := r.result.(*symbolResolver)
	return res, ok
}

func (r *searcher) Suggestions(ctx context.Context, args *SearchSuggestionsArgs) ([]*SearchSuggestion, error) {
	args.applyDefaultsAndConstraints()

	if len(r.query.Syntax.Expr) == 0 {
		return nil, nil
	}

	// Only suggest for type:file.
	typeValues, _ := r.query.StringValues(query.FieldType)
	for _, resultType := range typeValues {
		if resultType != "file" {
			return nil, nil
		}
	}

	var suggesters []func(ctx context.Context) ([]*SearchSuggestion, error)

	showRepoSuggestions := func(ctx context.Context) ([]*SearchSuggestion, error) {
		// * If query contains only a single term (or 1 repogroup: token and a single term), treat it as a repo field here and ignore the other repo queries.
		// * If only repo fields (except 1 term in query), show repo suggestions.

		var effectiveRepoFieldValues []string
		if len(r.query.Values(query.FieldDefault)) == 1 && (len(r.query.Fields) == 1 || (len(r.query.Fields) == 2 && len(r.query.Values(query.FieldRepoGroup)) == 1)) {
			effectiveRepoFieldValues = append(effectiveRepoFieldValues, asString(r.query.Values(query.FieldDefault)[0]))
		} else if len(r.query.Values(query.FieldRepo)) > 0 && ((len(r.query.Values(query.FieldRepoGroup)) > 0 && len(r.query.Fields) == 2) || (len(r.query.Values(query.FieldRepoGroup)) == 0 && len(r.query.Fields) == 1)) {
			effectiveRepoFieldValues, _ = r.query.RegexpPatterns(query.FieldRepo)
		}

		// If we have a query which is not valid, just ignore it since this is for a suggestion.
		i := 0
		for _, v := range effectiveRepoFieldValues {
			if _, err := regexp.Compile(v); err == nil {
				effectiveRepoFieldValues[i] = v
				i++
			}
		}
		effectiveRepoFieldValues = effectiveRepoFieldValues[:i]

		if len(effectiveRepoFieldValues) == 0 {
			return nil, nil
		}

		_, _, repos, _, err := r.resolveRepositories(ctx, effectiveRepoFieldValues)
		return repos, err
	}

	suggesters = append(suggesters, showRepoSuggestions)

	showFileSuggestions := func(ctx context.Context) ([]*SearchSuggestion, error) {
		// If only repos/repogroups and files are specified (and at most 1 term), then show file
		// suggestions.  If the query has a single term, then consider it to be a `file:` filter (to
		// make it easy to jump to files by just typing in their name, not `file:<their name>`).
		hasOnlyEmptyRepoField := len(r.query.Values(query.FieldRepo)) > 0 && allEmptyStrings(r.query.RegexpPatterns(query.FieldRepo)) && len(r.query.Fields) == 1
		hasRepoOrFileFields := len(r.query.Values(query.FieldRepoGroup)) > 0 || len(r.query.Values(query.FieldRepo)) > 0 || len(r.query.Values(query.FieldFile)) > 0
		if !hasOnlyEmptyRepoField && hasRepoOrFileFields && len(r.query.Values(query.FieldDefault)) <= 1 {
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()
			return r.suggestFilePaths(ctx, maxSearchSuggestions)
		}
		return nil, nil
	}
	suggesters = append(suggesters, showFileSuggestions)

	showSymbolMatches := func(ctx context.Context) (results []*SearchSuggestion, err error) {

		repoRevs, _, _, _, err := r.resolveRepositories(ctx, nil)
		if err != nil {
			return nil, err
		}

		p, err := r.getPatternInfo(nil)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		fileMatches, _, err := searchSymbols(ctx, &Args{Pattern: p, Repos: repoRevs, Query: r.query}, 7)
		if err != nil {
			return nil, err
		}

		results = make([]*SearchSuggestion, 0)
		for _, fileMatch := range fileMatches {
			for _, sr := range fileMatch.Symbols {
				score := 20
				if sr.symbol.Parent == "" {
					score++
				}
				if len(sr.symbol.Name) < 12 {
					score++
				}
				switch symbols.CtagsKindToLSPSymbolKind(sr.symbol.Kind) {
				case lsp.SKFunction, lsp.SKMethod:
					score += 2
				case lsp.SKClass:
					score += 3
				}
				if len(sr.symbol.Name) >= 4 && strings.Contains(strings.ToLower(sr.uri.String()), strings.ToLower(sr.symbol.Name)) {
					score++
				}
				results = append(results, newSearchResults(sr, score))
			}
		}

		sortSearchSuggestions(results)
		const maxBoostedSymbolResults = 3
		boost := maxBoostedSymbolResults
		if len(results) < boost {
			boost = len(results)
		}
		if boost > 0 {
			for i := 0; i < boost; i++ {
				results[i].score += 200
			}
		}

		return results, nil
	}
	suggesters = append(suggesters, showSymbolMatches)

	showFilesWithTextMatches := func(ctx context.Context) ([]*SearchSuggestion, error) {
		// If terms are specified, then show files that have text matches. Set an aggressive timeout
		// to avoid delaying repo and file suggestions for too long.
		ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()
		if len(r.query.Values(query.FieldDefault)) > 0 {
			results, err := r.doResults(ctx, "file") // only "file" result type
			if err == context.DeadlineExceeded {
				err = nil // don't log as error below
			}
			var suggestions []*SearchSuggestion
			if results != nil {
				if len(results.results) > int(*args.First) {
					results.results = results.results[:*args.First]
				}
				for i, res := range results.results {
					file := res.fileMatch
					suggestions = append(suggestions, newSearchResult(file, len(results.results)-i))
				}
			}
			return suggestions, err
		}
		return nil, nil
	}
	suggesters = append(suggesters, showFilesWithTextMatches)

	// Run suggesters.
	var (
		allSuggestions []*SearchSuggestion
		mu             sync.Mutex
		par            = parallel.NewRun(len(suggesters))
	)
	for _, suggester := range suggesters {
		par.Acquire()
		go func(suggester func(ctx context.Context) ([]*SearchSuggestion, error)) {
			defer par.Release()
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			suggestions, err := suggester(ctx)
			if err == nil {
				mu.Lock()
				allSuggestions = append(allSuggestions, suggestions...)
				mu.Unlock()
			} else {
				if errors.Cause(err) == context.DeadlineExceeded || errors.Cause(err) == context.Canceled {
					log15.Warn("search suggestions exceeded deadline (skipping)", "query", r.rawQuery())
				} else if !errcode.IsBadRequest(err) {
					// We exclude bad user input. Note that this means that we
					// may have some tokens in the input that are valid, but
					// typing something "bad" results in no suggestions from the
					// this suggester. In future we should just ignore the bad
					// token.
					par.Error(err)
				}
			}
		}(suggester)
	}
	if err := par.Wait(); err != nil {
		if len(allSuggestions) == 0 {
			return nil, err
		}
		// If we got partial results, only log the error and return partial results
		log15.Error("error getting search suggestions: ", "error", err)
	}

	// Eliminate duplicates.
	type key struct {
		repoName api.RepoName
		repoRev  string
		file     string
		symbol   string
	}
	seen := make(map[key]struct{}, len(allSuggestions))
	uniqueSuggestions := allSuggestions[:0]
	for _, s := range allSuggestions {
		var k key
		switch s := s.result.(type) {
		case *types.Repo:
			k.repoName = s.Name
		case *FileMatch:
			k.repoName = s.Repo.Name
			// We explicitely do not use gitCommitResolver.OID() to get the OID here
			// because it could significantly slow down search suggestions from zoekt as
			// it doesn't specify the commit the default branch is on. This result would in
			// computing this commit for each suggestion, which could be heavy.
			k.repoRev = string(s.CommitID)
			// Zoekt only searches the default branch and sets commit ID to an empty string. This
			// may cause duplicate suggestions when merging results from Zoekt and non-Zoekt sources
			// (that do specify a commit ID), because their key k (i.e., k in seen[k]) will not
			// equal.
			k.file = s.Path
		case *symbols.Symbol:
			k.repoName = s.Repo.Name
			k.symbol = s.symbol.Name + s.symbol.Parent
		default:
			panic(fmt.Sprintf("unhandled: %#v", s))
		}

		if _, dup := seen[k]; !dup {
			uniqueSuggestions = append(uniqueSuggestions, s)
			seen[k] = struct{}{}
		}
	}
	allSuggestions = uniqueSuggestions

	sortSearchSuggestions(allSuggestions)
	if len(allSuggestions) > int(*args.First) {
		allSuggestions = allSuggestions[:*args.First]
	}

	return allSuggestions, nil
}

func allEmptyStrings(ss1, ss2 []string) bool {
	for _, s := range ss1 {
		if s != "" {
			return false
		}
	}
	for _, s := range ss2 {
		if s != "" {
			return false
		}
	}
	return true
}

func sortSearchSuggestions(s []*SearchSuggestion) {
	sort.Slice(s, func(i, j int) bool {
		// Sort by score
		a, b := s[i], s[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		// Prefer shorter strings for the same match score
		// E.g. prefer gorilla/mux over gorilla/muxy, Microsoft/vscode over g3ortega/vscode-crystal
		if a.Length != b.Length {
			return a.Length < b.Length
		}

		// All else equal, sort alphabetically.
		return a.Label < b.Label
	})
}
