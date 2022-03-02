package commands

import (
	"bytes"
	"context"
	"fmt"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/Skye-31/WordleBot/models"
	"github.com/fogleman/gg"
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
		SetImage("attachment://stats.png").
		SetColor(0x4fffff)

	total := append(append(append(append(s.Four, s.Five...), s.Six...), s.Seven...), s.Eight...)
	for i, v := range [][]int{s.Four, s.Five, s.Six, s.Seven, s.Eight} {
		if len(v) == 0 {
			continue
		}
		l := analyzeLength(v)
		e.AddField(getTitle(i), l.String(), true)
	}

	totalAnalysis := analyzeLength(total)
	dc := gg.NewContext(400, 400)
	dc.SetHexColor("#2f3136")
	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()

	dc.SetHexColor("#fff")
	fontFace, err := models.LoadFontFace(models.FontBytes, 20)
	if err != nil {
		return m.SetContent("Error generating graph").Build()
	}
	dc.SetFontFace(fontFace)
	dc.DrawStringAnchored("9 -", 40, 325, 1, 0.25)
	dc.DrawStringAnchored("8 -", 40, 293.75, 1, 0.25)
	dc.DrawStringAnchored("7 -", 40, 262.5, 1, 0.25)
	dc.DrawStringAnchored("6 -", 40, 231.25, 1, 0.25)
	dc.DrawStringAnchored("5 -", 40, 200, 1, 0.25)
	dc.DrawStringAnchored("4 -", 40, 168.75, 1, 0.25)
	dc.DrawStringAnchored("3 -", 40, 137.5, 1, 0.25)
	dc.DrawStringAnchored("2 -", 40, 106.25, 1, 0.25)
	dc.DrawStringAnchored("1 -", 40, 75, 1, 0.25)

	dc.DrawStringAnchored("Total", 350, 360, 0.5, 1)
	fontFace, err = models.LoadFontFace(models.FontBytes, 26)
	if err != nil {
		return m.SetContent("Error generating graph").Build()
	}
	dc.SetFontFace(fontFace)
	dc.DrawStringAnchored("# of Guesses Distribution", float64(dc.Width()/2), 25, 0.5, 0.5)

	dc.SetHexColor("#2596ce")
	dc.DrawRectangle(50, 50, 2, 300)
	dc.DrawRectangle(50, 350, 300, 2)

	t := float64(totalAnalysis.Total - totalAnalysis.ValMap[0])
	for i := 9; i > 0; i-- {
		v := float64(totalAnalysis.ValMap[i])
		dc.DrawRectangle(50, 32+float64(i)*31.5, v/t*300, 20)
	}
	dc.Fill()

	var b bytes.Buffer
	if err := dc.EncodePNG(&b); err != nil {
		return m.SetContent("Error generating graph").Build()
	}
	f := discord.NewFile("stats.png", &b)
	return m.AddFiles(f).AddEmbeds(e.Build()).Build()
}

func analyzeLength(a []int) lengthAnalysis {
	l := lengthAnalysis{}
	if len(a) == 0 {
		return l
	}
	l.Total = 0
	l.ValMap = make(map[int]int)
	for i := range a {
		l.Total += 1
		l.ValMap[a[i]]++
	}
	l.Avg = floatToString(float64(l.Total) / float64(len(a)))
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
	ValMap map[int]int
}

func (l *lengthAnalysis) String() string {
	s := fmt.Sprintf("%d Games Total\nAverage: %s\n", l.Total, l.Avg)
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
	return map[int]string{
		0: "4",
		1: "5",
		2: "6",
		3: "7",
		4: "8",
	}[i] + " letter words"
}
