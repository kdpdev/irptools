package irp

/*

/*
func ParseFrequency(str string) (Frequency, error) {
	freq, err := strconv.ParseUint(strings.TrimSpace(str), 10, 32)
	if err != nil {
		return 0, errs.Errorf("failed to parse frequency: %w", err)
	}

	if freq > 56000 {
		return 0, errs.Errorf("frequency is too large: %v", freq)
	}

	tie := freq % 1000
	if tie >= 500 {
		freq = freq - tie + 1000
	} else if tie > 0 {
		freq = freq - tie
	}

	return Frequency(freq), nil
}
*/

/*

func ValidateDurations(durations DurationsData) custom_error.CustomError {
	return checks.Catch(func() {
		checks.ThrowCheckPositiveInt("len(durations)", len(durations))
		checks.ThrowCheckZeroInt("len(durations) % 2", len(durations)%2)
		for i, d := range durations {
			if d == 0 {
				checks.Throw(custom_error.MakeErrorf("durations[%v]==0", i))
			}
		}
	})
}
*/
