package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	separation := strings.Split(data, ",")
	if len(separation) != 3 {
		return 0, "", 0, fmt.Errorf("slice length is not equal to 3")
	}

	steps := separation[0]
	activity := separation[1]
	duration := separation[2]

	stepsTrain, err := strconv.Atoi(steps)
	if err != nil {
		return 0, "", 0, err
	}
	if stepsTrain <= 0 {
		return 0, "", 0, fmt.Errorf("steps must be greater than zero")
	}

	durationTrain, err := time.ParseDuration(duration)
	if err != nil {
		return 0, "", 0, err
	}
	if durationTrain <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be greater than zero")
	}

	return stepsTrain, activity, durationTrain, nil

}

func distance(steps int, height float64) float64 {
	stepLength := float64(height) * stepLengthCoefficient
	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	speedAverage := dist / hours
	return speedAverage
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	dist := distance(steps, height)
	durationAverage := meanSpeed(steps, height, duration)
	hours := duration.Hours()
	var calories float64

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", hours, dist, durationAverage, calories), nil
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", hours, dist, durationAverage, calories), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight")
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration")
	}
	durationAverage := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	caloriesRun := (weight * durationAverage * minutes) / minInH
	return caloriesRun, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight")
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration")
	}
	durationAverage := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * durationAverage * minutes) / minInH
	caloriesWalk := calories * walkingCaloriesCoefficient
	return caloriesWalk, nil
}
