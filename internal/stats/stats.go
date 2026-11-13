package stats

import (
	"github-api-client/internal/github"
	"sort"
)

type LanguageShare struct {
	Name    string
	Bytes   int
	Percent float64
}
type UserStats struct {
	Login          string
	Name           string
	Followers      int
	RepoCount      int
	TotalForks     int
	Languages      []LanguageShare
	ActivityByYear map[int]int
}

func Compute(user github.User, repos []github.Repo, langs []github.Languages) UserStats {
	return UserStats{
		Login:          user.Login,
		Name:           user.Name,
		Followers:      user.Followers,
		RepoCount:      len(repos),
		TotalForks:     totalForks(repos),
		Languages:      languageShares(langs),
		ActivityByYear: activityByYear(repos),
	}
}

func totalForks(repos []github.Repo) int {
	total := 0
	for _, r := range repos {
		total += r.ForksCount
	}
	return total
}

func languageShares(perRepo []github.Languages) []LanguageShare {
	bytesByLang := make(map[string]int)
	total := 0
	for _, langs := range perRepo {
		for name, n := range langs {
			bytesByLang[name] += n
			total += n
		}
	}

	shares := make([]LanguageShare, 0, len(bytesByLang))
	for name, n := range bytesByLang {
		shares = append(shares, LanguageShare{
			Name:    name,
			Bytes:   n,
			Percent: float64(n) * 100 / float64(total),
		})
	}

	sort.Slice(shares, func(i, j int) bool {
		if shares[i].Bytes != shares[j].Bytes {
			return shares[i].Bytes > shares[j].Bytes
		}
		return shares[i].Name < shares[j].Name
	})
	return shares
}

func activityByYear(repos []github.Repo) map[int]int {
	years := make(map[int]int)
	for _, r := range repos {
		years[r.CreatedAt.Year()]++
		years[r.UpdatedAt.Year()]++
	}
	return years
}
