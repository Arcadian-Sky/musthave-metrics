//go:build !coverage
// +build !coverage

package caretaker

import (
	"fmt"
	"os"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

// Caretaker - объект, который хранит Memento
type Caretaker struct {
	Memento *storage.Memento
}

// saveToFile сохраняет метрики в файл
func (app *Caretaker) SaveToFile(m *storage.Memento, filename string) error {
	fmt.Printf("m: %v\n", m)
	jsonData, err := m.MarshalJSON()
	if err != nil {
		return err
	}
	// fmt.Printf("jsonData: %v\n", jsonData)
	// Записываем JSON в файл.
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// loadMetrics загружает ранее сохраненные метрики при старте сервера
func (app *Caretaker) ReadFromFile(filename string) (*storage.Memento, error) {
	if filename == "" {
		return nil, nil
	}
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil
	}

	m := &storage.Memento{}
	err = m.UnmarshalJSON(jsonData)
	if err != nil {
		return nil, nil
	}

	return m, nil
}
