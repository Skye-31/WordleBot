package commands

import (
	"context"
	"fmt"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/uptrace/bun"
)

func stats(db *bun.DB, event *events.ApplicationCommandInteractionEvent) {
	user := event.User
	if u := event.SlashCommandInteractionData().Options.UserOption("user"); u != nil {
		user = u.User()
	}
	if user.ID != event.User.ID {
		mu := models.User{
			ID: user.ID,
		}
		if c, err := db.NewSelect().Model(&mu).WherePK().Where("public = ?", true).Count(context.TODO()); err != nil || c == 0 {
			_ = event.CreateMessage(discord.MessageCreate{Content: "This user's profile is not public.", Flags: 64})
			return
		}
	}
	mStats := models.UserStats{
		ID: user.ID,
	}
	if err := db.NewSelect().Model(&mStats).WherePK().Scan(context.TODO(), &mStats); err != nil {
		_ = event.CreateMessage(discord.MessageCreate{Content: "This user has no past stats.", Flags: 64})
		return
	}
	_ = event.CreateMessage(analyzeStats(mStats, user))
}

func analyzeStats(s models.UserStats, user *core.User) discord.MessageCreate {
	m := discord.NewMessageCreateBuilder()
	e := discord.NewEmbedBuilder().
		SetAuthor("Stats for "+user.Tag(), "", user.EffectiveAvatarURL(128)).
		SetColor(0x4fffff)

	total := append(append(append(append(s.Four, s.Five...), s.Six...), s.Seven...), s.Eight...)
	for i, v := range [][]int{total, s.Four, s.Five, s.Six, s.Seven, s.Eight} {
		if len(v) == 0 {
			continue
		}
		l := analyzeLength(v)
		e.AddField(getTitle(i), l.String(), true)
	}

	return m.AddEmbeds(e.Build()).Build()
}

func analyzeLength(a []int) lengthAnalysis {
	l := lengthAnalysis{Raw: a}
	if len(l.Raw) == 0 {
		return l
	}
	l.Total = 0
	l.ValMap = make(map[int]int)
	for i := range l.Raw {
		l.Total += l.Raw[i]
		l.ValMap[l.Raw[i]]++
	}
	l.Avg = floatToString(float64(l.Total) / float64(len(l.Raw)))
	if l.ValMap[0] > 0 {
		l.GaveUp = l.ValMap[0]
		delete(l.ValMap, 0)
	}
	return l
}

func floatToString(f float64) string {
	s := fmt.Sprintf("%.1f", f)
	if s[len(s)-1] == '0' {
		return s[:len(s)-2]
	}
	return s
}

type lengthAnalysis struct {
	Total  int
	Avg    string
	GaveUp int
	Raw    []int
	ValMap map[int]int
}

func (l *lengthAnalysis) String() string {
	s := fmt.Sprintf("%d Games Total\nAverage: %s\n", len(l.Raw), l.Avg)
	if l.GaveUp > 0 {
		s += fmt.Sprintf("Gave up: %d\n", l.GaveUp)
	}
	for i := 0; i < 9; i++ {
		if j, ok := l.ValMap[i]; ok {
			s += fmt.Sprintf("%d: %d\n", i, j)
		}

	}
	return s
}

func getTitle(i int) string {
	if i == 0 {
		return "Total"
	}
	return map[int]string{
		4: "4",
		5: "5",
		6: "6",
		7: "7",
		8: "8",
	}[i+3] + " letter words"
}
