package spentcalories

import (
	"errors"
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
	// Разделяем строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных: ожидает 'шаги,тип_активности,длительность'")
	}

	// Преобразуем первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("ошибка преобразования количества шагов: %w")
	}
	// Проверка на отрицательное значение шагов
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов не может быть меньше 0")
	}

	// Получаем вид активности
	activity := parts[1]
	if activity == "" {
		return 0, "", 0, errors.New("тип активности не указан")
	}

	if activity != "Бег" && activity != "Ходьба" {
		return 0, "", 0, errors.New("неизвестный тип тренировки")
	}

	// Преобразовываем третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("ошибка преобразования продолжительности: " + err.Error())
	}

	// Проверка на нулевое и отрицательное время
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность активности должна быть больше нуля")
	}

	// Если всё прошло без ошибок, возвращаем количество шагов, вид активности, продолжительность и nil (для ошибки)
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// Умножаем количество шагов на длину шага
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := distanceMeters / mInKm

	return distanceKilometers
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверка всех входных параметров больше 0

	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	distanceKm := distance(steps, height)

	// Переводим продолжительность в часы
	durationHours := duration.Hours()

	// Вычисляем среднюю скорость
	meanSpeed := distanceKm / durationHours

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	// Парсинг входных данных
	steps, activity, durationHours, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, durationHours)
	var calories float64

	// Расчёт для разных видов тренировок
	switch activity {
	case "Ходьба":
		caloriesBurned, err := WalkingSpentCalories(steps, weight, height, durationHours)
		if err != nil {
			return "", err
		}
		calories = caloriesBurned

	case "Бег":
		caloriesBurned, err := RunningSpentCalories(steps, weight, height, durationHours)
		if err != nil {
			return "", err
		}
		calories = caloriesBurned

	default:
		return "", errors.New("неизвестный тип тренировки: " + activity)
	}

	// Формируем итоговую строку
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		activity, durationHours.Hours(), distance, speed, calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if weight <= 0 || weight > 300 || height <= 0 || height > 2.5 || steps < 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры для расчета калорий при беге")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Дополнительно проверяем среднюю скорость
	if meanSpeed >= 10 { // Пример максимального значения средней скорости
		return 0, errors.New("средняя скорость слишком высока")
	}

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if weight <= 0 || weight > 300 || height <= 0 || height > 2.5 {
		return 0, errors.New("некорректные входные параметры")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	// Умножаем на корректирующий коэффициент
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
