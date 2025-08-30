package cmd

import (
    "os"
    "path/filepath"
    "testing"
)

func TestExecute_WithCFGFixtures(t *testing.T) {
    fixtures := filepath.Join("..", "fixtures", "input", "cfg")
    out := t.TempDir()

    // Сохраняем и восстанавливаем флаги
    of, ot, ov := formatFlag, typesFlag, verboseFlag
    defer func() { formatFlag, typesFlag, verboseFlag = of, ot, ov }()

    rootCmd.SetArgs([]string{fixtures, out})
    if err := Execute(); err != nil {
        t.Fatalf("Execute() failed for CFG fixtures: %v", err)
    }

    if _, err := os.Stat(filepath.Join(out, "objects.csv")); err != nil {
        t.Fatalf("objects.csv not created: %v", err)
    }
}

func TestExecute_WithEDTFixtures(t *testing.T) {
    fixtures := filepath.Join("..", "fixtures", "input", "edt")
    out := t.TempDir()

    of, ot, ov := formatFlag, typesFlag, verboseFlag
    defer func() { formatFlag, typesFlag, verboseFlag = of, ot, ov }()

    rootCmd.SetArgs([]string{fixtures, out, "--types", "enums"})
    if err := Execute(); err != nil {
        t.Fatalf("Execute() failed for EDT fixtures: %v", err)
    }

    if _, err := os.Stat(filepath.Join(out, "objects.csv")); err != nil {
        t.Fatalf("objects.csv not created: %v", err)
    }
}
