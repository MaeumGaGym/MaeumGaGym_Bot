package utils

import (
	"bufio"
	"errors"
	"strings"
)

type IssueSections struct {
	Details  string
	Label    string
	Priority string
}

func ParseIssueContent(content string) (IssueSections, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var section IssueSections
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "### ") {
			switch trimmedLine {
			case "### 이슈 상세 내용":
				currentSection = "details"
			case "### 라벨":
				currentSection = "labels"
			case "### 우선순위":
				currentSection = "priority"
			default:
				currentSection = ""
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "-") {
			trimmedLine = strings.TrimPrefix(trimmedLine, "- ")
		}

		switch currentSection {
		case "details":
			section.Details += line + "\n"
		case "labels":
			if section.Label == "" {
				section.Label = trimmedLine
			}
		case "priority":
			if section.Priority == "" {
				section.Priority = trimmedLine
			}
		}
	}

	if section.Details == "" || section.Label == "" || section.Priority == "" {
		return section, errors.New("모든 섹션이 채워져 있어야 합니다")
	}

	return section, nil
}
