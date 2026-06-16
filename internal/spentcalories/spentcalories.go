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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: expected 3 elements, got %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to parse steps count: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps count must be greater than 0, got %d", steps)
	}

	activityType := parts[1]
	if activityType == "" {
		return 0, "", 0, fmt.Errorf("activity type cannot be empty")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("failed to parse duration: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be greater than 0, got %v", duration)
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	var stepLen float64

	if height > 0 {

		stepLen = height * stepLengthCoefficient
	} else {

		stepLen = lenStep
	}

	totalDistanceMeters := float64(steps) * stepLen

	return totalDistanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)
	durationHours := duration.Hours()
	return distKm / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		log.Printf("Invalid steps value: %d. Steps must be greater than 0", steps)
		return 0, fmt.Errorf("steps must be greater than 0, got %d", steps)
	}

	if weight <= 0 {
		log.Printf("Invalid weight value: %.2f. Weight must be greater than 0", weight)
		return 0, fmt.Errorf("weight must be greater than 0, got %.2f", weight)
	}

	if height <= 0 {
		log.Printf("Invalid height value: %.2f. Height must be greater than 0", height)
		return 0, fmt.Errorf("height must be greater than 0, got %.2f", height)
	}

	if duration <= 0 {
		log.Printf("Invalid duration value: %v. Duration must be greater than 0", duration)
		return 0, fmt.Errorf("duration must be greater than 0, got %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		log.Printf("Invalid steps value: %d. Steps must be greater than 0", steps)
		return 0, fmt.Errorf("steps must be greater than 0, got %d", steps)
	}

	if weight <= 0 {
		log.Printf("Invalid weight value: %.2f. Weight must be greater than 0", weight)
		return 0, fmt.Errorf("weight must be greater than 0, got %.2f", weight)
	}

	if height <= 0 {
		log.Printf("Invalid height value: %.2f. Height must be greater than 0", height)
		return 0, fmt.Errorf("height must be greater than 0, got %.2f", height)
	}

	if duration <= 0 {
		log.Printf("Invalid duration value: %v. Duration must be greater than 0", duration)
		return 0, fmt.Errorf("duration must be greater than 0, got %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	var errCalories error

	switch activityType {
	case "Ходьба":
		calories, errCalories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCalories = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	if errCalories != nil {
		log.Println(errCalories)
		return "", errCalories
	}

	durationHours := duration.Hours()

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, durationHours, distKm, speed, calories,
	)
	return result, nil
}
