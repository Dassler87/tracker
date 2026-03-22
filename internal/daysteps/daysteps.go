package daysteps

import (
	"errors"
	"fmt"
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
	// Разделяем строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректный формат данных: ожидалось 'шаги,длительность'")
	}

	// Преобразуем первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("не удалось распарсить количество шагов")
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// Проверяем, что количество шагов не аномально большое
	if steps >= 20000 {
		return 0, 0, errors.New("количество шагов должно быть меньше 20000")
	}

	// Преобразовываем второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("количество времени должно быть больше 0")
	}

	// Если всё прошло без ошибок, возвращаем количество шагов, продолжительность и nil (для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка при парсинге данных:", err)
		return ""
	}

	// Проверяем, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := distanceMeters / mInKm

	// Вычисляем количество калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	// Формируем строку для возврата
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKilometers, calories)
	return result
}
