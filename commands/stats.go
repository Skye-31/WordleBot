package commands

import (
	"context"
	"fmt"
	"strconv"

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
	_ = event.CreateMessage(discord.NewMessageCreateBuilder().AddEmbeds(analyzeStats(mStats, user)).Build())
}

func analyzeStats(s models.UserStats, user *core.User) discord.Embed {
	e := discord.NewEmbedBuilder().
		SetAuthor("Stats for "+user.Tag(), "", user.EffectiveAvatarURL(128)).
		SetColor(0x4fffff)

	if h := handleArray(s.Four); h != "" {
		e.AddField("4 Letter Words", h, true)
	}
	if h := handleArray(s.Five); h != "" {
		e.AddField("5 Letter Words", h, true)
	}
	if h := handleArray(s.Six); h != "" {
		e.AddField("6 Letter Words", h, true)
	}
	if h := handleArray(s.Seven); h != "" {
		e.AddField("7 Letter Words", h, true)
	}
	if h := handleArray(s.Eight); h != "" {
		e.AddField("8 Letter Words", h, true)
	}

	return e.Build()
}

func handleArray(a []int) string {
	if len(a) == 0 {
		return ""
	}
	total := 0
	valMap := make(map[int]int)
	for i := range a {
		total += a[i]
		valMap[a[i]]++
	}
	avg := float64(total) / float64(len(a))
	s := "**Total**: " + strconv.Itoa(len(a)) + "\n" +
		"**Average**: " + floatToString(avg) + "\n"
	if valMap[0] > 0 {
		s += "**Gave up**: " + strconv.Itoa(valMap[0]) + "\n"
	}
	for i := 1; i < 9; i++ {
		if valMap[i] > 0 {
			s += "**" + strconv.Itoa(i) + "**: " + strconv.Itoa(valMap[i]) + "\n"
		}
	}

	return s
}

func floatToString(f float64) string {
	s := fmt.Sprintf("%.1f", f)
	if s[len(s)-1] == '0' {
		return s[:len(s)-2]
	}
	return s

}
