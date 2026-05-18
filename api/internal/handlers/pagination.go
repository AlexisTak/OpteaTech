package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type pagination struct {
	offset int
	limit  int
	valid  bool
}

func parsePagination(c fiber.Ctx) pagination {
	startRaw := c.Query("_start")
	endRaw := c.Query("_end")
	if startRaw != "" && endRaw != "" {
		start, errStart := strconv.Atoi(startRaw)
		end, errEnd := strconv.Atoi(endRaw)
		if errStart == nil && errEnd == nil && start >= 0 && end > start {
			return pagination{offset: start, limit: end - start, valid: true}
		}
	}

	pageRaw := c.Query("page")
	limitRaw := c.Query("limit")
	if pageRaw != "" && limitRaw != "" {
		page, errPage := strconv.Atoi(pageRaw)
		limit, errLimit := strconv.Atoi(limitRaw)
		if errPage == nil && errLimit == nil && page > 0 && limit > 0 {
			return pagination{offset: (page - 1) * limit, limit: limit, valid: true}
		}
	}

	return pagination{offset: 0, limit: 0, valid: false}
}

func applySlicePagination[T any](items []T, p pagination) []T {
	if !p.valid {
		return items
	}
	if p.offset >= len(items) {
		return []T{}
	}
	end := p.offset + p.limit
	if end > len(items) {
		end = len(items)
	}
	return items[p.offset:end]
}

func setTotalCountHeader(c fiber.Ctx, total int) {
	c.Set("X-Total-Count", strconv.Itoa(total))
}

func buildOrderClause(c fiber.Ctx, allowed map[string]string, defaultColumn string, defaultDirection string) string {
	sortKey := strings.TrimSpace(c.Query("_sort"))
	order := strings.ToUpper(strings.TrimSpace(c.Query("_order")))

	column := defaultColumn
	if mapped, ok := allowed[sortKey]; ok {
		column = mapped
	}

	direction := strings.ToUpper(defaultDirection)
	if order == "ASC" || order == "DESC" {
		direction = order
	}

	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}

	return column + " " + direction
}

func buildWhereClause(c fiber.Ctx, allowed map[string]string, startIndex int) (string, []interface{}, int) {
	queries := c.Queries()
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	index := startIndex

	for key, value := range queries {
		if value == "" {
			continue
		}
		if isReservedQueryKey(key) {
			continue
		}

		suffix := ""
		baseKey := key
		for _, candidate := range []string{"_contains", "_in", "_gte", "_lte"} {
			if strings.HasSuffix(key, candidate) {
				suffix = candidate
				baseKey = strings.TrimSuffix(key, candidate)
				break
			}
		}

		column, ok := allowed[baseKey]
		if !ok {
			continue
		}

		switch suffix {
		case "_contains":
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", column, index))
			args = append(args, "%"+value+"%")
			index++
		case "_in":
			parts := strings.Split(value, ",")
			clean := make([]string, 0, len(parts))
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" {
					clean = append(clean, trimmed)
				}
			}
			if len(clean) == 0 {
				continue
			}
			conditions = append(conditions, fmt.Sprintf("%s = ANY($%d)", column, index))
			args = append(args, clean)
			index++
		case "_gte":
			conditions = append(conditions, fmt.Sprintf("%s >= $%d", column, index))
			args = append(args, value)
			index++
		case "_lte":
			conditions = append(conditions, fmt.Sprintf("%s <= $%d", column, index))
			args = append(args, value)
			index++
		default:
			conditions = append(conditions, fmt.Sprintf("%s = $%d", column, index))
			args = append(args, value)
			index++
		}
	}

	if len(conditions) == 0 {
		return "", args, index
	}

	return "WHERE " + strings.Join(conditions, " AND "), args, index
}

func isReservedQueryKey(key string) bool {
	switch key {
	case "_start", "_end", "page", "limit", "_sort", "_order":
		return true
	default:
		return false
	}
}
