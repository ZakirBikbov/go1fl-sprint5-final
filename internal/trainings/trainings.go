package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return errors.New("the string must contain exactly two commas and be split into three parts")
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

	activityType := strings.TrimSpace(parts[1])
	if activityType == "" {
		return errors.New("empty activity type is not allowed")
	}

	rawDuration := parts[2]
	if strings.HasPrefix(rawDuration, " ") || strings.HasSuffix(rawDuration, " ") {
		return errors.New("duration value contains leading or trailing whitespace")
	}

	durationStr := strings.TrimSpace(rawDuration)

	if strings.Contains(durationStr, " ") {
		return errors.New("duration value contains invalid whitespace")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration: failed to parse: %w", err)
	}

	if duration <= 0 {
		return fmt.Errorf("duration must be a positive number, but received: %s", duration)
	}

	t.Steps = steps
	t.TrainingType = activityType
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	dist := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)
	hours := t.Duration.Hours()

	var (
		calories float64
		err      error
	)
	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба", "Спортходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", fmt.Errorf("unknown activity type: %s", t.TrainingType)
	}
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, hours, dist, speed, calories,
	)

	return result, nil
}
