package db

import "testing"

func TestEmbeddedMigrationsAreOrderedAndUnique(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(migrations) == 0 {
		t.Fatal("expected at least one migration")
	}
	for i, m := range migrations {
		if m.sql == "" {
			t.Fatalf("migration %s is empty", m.name)
		}
		if i > 0 && m.version <= migrations[i-1].version {
			t.Fatalf("migration %s is out of order", m.name)
		}
	}
}
