package daysteps

import (
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
		return fmt.Errorf("ошибка: строка должна содержать ровно одну запятую и делиться на две части")
	}

	rawSteps := parts[0]
	rawDuration := parts[1]

	if strings.TrimSpace(rawSteps) != rawSteps {
		return fmt.Errorf("ошибка: значение steps содержит недопустимые пробелы")
	}
	if strings.TrimSpace(rawDuration) != rawDuration {
		return fmt.Errorf("ошибка: значение duration содержит недопустимые пробелы")
	}

	steps, err := strconv.Atoi(rawSteps)
	if err != nil {
		return fmt.Errorf("ошибка преобразования steps: %w", err)
	}

	if steps <= 0 {
		return fmt.Errorf("steps должен быть положительным числом, но получено: %d", steps)
	}

	duration, err := time.ParseDuration(rawDuration)
	if err != nil {
		return fmt.Errorf("ошибка преобразования duration: %w", err)
	}

	if duration <= 0 {
		return fmt.Errorf("duration должен быть положительным числом, но получено: %d", steps)
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

	var b strings.Builder
	fmt.Fprintf(&b, "Количество шагов: %d.\n", ds.Steps)
	fmt.Fprintf(&b, "Дистанция составила %.2f км.\n", distance)
	fmt.Fprintf(&b, "Вы сожгли %.2f ккал.\n", calories)

	return b.String(), nil
}
