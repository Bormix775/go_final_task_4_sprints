package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		log.Printf("Invalid data format: expected 2 elements, got %d", len(parts))
		return 0, 0, fmt.Errorf("invalid data format: expected 2 elements, got %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("Error converting steps count: %v", err)
		return 0, 0, fmt.Errorf("error converting steps count: %v", err)
	}
	if steps <= 0 {
		log.Printf("Steps count must be greater than 0, got %d", steps)
		return 0, 0, fmt.Errorf("steps count must be greater than 0")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		log.Printf("Error parsing duration: %v", err)
		return 0, 0, fmt.Errorf("error parsing duration: %v", err)
	}
	if duration <= 0 {
		log.Printf("Duration must be greater than 0, got %v", duration)
		return 0, 0, fmt.Errorf("duration must be greater than 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
