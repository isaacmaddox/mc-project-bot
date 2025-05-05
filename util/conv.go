package util

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	AMOUNT_SHULKER = 64 * 27
	AMOUNT_STACK   = 64
)

func From_unit(measure string) int {
	divided := strings.Split(measure, " ")

	if len(divided) == 0 {
		return 0
	}

	amount := Extract(strconv.Atoi(strings.Trim(divided[0], " ")))

	if len(divided) == 1 {
		return amount
	}

	unit := strings.Trim(divided[1], " ")

	switch unit {
	case "stacks", "stack":
		if len(divided) > 2 {
			return AMOUNT_STACK*amount + From_unit(strings.Join(divided[2:], " "))
		}
		return AMOUNT_STACK * amount
	case "shulkers", "shulker":
		if len(divided) > 2 {
			return AMOUNT_SHULKER*amount + From_unit(strings.Join(divided[2:], " "))
		}
		return AMOUNT_SHULKER * amount
	}

	return amount
}

func To_unit(amount int) string {
	if amount == 0 {
		return ""
	}

	if amount < AMOUNT_STACK {
		return fmt.Sprintf("%d items", amount)
	}

	if amount < AMOUNT_SHULKER {
		number := amount / AMOUNT_STACK
		s := "s"

		if number == 1 {
			s = ""
		}

		if amount%AMOUNT_STACK == 0 {
			return fmt.Sprintf("%d stack%s", number, s)
		}

		return fmt.Sprintf("%d stack%s, %s", number, s, To_unit(amount%AMOUNT_STACK))
	}

	number := amount / AMOUNT_SHULKER
	s := "s"

	if number == 1 {
		s = ""
	}

	if amount%AMOUNT_SHULKER == 0 {
		return fmt.Sprintf("%d shulker%s", number, s)
	}

	return fmt.Sprintf("%d shulker%s, %s", number, s, To_unit(amount%AMOUNT_SHULKER))
}

func Make_progress(amount, goal, quality int) (bar string) {
	floatProgress := float32(amount) / float32(goal) * float32(quality)
	partialProgress := int(floatProgress*100) % 100

	progress := int(floatProgress)

	for range progress {
		bar += ":white_large_square:"
	}

	if partialProgress <= 25 {
		bar += ":white_small_square:"
	} else if partialProgress <= 50 {
		bar += ":white_medium_small_square:"
	} else if partialProgress <= 75 {
		bar += ":white_medium_square:"
	} else {
		bar += ":black_large_square:"
	}

	for range quality - progress - 1 {
		bar += ":black_large_square:"
	}

	return
}
