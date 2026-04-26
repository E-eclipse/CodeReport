package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type ignoreRule struct {
	pattern  string
	negate   bool
	dirOnly  bool
	anchored bool
	hasSlash bool
}

type IgnoreMatcher struct {
	rules []ignoreRule
}

func NewIgnoreMatcher(root string) *IgnoreMatcher {
	matcher := &IgnoreMatcher{}
	matcher.addDefaults()
	matcher.loadFile(filepath.Join(root, ".codereportignore"))
	return matcher
}

func (m *IgnoreMatcher) ShouldIgnore(relPath string, isDir bool) bool {
	relPath = filepath.ToSlash(filepath.Clean(relPath))
	if relPath == "." || relPath == "" {
		return false
	}

	ignored := false
	for _, rule := range m.rules {
		if rule.matches(relPath, isDir) {
			ignored = !rule.negate
		}
	}
	return ignored
}

func (m *IgnoreMatcher) addDefaults() {
	defaultPatterns := []string{
		".git/",
		".idea/",
		".vscode/",
		"__pycache__/",
		".pytest_cache/",
		".mypy_cache/",
		".ruff_cache/",
		".tox/",
		".cache/",
		".venv/",
		"venv/",
		"env/",
		"node_modules/",
		"vendor/",
		"dist/",
		"build/",
		"target/",
		"out/",
		"bin/",
		"obj/",
		".next/",
		".nuxt/",
		"coverage/",
		"htmlcov/",
		"*.pyc",
		"*.pyo",
		"*.pyd",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
		"*.class",
		"*.jar",
		"*.war",
		"*.ear",
		"*.zip",
		"*.tar",
		"*.gz",
		"*.rar",
		"*.7z",
		"*.doc",
		"*.docx",
		"*.xlsx",
		"*.pptx",
		"*.pdf",
		"*.png",
		"*.jpg",
		"*.jpeg",
		"*.gif",
		"*.svg",
		"*.webp",
		"*.ico",
		"*.mp3",
		"*.mp4",
		"*.webm",
		"*.ttf",
		"*.woff",
		"*.woff2",
		"*.eot",
		"*.otf",
		"*.map",
		"*.min.js",
		"*.min.css",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"poetry.lock",
		"Pipfile.lock",
		"*.sqlite",
		"*.sqlite3",
		"*.db",
		"*.log",
		".env",
		".DS_Store",
	}

	for _, pattern := range defaultPatterns {
		m.addRule(pattern)
	}
}

func (m *IgnoreMatcher) loadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.addRule(scanner.Text())
	}
}

func (m *IgnoreMatcher) addRule(line string) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}

	rule := ignoreRule{}
	if strings.HasPrefix(line, "!") {
		rule.negate = true
		line = strings.TrimSpace(strings.TrimPrefix(line, "!"))
	}
	if line == "" {
		return
	}

	line = filepath.ToSlash(line)
	if strings.HasPrefix(line, "/") {
		rule.anchored = true
		line = strings.TrimPrefix(line, "/")
	}
	if strings.HasSuffix(line, "/") {
		rule.dirOnly = true
		line = strings.TrimSuffix(line, "/")
	}

	rule.pattern = strings.Trim(line, "/")
	rule.hasSlash = strings.Contains(rule.pattern, "/")
	if rule.pattern != "" {
		m.rules = append(m.rules, rule)
	}
}

func (r ignoreRule) matches(relPath string, isDir bool) bool {
	if r.dirOnly && !isDir {
		return false
	}

	if r.anchored || r.hasSlash {
		return matchPattern(r.pattern, relPath)
	}

	parts := strings.Split(relPath, "/")
	for _, part := range parts {
		if matchPattern(r.pattern, part) {
			return true
		}
	}
	return false
}

func matchPattern(pattern, name string) bool {
	matched, err := filepath.Match(pattern, name)
	if err == nil && matched {
		return true
	}
	return pattern == name || strings.HasPrefix(name, pattern+"/")
}
