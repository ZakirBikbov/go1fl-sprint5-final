package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("the string must contain exactly one comma and be split into two parts")
	}

	rawSteps := parts[0]
	if strings.HasPrefix(rawSteps, " ") || strings.HasSuffix(rawSteps, " ") {
		return errors.New("steps value contains leading or trailing whitespace")
	}

	stepsStr := strings.TrimSpace(rawSteps)

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("invalid steps: failed to parse as integer: %w", err)
	}

	if steps <= 0 {
		return fmt.Errorf("steps must be a positive number, but received: %d", steps)
	}

	rawDuration := parts[1]
	if strings.HasPrefix(rawDuration, " ") || strings.HasSuffix(rawDuration, " ") {
		return errors.New("duration value contains leading or trailing whitespace")
	}

	durationStr := strings.TrimSpace(rawDuration)

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration: failed to parse: %w", err)
	}

	if duration <= 0 {
		return fmt.Errorf("duration must be a positive number, but received: %s", duration)
	}

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories,
	)

	return result, nil
}
