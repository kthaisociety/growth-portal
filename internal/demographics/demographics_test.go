package demographics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kthaisociety/growth-portal/internal/lumaservice"
)

func answer(id, label, typ string, values ...string) lumaservice.Answer {
	return lumaservice.Answer{QuestionID: id, Label: label, Type: typ, Values: values}
}

func TestSummarize(t *testing.T) {
	members := []lumaservice.Member{
		{Answers: []lumaservice.Answer{
			answer("year", "Graduation year", "select", "2027"),
			answer("areas", "Areas of interest", "select", "NLP", "Robotics"),
			answer("bio", "About you", "text", "I like robots"),
		}},
		{Answers: []lumaservice.Answer{
			answer("year", "Graduation year", "select", "2028"),
			answer("areas", "Areas of interest", "select", "NLP"),
		}},
		{Answers: []lumaservice.Answer{
			answer("year", "Graduation year", "select", "2027"),
			answer("areas", "Areas of interest", "select"),
		}},
	}

	got := Summarize(members)

	assert.Equal(t, Summary{
		MemberCount: 3,
		Questions: []QuestionSummary{
			{QuestionID: "year", Label: "Graduation year", Answered: 3, Values: []ValueCount{
				{Value: "2027", Count: 2},
				{Value: "2028", Count: 1},
			}},
			{QuestionID: "areas", Label: "Areas of interest", Answered: 2, Values: []ValueCount{
				{Value: "NLP", Count: 2},
				{Value: "Robotics", Count: 1},
			}},
		},
	}, got)
}

func TestSummarizeNoMembers(t *testing.T) {
	assert.Equal(t, Summary{MemberCount: 0, Questions: []QuestionSummary{}}, Summarize(nil))
}
