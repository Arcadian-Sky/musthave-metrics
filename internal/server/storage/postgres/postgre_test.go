package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

func NewTestPostgresStorage(db *sql.DB) *PostgresStorage {
	p := &PostgresStorage{db: db, retriableErrorMap: map[string]bool{pgerrcode.UniqueViolation: true}, maxRetries: 3, initialDelay: time.Second}
	return p
}

func TestPostgreStorage_GetMetrics(t *testing.T) {

	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	p := NewTestPostgresStorage(db)

	// Создание строки результата для mock запроса
	rows := sqlmock.NewRows([]string{"name", "type", "counter", "gauge"}).
		AddRow("metric1", "counter", 10, nil).
		AddRow("metric2", "gauge", nil, 20.0)

	// Ожидаемый результат выполнения mock запроса
	mock.ExpectQuery("SELECT name, type, counter, gauge FROM .*").
		WillReturnRows(rows)

	// Выполнение тестируемого метода
	metrics := p.GetMetrics(context.Background())
	fmt.Printf("metrics: %v\n", metrics)
	// Ожидаемый результат
	expectedMetrics := map[storage.MetricType]map[string]interface{}{
		storage.Counter: {
			"metric1": int64(10),
		},
		storage.Gauge: {
			"metric2": 20.0,
		},
	}

	// Проверка на соответствие ожидаемому результату
	assert.Equal(t, expectedMetrics, metrics)
}

func TestPostgreStorage_GetJSONMetrics(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	// Создание пустого среза метрик
	var metrics []models.Metrics

	err = p.GetJSONMetrics(context.Background(), &metrics)
	if err != nil {
		t.Fatalf("ошибка при получении метрик: %v", err)
	}

	// Ожидается, что метод вернет nil без каких-либо действий
}

func TestPostgreStorage_GetMetric(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	// Ожидаемый результат выполнения mock запроса
	mock.ExpectQuery("SELECT name, gauge FROM .*").
		WillReturnRows(sqlmock.NewRows([]string{"name", "gauge"}).
			AddRow("metric1", 10.5).
			AddRow("metric2", 20.0))

	// Выполнение тестируемого метода для типа метрики Gauge
	metrics := p.GetMetric(context.Background(), storage.Gauge)
	if metrics == nil {
		t.Fatalf("метрики не были получены")
	}

	// Ожидаемый результат
	expectedMetrics := map[string]interface{}{
		"metric1": 10.5,
		"metric2": 20.0,
	}

	// Проверка на соответствие ожидаемому результату
	assert.Equal(t, expectedMetrics, metrics)
}

func TestPostgreStorage_GetJSONMetric(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage
	p := NewTestPostgresStorage(db)

	// Создание тестовой метрики
	testMetric := &models.Metrics{
		ID:    "test",
		MType: "gauge",
	}

	// Ожидаемый результат сканирования строк
	mock.ExpectQuery("SELECT name, type, counter, gauge FROM .*").WithArgs(testMetric.ID, testMetric.MType).
		WillReturnRows(sqlmock.NewRows([]string{"name", "type", "counter", "gauge"}).AddRow("test", "gauge", nil, 42.0))

	// Выполнение тестируемого метода
	err = p.GetJSONMetric(context.Background(), testMetric)

	// Проверка на отсутствие ошибки
	assert.NoError(t, err)

	// Проверка на соответствие ожидаемого значения
	assert.Equal(t, 42.0, *testMetric.Value)
}

func TestPostgreStorage_UpdateMetric(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage
	p := NewTestPostgresStorage(db)

	// Ожидаемый запрос SQL
	mock.ExpectExec("INSERT INTO .*").WithArgs("test", 42).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Выполнение тестируемого метода
	err = p.UpdateMetric(context.Background(), "gauge", "test", "42")

	// Проверка на отсутствие ошибки
	assert.NoError(t, err)
}

func TestPostgreStorage_UpdateJSONMetric(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	var v = float64(10.5)
	// Создание метрики для обновления
	metric := &models.Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: &v,
	}

	// Ожидаемый результат выполнения mock запроса
	mock.ExpectExec("INSERT INTO .*").
		WithArgs("metric1", metric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Выполнение тестируемого метода
	err = p.UpdateJSONMetric(context.Background(), metric)
	if err != nil {
		t.Fatalf("ошибка при обновлении метрики: %v", err)
	}
}

func TestPostgreStorage_UpdateJSONMetrics(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	// Создание тестовых метрик
	metrics := []models.Metrics{
		{ID: "metric1", MType: "counter", Delta: new(int64)},
		{ID: "metric2", MType: "gauge", Value: new(float64)},
	}

	// Ожидаемые SQL запросы
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Выполнение тестируемого метода
	err = p.UpdateJSONMetrics(context.Background(), &metrics)
	if err != nil {
		t.Fatalf("ошибка при обновлении метрик: %v", err)
	}

	// Проверка выполнения всех ожидаемых SQL запросов
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ожидаемые SQL запросы не были выполнены: %s", err)
	}
}

func TestPostgreStorage_SetMetrics(t *testing.T) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	// Подготовка тестовых данных
	metrics := map[storage.MetricType]map[string]interface{}{
		storage.Gauge: {
			"metric1": 10.0,
			"metric2": 20.0,
		},
		storage.Counter: {
			"metric3": 100,
			"metric4": 200,
		},
	}

	// Ожидаемые SQL запросы
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Выполнение тестируемого метода
	p.SetMetrics(context.Background(), metrics)

	// Проверка выполнения всех ожидаемых SQL запросов
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ожидаемые SQL запросы не были выполнены: %s", err)
	}
}

func BenchmarkPostgresStorage_GetMetric(b *testing.B) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, _, err := sqlmock.New()
	if err != nil {
		b.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.GetMetric(ctx, storage.Gauge)
	}
}

func BenchmarkPostgresStorage_SetMetrics(b *testing.B) {
	// Создание mock базы данных и отложенное закрытие соединения
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()

	// Создание экземпляра PostgresStorage с помощью конструктора NewPostgresStorage
	p := NewTestPostgresStorage(db)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		metrics := map[storage.MetricType]map[string]interface{}{
			storage.Gauge: {
				"metric1": 10.0,
				"metric2": 20.0,
			},
			storage.Counter: {
				"metric3": 100,
				"metric4": 200,
			},
		}

		// Выполнение тестируемого метода
		p.SetMetrics(ctx, metrics)
	}
}
