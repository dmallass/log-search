package search

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Allowed log directory path
const (
	AllowedLogDir = "/var/log/"
)

func ValidateFilePath(path string) (string, error) {
	if strings.Contains(path, "\x00") {
		return "", fmt.Errorf("invalid filepath, null bytes are found in the file path")
	}

	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid  filepath %v", err)
	}

	if !strings.HasPrefix(absPath, "/var/log/") && absPath != AllowedLogDir {
		return "", fmt.Errorf("invalid filepath, not an allowed filepath")
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("filepath doesn't exist %+v", err)
		}
		return "", fmt.Errorf("cannot access file %+v", err)
	}
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file %+v", err)
	}
	return absPath, nil
}

// LogLine represents a parsed log line with its timestamp
type LogLine struct {
	line      string
	timestamp time.Time
}

func RunRipgrep(logfile string, searchQuery string, searchMode string) ([]string, error) {
	threads := fmt.Sprintf("%d", runtime.NumCPU())

	if searchMode != "regex" {
		searchQuery = regexp.QuoteMeta(searchQuery)
	}

	args := []string{
		"-i",                 // case insensitive
		"--threads", threads, // use all available CPU cores
		"--color=never",    // Ensures output is plain text.
		"--mmap",           // use memory-mapped I/O
		"--no-ignore",      // don't respect .gitignore
		"--no-heading",     // don't show file names
		"--no-line-number", // don't show line numbers
		"--no-unicode",     // optimize for ASCII-only searches
		"--pre-glob=*.log", // only search log files
		"--no-filename",    // don't show filenames in output
		"--no-pcre2",       // disable PCRE2 for faster regex
		"--no-config",      // don't read configuration files
		"--no-require-git", // don't require git
	}

	if strings.EqualFold(searchMode, "fulltext") {
		args = append(args, "-F", searchQuery, logfile)
	} else {
		args = append(args, searchQuery, logfile)
	}

	cmd := exec.Command("rg", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	processState := cmd.ProcessState
	if processState.ExitCode() == 1 {
		return []string{}, nil // No matches found
	} else if processState.ExitCode() == 2 {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = "ripgrep encountered an error"
		}
		return nil, fmt.Errorf("ripgrep error: %s", errMsg)
	}
	if err != nil {
		return nil, fmt.Errorf("command error: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	var allResults []LogLine
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Extract timestamp (Assumed format: "YYYY-MM-DDTHH:MM:SS Log Message")
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339Nano, parts[0])
		if err != nil {
			continue // Skip if timestamp is invalid
		}

		allResults = append(allResults, LogLine{
			line:      line,
			timestamp: timestamp,
		})
	}

	// Sort results by timestamp in reverse order
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].timestamp.After(allResults[j].timestamp)
	})

	// Convert back to string slice
	var sortedLines []string
	for _, result := range allResults {
		sortedLines = append(sortedLines, result.line)
	}
	return sortedLines, nil
}

func paginateResults(results []string, page, limit int) []string {
	if len(results) <= 0 {
		return results
	}
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}
	end := page * limit
	if end > len(results) {
		end = len(results)
	}
	return results[start:end]
}

func searchLogfileHandler(c *gin.Context) {
	t1 := time.Now().UTC()
	logfile := c.Query("logfile")
	searchTerm := c.Query("q")
	searchMode := c.DefaultQuery("searchMode", "fulltext")

	// Validate inputs
	if logfile == "" || len(searchTerm) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "logfile and search query are required"})
		return
	}

	// Validate search mode
	if searchMode != "fulltext" && searchMode != "regex" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search mode: must be 'fulltext' or 'regex'"})
		return
	}

	validatedPath, err := ValidateFilePath(logfile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter (must be between 1 and 1000)"})
		return
	}

	results, err := RunRipgrep(validatedPath, searchTerm, searchMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedResults := paginateResults(results, page, limit)
	t2 := time.Now().UTC()
	c.JSON(http.StatusOK, gin.H{
		"total":      len(results),
		"page":       page,
		"limit":      limit,
		"query":      searchTerm,
		"logfile":    logfile,
		"searchMode": searchMode,
		"results":    paginatedResults,
		"duration":   fmt.Sprintf("%dms", t2.Sub(t1).Milliseconds()),
	})
}

func RunWebserver() error {
	r := gin.Default()
	r.GET("/search", searchLogfileHandler)
	if err := r.Run(":8080"); err != nil {
		return err
	}
	return nil
}
