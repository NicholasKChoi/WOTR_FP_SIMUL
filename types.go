package main

type CharDiceNeededTable struct {
	Results []struct {
		NumCharDie int
		Frequency  int
		Sum        int
	}
}

type CorruptionInflictedTable struct {
	Results []struct {
		NumCorruptionInflicted int
		Frequency              int
		Sum                    int
	}
}

type RevealsTable struct {
	Results []struct {
		NumReveals int
		Frequency  int
		Sum        int
	}
}
