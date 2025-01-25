package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Config struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	Name          string   `json:"name"`
	Parent        string   `json:"parent,omitempty"`
	Rows          int      `json:"rows,omitempty"`
	RowsPerParent int      `json:"rows_per_parent,omitempty"`
	Columns       []Column `json:"columns"`
}

type Column struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type,omitempty"`
	Generator map[string]interface{} `json:"generator"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	configFolder := "./config"
	fileList, err := listConfigFiles(configFolder)
	if err != nil {
		panic(fmt.Sprintf("Failed to list configuration files: %v", err))
	}

	fmt.Println("Available configuration files:")
	for i, file := range fileList {
		fmt.Printf("[%d] %s\n", i+1, file)
	}

	var selectedFile string
	for {
		fmt.Print("Enter the number of the configuration file to use: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err == nil && choice > 0 && choice <= len(fileList) {
			selectedFile = fileList[choice-1]
			break
		} else {
			fmt.Println("Invalid choice. Please enter a valid number.")
		}
	}

	configFilePath := filepath.Join(configFolder, selectedFile)
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read configuration file: %v", err))
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}

	fmt.Print("Enter the output directory for the CSV files (default: ./output): ")
	outputDir, _ := reader.ReadString('\n')
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		outputDir = "./output"
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create output directory: %v", err))
	}

	tableCounters := make(map[string]int)
	for _, table := range config.Tables {
		for _, col := range table.Columns {
			if _, ok := col.Generator["table_increment"].(bool); ok {
				tableCounters[table.Name] = 1
				break
			}
		}
	}

	writers, files := createWriters(config.Tables, outputDir)
	defer cleanup(writers, files)

	// Generate data
	for _, table := range config.Tables {
		if table.Parent != "" {
			continue
		}
		generateTableWithProgress(table, config.Tables, writers, tableCounters)
	}

	fmt.Printf("Finish generating data. Output can be found on %s folder", outputDir)
}

func listConfigFiles(folder string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			files = append(files, info.Name())
		}
		return nil
	})
	return files, err
}

func createWriters(tables []Table, outputDir string) (map[string]*csv.Writer, map[string]*os.File) {
	writers := make(map[string]*csv.Writer)
	files := make(map[string]*os.File)

	for _, table := range tables {
		filePath := filepath.Join(outputDir, table.Name + ".csv")
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		files[table.Name] = file
		writer := csv.NewWriter(file)

		// Write header
		headers := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			headers[i] = col.Name
		}
		writer.Write(headers)
		writers[table.Name] = writer
	}

	return writers, files
}

func cleanup(writers map[string]*csv.Writer, files map[string]*os.File) {
	for _, w := range writers {
		w.Flush()
	}
	for _, f := range files {
		f.Close()
	}
}

func generateTableWithProgress(table Table, allTables []Table, writers map[string]*csv.Writer, counters map[string]int) {
	rowCount := table.Rows
	fmt.Printf("Generating data for table (including child): %s (%d rows)\n", table.Name, rowCount)

	bar := progressbar.NewOptions(rowCount,
		progressbar.OptionSetDescription(fmt.Sprintf("Table: %s", table.Name)),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionClearOnFinish(),
	)

	writer := writers[table.Name]

	for i := 0; i < table.Rows; i++ {
		row := make(map[string]string)
		record := make([]string, len(table.Columns))

		for colIdx, column := range table.Columns {
			value := generateColumnValue(column, i, nil, table.Name, counters)
			record[colIdx] = value
			row[column.Name] = value
		}

		writer.Write(record)

		// Generate child tables
		for _, childTable := range getChildTables(table.Name, allTables) {
			generateChildTable(childTable, writers, row, counters, allTables)
		}
		bar.Add(1)
		// time.Sleep(50 * time.Millisecond) // uncomment this to simulate progress bar
	}

	fmt.Printf("Completed data generation for table: %s\n", table.Name)
}

func generateChildTable(table Table, writers map[string]*csv.Writer, parentRow map[string]string, counters map[string]int, allTables []Table) {
	writer := writers[table.Name]

	for i := 0; i < table.RowsPerParent; i++ {
		childRow := make(map[string]string)
		record := make([]string, len(table.Columns))

		for colIdx, column := range table.Columns {
			value := generateColumnValue(column, i, parentRow, table.Name, counters)
			record[colIdx] = value
			childRow[column.Name] = value
		}

		writer.Write(record)

		// Process grandchildren
		for _, grandchildTable := range getChildTables(table.Name, allTables) {
			generateChildTable(grandchildTable, writers, childRow, counters, allTables)
		}
	}
}

func generateColumnValue(col Column, index int, parentRow map[string]string, tableName string, counters map[string]int) string {
	gen := col.Generator

	if hc, ok := gen["hardcoded"].(string); ok {
		return hc
	}

	if _, ok := gen["increment"]; ok {
		return strconv.Itoa(index + 1)
	}

	if _, ok := gen["table_increment"].(bool); ok {
		current := counters[tableName]
		counters[tableName]++
		return strconv.Itoa(current)
	}

	if pk, ok := gen["parent_key"].(string); ok {
		return parentRow[pk]
	}

	if randomCfg, ok := gen["random"].(map[string]interface{}); ok {
		length := 10
		if l, ok := randomCfg["length"].(float64); ok {
			length = int(l)
		}

		prefix, _ := randomCfg["prefix"].(string)
		suffix, _ := randomCfg["suffix"].(string)

		return prefix + randomString(length) + suffix
	}

	if predefinedList, ok := gen["predefined_list"].([]interface{}); ok {
		if len(predefinedList) > 0 {
			currentIndex := index % len(predefinedList)
			if value, ok := predefinedList[currentIndex].(string); ok {
				return value
			}
		}
	}

	return ""
}

func randomString(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}
	return sb.String()
}

func getChildTables(parent string, allTables []Table) []Table {
	var children []Table
	for _, t := range allTables {
		if t.Parent == parent {
			children = append(children, t)
		}
	}
	return children
}