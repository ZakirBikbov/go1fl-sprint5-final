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
		return errors.New("ошибка: строка должна содержать ровно две запятые и делиться на три части")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return errors.New("ошибка: некорректное количество шагов")
	}

	activityType := strings.TrimSpace(parts[1])
	if activityType == "" {
		return errors.New("ошибка: пустой тип тренировки")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		return errors.New("ошибка: некорректная продолжительность")
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
		return "", fmt.Errorf("неизвестный тип тренировки: %s", t.TrainingType)
	}
	if err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Тип тренировки: %s\n", t.TrainingType)
	fmt.Fprintf(&b, "Длительность: %.2f ч.\n", hours)
	fmt.Fprintf(&b, "Дистанция: %.2f км.\n", dist)
	fmt.Fprintf(&b, "Скорость: %.2f км/ч\n", speed)
	fmt.Fprintf(&b, "Сожгли калорий: %.2f\n", calories)

	return b.String(), nil
}
