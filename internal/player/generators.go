package player

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

var generatorAliases = map[string]string{
	"phone": "phone", "tel": "phone", "телефон": "phone",
	"first_name": "first_name", "firstname": "first_name", "имя": "first_name", "name": "first_name",
	"last_name": "last_name", "lastname": "last_name", "surname": "last_name", "фамилия": "last_name",
	"patronymic": "patronymic", "middlename": "patronymic", "middle_name": "patronymic", "отчество": "patronymic",
	"address": "address", "адрес": "address",
	"inn": "inn", "инн": "inn",
	"bank_account": "bank_account", "account": "bank_account", "rs": "bank_account",
	"р/с": "bank_account", "расчетный счет": "bank_account", "расчётный счёт": "bank_account", "расчетныйсчет": "bank_account",
	"ogrnip": "ogrnip", "огрнип": "ogrnip",
}

var (
	maleFirst   = []string{"Иван", "Алексей", "Дмитрий", "Сергей", "Андрей", "Михаил", "Павел"}
	femaleFirst = []string{"Анна", "Мария", "Елена", "Ольга", "Наталья", "Татьяна", "Ирина"}
	maleLast    = []string{"Иванов", "Петров", "Смирнов", "Кузнецов", "Попов", "Соколов", "Лебедев"}
	femaleLast  = []string{"Иванова", "Петрова", "Смирнова", "Кузнецова", "Попова", "Соколова", "Лебедева"}
	malePatr    = []string{"Иванович", "Петрович", "Сергеевич", "Алексеевич", "Дмитриевич", "Андреевич"}
	femalePatr  = []string{"Ивановна", "Петровна", "Сергеевна", "Алексеевна", "Дмитриевна", "Андреевна"}
	streets     = []string{"Ленина", "Пушкина", "Гагарина", "Советская", "Мира", "Тверская", "Садовая"}
	cities      = []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск", "Екатеринбург"}
)

type personBundle struct {
	first, last, patronymic string
}

func NormalizeGeneratorName(raw string) (string, bool) {
	key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(raw), "ё", "е"))
	key = regexp.MustCompile(`\s+`).ReplaceAllString(key, " ")
	compact := strings.ReplaceAll(key, " ", "")
	if v, ok := generatorAliases[key]; ok {
		return v, true
	}
	if v, ok := generatorAliases[compact]; ok {
		return v, true
	}
	return "", false
}

func (c *RunContext) generateCanonical(canonical string) (string, error) {
	if v, ok := c.values[canonical]; ok {
		return v, nil
	}
	if c.person == nil {
		c.person = c.newPerson()
	}
	var value string
	switch canonical {
	case "phone":
		value = fmt.Sprintf("+79%09d", c.rng.Intn(1_000_000_000))
	case "first_name":
		value = c.person.first
	case "last_name":
		value = c.person.last
	case "patronymic":
		value = c.person.patronymic
	case "address":
		value = fmt.Sprintf("г. %s, ул. %s, д. %d, кв. %d",
			cities[c.rng.Intn(len(cities))],
			streets[c.rng.Intn(len(streets))],
			c.rng.Intn(120)+1,
			c.rng.Intn(200)+1,
		)
	case "inn":
		value = generateINN(c.rng)
	case "bank_account":
		value = "40817810" + randomDigits(c.rng, 12)
	case "ogrnip":
		value = generateOGRNIP(c.rng)
	default:
		return "", fmt.Errorf("unknown generator %q", canonical)
	}
	c.values[canonical] = value
	return value, nil
}

func (c *RunContext) newPerson() *personBundle {
	female := c.rng.Intn(2) == 0
	if female {
		return &personBundle{
			first:      femaleFirst[c.rng.Intn(len(femaleFirst))],
			last:       femaleLast[c.rng.Intn(len(femaleLast))],
			patronymic: femalePatr[c.rng.Intn(len(femalePatr))],
		}
	}
	return &personBundle{
		first:      maleFirst[c.rng.Intn(len(maleFirst))],
		last:       maleLast[c.rng.Intn(len(maleLast))],
		patronymic: malePatr[c.rng.Intn(len(malePatr))],
	}
}

func checksumDigit(digits []int, weights []int) int {
	total := 0
	for i, d := range digits {
		total += d * weights[i]
	}
	return total % 11 % 10
}

func generateINN(rng *rand.Rand) string {
	digits := make([]int, 10)
	for i := range digits {
		digits[i] = rng.Intn(10)
	}
	digits = append(digits, checksumDigit(digits, []int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8}))
	digits = append(digits, checksumDigit(digits, []int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8}))
	out := make([]byte, len(digits))
	for i, d := range digits {
		out[i] = byte('0' + d)
	}
	return string(out)
}

func generateOGRNIP(rng *rand.Rand) string {
	digits := make([]int, 15)
	digits[0] = 3
	for i := 1; i < 14; i++ {
		digits[i] = rng.Intn(10)
	}
	num := 0
	for i := 0; i < 14; i++ {
		num = num*10 + digits[i]
	}
	digits[14] = num % 13 % 10
	out := make([]byte, 15)
	for i, d := range digits {
		out[i] = byte('0' + d)
	}
	return string(out)
}
