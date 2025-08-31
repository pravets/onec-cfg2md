package cmd

import (
	"bytes"
	"io"
	"onec-cfg2md/pkg/model"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestParseObjectTypes(t *testing.T) {
	testCases := []struct {
		name          string
		typesStr      string
		expectedTypes []model.ObjectType
		expectError   bool
	}{
		{
			name:     "Valid types",
			typesStr: "documents,catalogs",
			expectedTypes: []model.ObjectType{
				model.ObjectTypeDocument,
				model.ObjectTypeCatalog,
			},
			expectError: false,
		},
		{
			name:          "Single type",
			typesStr:      "enums",
			expectedTypes: []model.ObjectType{model.ObjectTypeEnum},
			expectError:   false,
		},
		{
			name:          "Empty string",
			typesStr:      "",
			expectedTypes: []model.ObjectType{model.ObjectTypeDocument},
			expectError:   false,
		},
		{
			name:          "Invalid type",
			typesStr:      "invalidtype",
			expectedTypes: nil,
			expectError:   true,
		},
		{
			name:          "Mixed valid and invalid types",
			typesStr:      "documents,invalidtype",
			expectedTypes: nil,
			expectError:   true,
		},
		{
			name:     "All types",
			typesStr: "documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes",
			expectedTypes: []model.ObjectType{
				model.ObjectTypeDocument,
				model.ObjectTypeCatalog,
				model.ObjectTypeAccumulationRegister,
				model.ObjectTypeInformationRegister,
				model.ObjectTypeEnum,
				model.ObjectTypeChartOfCharacteristicTypes,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTypes, err := parseObjectTypes(tc.typesStr)

			if (err != nil) != tc.expectError {
				t.Fatalf("Expected error: %v, but got: %v", tc.expectError, err)
			}

			if !reflect.DeepEqual(actualTypes, tc.expectedTypes) {
				t.Errorf("Expected types: %v, but got: %v", tc.expectedTypes, actualTypes)
			}
		})
	}
}

// Helper to capture stdout
func captureOutput(f func()) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestRootCmd_Scenarios(t *testing.T) {
	cfgSourceDir := filepath.Join("..", "fixtures", "input", "cfg")
	edtSourceDir := filepath.Join("..", "fixtures", "input", "edt")

	testCases := []struct {
		name          string
		args          []string
		setup         func(t *testing.T) (sourceDir string, cleanup func())
		expectError   bool
		errorContains string
		checkOutput   func(t *testing.T, outputDir string, stdout string)
	}{
		{
			name: "Successful CFG conversion (documents only)",
			args: []string{"--types", "documents"},
			setup: func(t *testing.T) (string, func()) {
				return cfgSourceDir, func() {}
			},
			expectError: false,
			checkOutput: func(t *testing.T, outputDir, stdout string) {
				if _, err := os.Stat(filepath.Join(outputDir, "objects.csv")); err != nil {
					t.Errorf("objects.csv not created: %v", err)
				}
				if _, err := os.Stat(filepath.Join(outputDir, "Документ_Заказ.md")); err != nil {
					t.Errorf("Документ_Заказ.md not created: %v", err)
				}
				if _, err := os.Stat(filepath.Join(outputDir, "Перечисление_СостоянияЗаказов.md")); err == nil {
					t.Errorf("Перечисление_СостоянияЗаказов.md should not be created")
				}
			},
		},
		{
			name: "Successful EDT conversion (enums only)",
			args: []string{"--types", "enums"},
			setup: func(t *testing.T) (string, func()) {
				return edtSourceDir, func() {}
			},
			expectError: false,
			checkOutput: func(t *testing.T, outputDir, stdout string) {
				if _, err := os.Stat(filepath.Join(outputDir, "objects.csv")); err != nil {
					t.Errorf("objects.csv not created: %v", err)
				}
				if _, err := os.Stat(filepath.Join(outputDir, "Перечисление_СостоянияЗаказов.md")); err != nil {
					t.Errorf("Перечисление_СостоянияЗаказов.md not created: %v", err)
				}
			},
		},
		{
			name: "Non-existent source directory",
			args: []string{},
			setup: func(t *testing.T) (string, func()) {
				return "non_existent_dir", func() {}
			},
			expectError:   true,
			errorContains: "ошибка определения формата",
		},
		{
			name: "Invalid format flag",
			args: []string{"--format", "xml"},
			setup: func(t *testing.T) (string, func()) {
				return cfgSourceDir, func() {}
			},
			expectError:   true,
			errorContains: "неподдерживаемый формат 'xml'",
		},
		{
			name: "Mismatched format flag (cfg for edt)",
			args: []string{"--format", "cfg"},
			setup: func(t *testing.T) (string, func()) {
				return edtSourceDir, func() {}
			},
			expectError:   true,
			errorContains: "не содержит файл Configuration.xml", // error from detector.ValidateFormat
		},
		{
			name: "No objects found of specified type",
			args: []string{"--types", "documents"},
			setup: func(t *testing.T) (string, func()) {
				// Create an empty dir with a valid config file but no objects
				tmpDir, err := os.MkdirTemp("", "empty-cfg")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				configContent := `<?xml version="1.0" encoding="UTF-8"?><MetaDataObject></MetaDataObject>`
				if err := os.WriteFile(filepath.Join(tmpDir, "Configuration.xml"), []byte(configContent), 0644); err != nil {
					t.Fatalf("Failed to write dummy config: %v", err)
				}
				return tmpDir, func() { _ = os.RemoveAll(tmpDir) }
			},
			expectError: false,
			checkOutput: func(t *testing.T, outputDir, stdout string) {
				if !strings.Contains(stdout, "Объекты указанных типов не найдены") {
					t.Errorf("Expected 'No objects found' message, got: %s", stdout)
				}
				// Check that dir is empty
				files, err := os.ReadDir(outputDir)
				if err != nil {
					t.Fatalf("Could not read output dir: %v", err)
				}
				if len(files) > 0 {
					t.Errorf("Output directory should be empty, but contains %d files", len(files))
				}
			},
		},
		{
			name: "Verbose flag output",
			args: []string{"--types", "documents", "-v"},
			setup: func(t *testing.T) (string, func()) {
				return cfgSourceDir, func() {}
			},
			expectError: false,
			checkOutput: func(t *testing.T, outputDir, stdout string) {
				expectedSubstrings := []string{
					"Определен формат: cfg",
					"Типы объектов для обработки:",
					"Начинаем парсинг метаданных...",
					"Генерируем Markdown файлы...",
					"Генерируем CSV каталог...",
				}
				for _, sub := range expectedSubstrings {
					if !strings.Contains(stdout, sub) {
						t.Errorf("Verbose output missing expected string '%s'", sub)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new command for each test case to avoid flag redefinition panic
			cmd := &cobra.Command{
				Use:  "onec-cfg2md <source_directory> <output_directory>",
				Args: cobra.ExactArgs(2),
				RunE: runConversion,
			}
			cmd.Flags().StringVar(&formatFlag, "format", "", "")
			cmd.Flags().StringVar(&typesFlag, "types", "documents,catalogs,accumulationregisters,informationregisters,enums,chartsofcharacteristictypes", "")
			cmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "")

			outputDir, err := os.MkdirTemp("", "output-*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.RemoveAll(outputDir) }()

			sourceDir, cleanup := tc.setup(t)
			defer cleanup()

			args := append([]string{sourceDir, outputDir}, tc.args...)
			cmd.SetArgs(args)

			var execErr error
			stdout, err := captureOutput(func() {
				execErr = cmd.Execute()
			})
			if err != nil {
				t.Fatalf("Failed to capture output: %v", err)
			}

			if tc.expectError {
				if execErr == nil {
					t.Fatal("Expected an error, but got none")
				}
				if tc.errorContains != "" && !strings.Contains(execErr.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s', but got: %v", tc.errorContains, execErr)
				}
			} else {
				if execErr != nil {
					t.Fatalf("Expected no error, but got: %v", execErr)
				}
			}

			if tc.checkOutput != nil {
				tc.checkOutput(t, outputDir, stdout)
			}
		})
	}
}
