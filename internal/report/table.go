package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github-api-client/internal/stats"
)

func Print(w io.Writer, all []stats.UserStats) {
	logins := make([]string, len(all))
	for i, s := range all {
		logins[i] = s.Login
	}
	header := append([]string{"Metric"}, logins...)

	summary := [][]string{
		cells("Name", all, func(s stats.UserStats) string { return s.Name }),
		cells("Followers", all, func(s stats.UserStats) string { return fmt.Sprint(s.Followers) }),
		cells("Repositories", all, func(s stats.UserStats) string { return fmt.Sprint(s.RepoCount) }),
		cells("Total forks", all, func(s stats.UserStats) string { return fmt.Sprint(s.TotalForks) }),
	}
	table(w, "USER SUMMARY", header, summary)

	table(w, "LANGUAGE DISTRIBUTION (bytes, share of user's code)",
		append([]string{"Language"}, logins...), languageRows(all))

	table(w, "ACTIVITY BY YEAR (repo created/updated events)",
		append([]string{"Year"}, logins...), activityRows(all))
}

func languageRows(all []stats.UserStats) [][]string {
	totals := make(map[string]int)
	perUser := make([]map[string]stats.LanguageShare, len(all))
	for i, s := range all {
		perUser[i] = make(map[string]stats.LanguageShare)
		for _, l := range s.Languages {
			perUser[i][l.Name] = l
			totals[l.Name] += l.Bytes
		}
	}

	names := make([]string, 0, len(totals))
	for n := range totals {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if totals[names[i]] != totals[names[j]] {
			return totals[names[i]] > totals[names[j]]
		}
		return names[i] < names[j]
	})

	rows := make([][]string, 0, len(names))
	for _, n := range names {
		row := []string{n}
		for i := range all {
			if l, ok := perUser[i][n]; ok {
				row = append(row, fmt.Sprintf("%d (%.1f%%)", l.Bytes, l.Percent))
			} else {
				row = append(row, "-")
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func activityRows(all []stats.UserStats) [][]string {
	seen := make(map[int]bool)
	for _, s := range all {
		for y := range s.ActivityByYear {
			seen[y] = true
		}
	}
	years := make([]int, 0, len(seen))
	for y := range seen {
		years = append(years, y)
	}
	sort.Ints(years)

	rows := make([][]string, 0, len(years))
	for _, y := range years {
		row := []string{fmt.Sprint(y)}
		for _, s := range all {
			row = append(row, fmt.Sprint(s.ActivityByYear[y]))
		}
		rows = append(rows, row)
	}
	return rows
}

func cells(label string, all []stats.UserStats, f func(stats.UserStats) string) []string {
	row := []string{label}
	for _, s := range all {
		row = append(row, f(s))
	}
	return row
}

func table(w io.Writer, title string, header []string, rows [][]string) {
	fmt.Fprintf(w, "\n%s\n", title)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(header, "\t"))

	sep := make([]string, len(header))
	for i, h := range header {
		sep[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(tw, strings.Join(sep, "\t"))

	for _, r := range rows {
		fmt.Fprintln(tw, strings.Join(r, "\t"))
	}
	tw.Flush()
}
