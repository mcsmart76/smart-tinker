package bowling

// Score a game of bowling.
func Score(rolls []int) int {
	score := 0

	for _, roll := range rolls {
		if roll < 0 || roll > 10 {
			// Illegal number of pins.
			return -1
		}
	}

	for frame := 1; frame <= 10; frame++ {
		if len(rolls) < 2 {
			// Incomplete frame.
			return -1
		}
		if rolls[0] != 10 && (rolls[0] + rolls[1] > 10) {
			// Knocked down more than 10 pins in a frame.
			return -1
		}

		next_frame_index := -1
		if rolls[0] == 10 {  // Strike.
			if len(rolls) < 3 {
				// Incomplete frame.
				return -1
			}
			score += rolls[0] + rolls[1] + rolls[2]

			if frame == 10 {
				next_frame_index = 3
			} else {
				next_frame_index = 1
			}
		} else if rolls[0] + rolls[1] == 10 {  // Spare.
			if len(rolls) < 3 {
				// Incomplete frame.
				return -4
			}
			score += rolls[0] + rolls[1] + rolls[2]

			if frame == 10 {
				next_frame_index = 3
			} else {
				next_frame_index = 2
			}
		} else {
			score += rolls[0] + rolls[1]
			next_frame_index = 2
		}
		if next_frame_index < 1 || next_frame_index > 3 {
			panic("Invalid advancement of frame index")
		}
		rolls = rolls[next_frame_index:]
	}

	if len(rolls) != 0 {
		// Remaining rolls.
		return -1
	}

	return score
}
