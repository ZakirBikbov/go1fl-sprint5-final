package spentenergy

import (
	"fmt"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	var errs []string

	if steps <= 0 {
		errs = append(errs, "steps")
	}
	if weight <= 0 {
		errs = append(errs, "weight")
	}
	if height <= 0 {
		errs = append(errs, "height")
	}
	if duration <= 0 {
		errs = append(errs, "duration")
	}

	if len(errs) > 0 {
		return 0, fmt.Errorf("error: the following parameters must be greater than 0: %s", strings.Join(errs, ", "))
	}

	return (weight * MeanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return Distance(steps, height) / duration.Hours()
}

func Distance(steps int, height float64) float64 {
	return (float64(steps) * height * stepLengthCoefficient) / mInKm
}
