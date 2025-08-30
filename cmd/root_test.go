package cmd

import (
	"ones-cfg2md/pkg/model"
	"os"
	"path/filepath"
	"reflect"
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

func TestRootCmd(t *testing.T) {
	// Define the source directory from fixtures
	sourceDir := filepath.Join("..", "fixtures", "input", "cfg")

	// Create a temporary directory for the output
	outputDir, err := os.MkdirTemp("", "output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputDir)

	// Execute the command
	rootCmd.SetArgs([]string{sourceDir, outputDir, "--types", "documents"})
	// a little hack to avoid `os.Exit(1)` in runConversion func
	var actualErr error
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		actualErr = runConversion(cmd, args)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if actualErr != nil {
		t.Fatal(actualErr)
	}

	// Check if the output directory contains the expected files
	expectedFile := filepath.Join(outputDir, "objects.csv")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file '%s' to be created, but it was not", expectedFile)
	}

	expectedMdFile := filepath.Join(outputDir, "Документ_Заказ.md")
	if _, err := os.Stat(expectedMdFile); os.IsNotExist(err) {
		t.Errorf("Expected file '%s' to be created, but it was not", expectedMdFile)
	}
}
