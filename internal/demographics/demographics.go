// Package demographics turns Community tier members into aggregate counts. It's pure: no I/O, so it's
// tested with plain values.
package demographics

import (
	"cmp"
	"slices"

	"github.com/kthaisociety/growth-portal/internal/lumaservice"
)

// countedTypes are the question types summarised. Only choice questions: free-text answers could
// identify people, and sign-in is still simulated (docs/pcr/0006).
var countedTypes = map[string]bool{"select": true}

type Summary struct {
	MemberCount int               `json:"member_count"`
	Questions   []QuestionSummary `json:"questions"`
}

type QuestionSummary struct {
	QuestionID string `json:"question_id"`
	Label      string `json:"label"`
	// Answered is how many members answered. Multi-select answers count once per value, so the value
	// counts can add up to more than this.
	Answered int          `json:"answered"`
	Values   []ValueCount `json:"values"`
}

type ValueCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// Summarize counts answers per choice question. Questions keep the order they first appear in; values
// are sorted by count, most common first.
func Summarize(members []lumaservice.Member) Summary {
	type question struct {
		label    string
		answered int
		counts   map[string]int
	}
	var order []string
	questions := map[string]*question{}

	for _, m := range members {
		for _, a := range m.Answers {
			if !countedTypes[a.Type] || len(a.Values) == 0 {
				continue
			}
			q, ok := questions[a.QuestionID]
			if !ok {
				q = &question{label: a.Label, counts: map[string]int{}}
				questions[a.QuestionID] = q
				order = append(order, a.QuestionID)
			}
			q.answered++
			for _, v := range a.Values {
				q.counts[v]++
			}
		}
	}

	summary := Summary{MemberCount: len(members), Questions: []QuestionSummary{}}
	for _, id := range order {
		q := questions[id]
		values := make([]ValueCount, 0, len(q.counts))
		for v, n := range q.counts {
			values = append(values, ValueCount{Value: v, Count: n})
		}
		slices.SortFunc(values, func(a, b ValueCount) int {
			return cmp.Or(cmp.Compare(b.Count, a.Count), cmp.Compare(a.Value, b.Value))
		})
		summary.Questions = append(summary.Questions, QuestionSummary{
			QuestionID: id,
			Label:      q.label,
			Answered:   q.answered,
			Values:     values,
		})
	}
	return summary
}
