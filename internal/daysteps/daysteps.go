package daysteps

import (
	"errors"
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
	// TODO: реализовать функцию
	// Разделяем строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат данных: ожидалось 'шаги,длительность'")
	}

	// Преобразуем первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов %d должно быть больше 0", steps)
	}

	// Преобразовываем второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга времени: %v", err)
	}

	// Проверяем, что длительность больше нуля.
	if duration <= 0 {
		return 0, 0, errors.New("количество времени должно быть больше 0")
	}

	// Если всё прошло без ошибок, возвращаем количество шагов, продолжительность и nil (для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	//  Если была ошибка ИЛИ количество шагов некорректно
	if err != nil || steps <= 0 {
		if err != nil {
			log.Println("Ошибка при обработке данных:", err)
		} else {
			fmt.Println("Некорректное количество шагов:", steps)
		}
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := distanceMeters / mInKm

	// Вычисляем количество калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка расчёта калорий: %v", err)
		return ""
	}

	// Формируем строку для возврата
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKilometers, calories)
	return result
}
