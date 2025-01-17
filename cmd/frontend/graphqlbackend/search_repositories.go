package graphqlbackend

import (
	"context"
	"math"
	"regexp"

	"github.com/sourcegraph/sourcegraph/pkg/api"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/pkg/search"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/pkg/search/query"
	searchbackend "github.com/sourcegraph/sourcegraph/pkg/search/backend"
)

var mockSearchRepositories func(args *search.Args) ([]searchResultResolver, *searchResultsCommon, error)
var repoIcon = "data:image/svg+xml;base64,PHN2ZyB2ZXJzaW9uPSIxLjEiIGlkPSJMYXllcl8xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB4PSIwcHgiIHk9IjBweCIKCSB2aWV3Qm94PSIwIDAgNjQgNjQiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDY0IDY0OyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+Cjx0aXRsZT5JY29ucyA0MDA8L3RpdGxlPgo8Zz4KCTxwYXRoIGQ9Ik0yMywyMi40YzEuMywwLDIuNC0xLjEsMi40LTIuNHMtMS4xLTIuNC0yLjQtMi40Yy0xLjMsMC0yLjQsMS4xLTIuNCwyLjRTMjEuNywyMi40LDIzLDIyLjR6Ii8+Cgk8cGF0aCBkPSJNMzUsMjYuNGMxLjMsMCwyLjQtMS4xLDIuNC0yLjRzLTEuMS0yLjQtMi40LTIuNHMtMi40LDEuMS0yLjQsMi40UzMzLjcsMjYuNCwzNSwyNi40eiIvPgoJPHBhdGggZD0iTTIzLDQyLjRjMS4zLDAsMi40LTEuMSwyLjQtMi40cy0xLjEtMi40LTIuNC0yLjRzLTIuNCwxLjEtMi40LDIuNFMyMS43LDQyLjQsMjMsNDIuNHoiLz4KCTxwYXRoIGQ9Ik01MCwxNmgtMS41Yy0wLjMsMC0wLjUsMC4yLTAuNSwwLjV2MzVjMCwwLjMtMC4yLDAuNS0wLjUsMC41aC0yN2MtMC41LDAtMS0wLjItMS40LTAuNmwtMC42LTAuNmMtMC4xLTAuMS0wLjEtMC4yLTAuMS0wLjQKCQljMC0wLjMsMC4yLTAuNSwwLjUtMC41SDQ0YzEuMSwwLDItMC45LDItMlYxMmMwLTEuMS0wLjktMi0yLTJIMTRjLTEuMSwwLTIsMC45LTIsMnYzNi4zYzAsMS4xLDAuNCwyLjEsMS4yLDIuOGwzLjEsMy4xCgkJYzEuMSwxLjEsMi43LDEuOCw0LjIsMS44SDUwYzEuMSwwLDItMC45LDItMlYxOEM1MiwxNi45LDUxLjEsMTYsNTAsMTZ6IE0xOSwyMGMwLTIuMiwxLjgtNCw0LTRjMS40LDAsMi44LDAuOCwzLjUsMgoJCWMxLjEsMS45LDAuNCw0LjMtMS41LDUuNFYzM2MxLTAuNiwyLjMtMC45LDQtMC45YzEsMCwyLTAuNSwyLjgtMS4zQzMyLjUsMzAsMzMsMjkuMSwzMywyOHYtMC42Yy0xLjItMC43LTItMi0yLTMuNQoJCWMwLTIuMiwxLjgtNCw0LTRjMi4yLDAsNCwxLjgsNCw0YzAsMS41LTAuOCwyLjctMiwzLjVoMGMtMC4xLDIuMS0wLjksNC40LTIuNSw2Yy0xLjYsMS42LTMuNCwyLjQtNS41LDIuNWMtMC44LDAtMS40LDAuMS0xLjksMC4zCgkJYy0wLjIsMC4xLTEsMC44LTEuMiwwLjlDMjYuNiwzOCwyNywzOC45LDI3LDQwYzAsMi4yLTEuOCw0LTQsNHMtNC0xLjgtNC00YzAtMS41LDAuOC0yLjcsMi0zLjRWMjMuNEMxOS44LDIyLjcsMTksMjEuNCwxOSwyMHoiLz4KPC9nPgo8L3N2Zz4K"

