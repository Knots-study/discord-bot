package _package

type Emojis struct {
	// New 🆕
	New string
	// Numbers "0⃣", "1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣", "🔟"
	Numbers []string
}

func InitEmojis() Emojis {
	emojis := Emojis{
		New:     "🆕",
		Numbers: []string{"0⃣", "1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣", "🔟"},
	}
	return emojis
}