// searchRepositories searches for repositories by name.
//
// For a repository to match a query, the repository's name must match all of the repo: patterns AND the
// default patterns (i.e., the patterns that are not prefixed with any search field).
func searchRepositories(ctx context.Context, args *search.Args, limit int32) (res []searchResultResolver, common *searchResultsCommon, err error) {
	if mockSearchRepositories != nil {
		return mockSearchRepositories(args)
	}

	fieldWhitelist := map[string]struct{}{
		query.FieldRepo:        {},
		query.FieldRepoGroup:   {},
		query.FieldType:        {},
		query.FieldDefault:     {},
		query.FieldIndex:       {},
		query.FieldCount:       {},
		query.FieldMax:         {},
		query.FieldTimeout:     {},
		query.FieldFork:        {},
		query.FieldArchived:    {},
		query.FieldRepoHasFile: {},
	}
	// Don't return repo results if the search contains fields that aren't on the whitelist.
	// Matching repositories based whether they contain files at a certain path (etc.) is not yet implemented.
	for field := range args.Query.Fields {
		if _, ok := fieldWhitelist[field]; !ok {
			return nil, nil, nil
		}
	}

	pattern, err := regexp.Compile(args.Pattern.Pattern)
	if err != nil {
		return nil, nil, err
	}

	// Filter args.Repos by matching their names against the query pattern.
	common = &searchResultsCommon{}
	repos := make([]*search.RepositoryRevisions, 0, len(args.Repos))
	for _, r := range args.Repos {
		if pattern.MatchString(string(r.Repo.Name)) {
			repos = append(repos, r)
		}
	}

	// Filter the repos if there is a repohasfile: or -repohasfile field.
	if len(args.Pattern.FilePatternsReposMustExclude) > 0 || len(args.Pattern.FilePatternsReposMustInclude) > 0 {
		repos, err = reposToAdd(ctx, args.Zoekt, repos, args.Pattern)
		if err != nil {
			return nil, nil, err
		}
	}

	// Convert the repos to RepositoryResolvers.
	results := make([]searchResultResolver, 0, len(repos))
	for _, r := range repos {
		if len(results) == int(limit) {
			common.limitHit = true
			break
		}
		results = append(results, &RepositoryResolver{repo: r.Repo, icon: repoIcon})
	}

	return results, common, nil
}

// reposToAdd determines which repositories should be included in the result set based on whether they fit in the subset
// of repostiories specified in the query's `repohasfile` and `-repohasfile` fields if they exist.
func reposToAdd(ctx context.Context, zoekt *searchbackend.Zoekt, repos []*search.RepositoryRevisions, pattern *search.PatternInfo) ([]*search.RepositoryRevisions, error) {
	matchingIDs := make(map[api.RepoID]bool)
	if len(pattern.FilePatternsReposMustInclude) > 0 {
		for _, pattern := range pattern.FilePatternsReposMustInclude {
			// The high FileMatchLimit here is to make sure we get all the repo matches we can. Setting it to
			// len(repos) could mean we miss some repos since there could be for example len(repos) file matches in
			// the first repo and some more in other repos.
			p := search.PatternInfo{IsRegExp: true, FileMatchLimit: math.MaxInt32, IncludePatterns: []string{pattern}, PathPatternsAreRegExps: true, PathPatternsAreCaseSensitive: false, PatternMatchesContent: true, PatternMatchesPath: true}
			q, err := query.ParseAndCheck("file:" + pattern)
			if err != nil {
				return nil, err
			}
			newArgs := search.Args{Pattern: &p, Repos: repos, Query: q, UseFullDeadline: true, Zoekt: zoekt}
			matches, _, err := searchFilesInRepos(ctx, &newArgs)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				matchingIDs[m.repo.ID] = true
			}
		}
	} else {
		// Default to including all the repos, then excluding some of them below.
		for _, r := range repos {
			matchingIDs[r.Repo.ID] = true
		}
	}

	if len(pattern.FilePatternsReposMustExclude) > 0 {
		for _, pattern := range pattern.FilePatternsReposMustExclude {
			p := search.PatternInfo{IsRegExp: true, FileMatchLimit: math.MaxInt32, IncludePatterns: []string{pattern}, PathPatternsAreRegExps: true, PathPatternsAreCaseSensitive: false, PatternMatchesContent: true, PatternMatchesPath: true}
			q, err := query.ParseAndCheck("file:" + pattern)
			if err != nil {
				return nil, err
			}
			newArgs := search.Args{Pattern: &p, Repos: repos, Query: q, UseFullDeadline: true, Zoekt: zoekt}
			matches, _, err := searchFilesInRepos(ctx, &newArgs)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				matchingIDs[m.repo.ID] = false
			}
		}
	}

	var rsta []*search.RepositoryRevisions
	for _, r := range repos {
		if matchingIDs[r.Repo.ID] {
			rsta = append(rsta, r)
		}
	}

	return rsta, nil
}
