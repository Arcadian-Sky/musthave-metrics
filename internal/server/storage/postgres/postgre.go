package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
)

// PostgresStorage представляет хранилище метрик в PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новый экземпляр PostgresStorage
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	p := &PostgresStorage{db: db}
	err := p.CreateMetricsTable()
	if err != nil {
		panic(err)
	}
	return p

}

func (p *PostgresStorage) getTableName() string {
	return "metrics"
}

func (p *PostgresStorage) Ping() error {
	err := p.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) CreateMetricsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS metrics (
			"name" varchar NOT NULL,
			counter int NULL,
			gauge double precision NULL,
			"type" varchar NOT NULL,
			CONSTRAINT constraint_name_type UNIQUE (name, type)
		);
    `
	_, err := p.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы метрик: %v", err)
	}
	return nil
}

func (p *PostgresStorage) GetJSONMetric(metric *models.Metrics) error {
	query := fmt.Sprintf("SELECT name, type, counter, gauge FROM %s WHERE name = $1 AND type = $2", p.getTableName())
	err := p.db.QueryRow(query, metric.ID, metric.MType).Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return fmt.Errorf("ошибка при извлечении метрики из базы данных: %v", err)
	}
	return nil
}

func (p *PostgresStorage) UpdateJSONMetric(metric *models.Metrics) error {
	query := fmt.Sprintf("INSERT INTO %s (name, type, counter, gauge) VALUES ($1, $2, $3, $4) ON CONFLICT (name, type) DO UPDATE SET counter = $3, gauge = $4", p.getTableName())
	_, err := p.db.Exec(query, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении метрики в базе данных: %v", err)
	}
	return nil
}

func (p *PostgresStorage) UpdateMetric(mtype string, name string, value string) error {
	// Получаем тип метрики
	metricType, err := storage.GetMetricTypeByCode(mtype)
	if err != nil {
		return err
	}

	// Определяем запрос SQL в зависимости от типа метрики
	var query string
	switch metricType {
	case storage.Gauge:
		query = fmt.Sprintf(`
            INSERT INTO %s (name, type, gauge)
            VALUES ($1, 'gauge', $2)
            ON CONFLICT (name, type) DO UPDATE
            SET gauge = EXCLUDED.gauge
        `, p.getTableName())
	case storage.Counter:
		query = fmt.Sprintf(`
            INSERT INTO %s (name, type, counter)
            VALUES ($1, 'counter', $2)
            ON CONFLICT (name, type) DO UPDATE
            SET counter = EXCLUDED.counter
        `, p.getTableName())
	default:
		return fmt.Errorf("неподдерживаемый тип метрики: %s", mtype)
	}

	// Выполняем SQL запрос
	_, err = p.db.Exec(query, name, value)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении метрики в базе данных: %v", err)
	}
	return nil
}

func (p *PostgresStorage) GetMetric(mtype storage.MetricType) map[string]interface{} {
	var query string
	var metrics map[string]interface{}
	fmt.Printf("mtype: %v\n", mtype)
	switch mtype {
	case storage.Gauge:
		query = fmt.Sprintf("SELECT name, gauge FROM %s WHERE type = 'gauge'", p.getTableName())
	case storage.Counter:
		query = fmt.Sprintf("SELECT name, counter FROM %s WHERE type = 'counter'", p.getTableName())
	default:
		return nil //, fmt.Errorf("неподдерживаемый тип метрики: %s", metricType)
	}

	rows, err := p.db.Query(query)
	if err != nil {
		return nil //, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()
	// Проверка наличия ошибок при чтении строк результата
	if rows.Err() != nil {
		return nil
		//, fmt.Errorf("ошибка при чтении строк результата: %v", rows.Err())
	}

	metrics = make(map[string]interface{})
	for rows.Next() {
		var name string
		var value interface{}

		if mtype == storage.Gauge {
			var gauge float64
			err := rows.Scan(&name, &gauge)
			if err != nil {
				return nil //, fmt.Errorf("ошибка при сканировании строки: %v", err)
			}
			value = gauge
		} else if mtype == storage.Counter {
			var counter int
			err := rows.Scan(&name, &counter)
			if err != nil {
				return nil //, fmt.Errorf("ошибка при сканировании строки: %v", err)
			}
			value = counter
		}
		metrics[name] = value
	}
	return metrics //, nil
}

// TODO:Add support error handling
func (p *PostgresStorage) SetMetrics(metrics map[storage.MetricType]map[string]interface{}) {
	// Начало транзакции
	tx, err := p.db.Begin()
	if err != nil {
		panic(fmt.Errorf("ошибка при начале транзакции: %v", err))
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println("Ошибка отката транзакции:", err)
		}
	}()

	// Цикл по всем типам метрик
	for metricType, metricValues := range metrics {
		// Определение имени столбца для текущего типа метрики
		var columnName string
		switch metricType {
		case storage.Gauge:
			columnName = "gauge"
		case storage.Counter:
			columnName = "counter"
		default:
			return
			// return fmt.Errorf("неподдерживаемый тип метрики: %v", metricType)
		}

		// Цикл по всем метрикам текущего типа
		for name, value := range metricValues {
			// Формирование запроса SQL для обновления или вставки метрики
			query := fmt.Sprintf(`
                INSERT INTO %s (name, type, %s)
                VALUES ($1, $2, $3)
                ON CONFLICT (name, type) DO UPDATE
                SET %s = EXCLUDED.%s
            `, p.getTableName(), columnName, columnName, columnName)

			// Выполнение запроса SQL внутри транзакции
			_, err := tx.Exec(query, name, string(metricType), value)
			if err != nil {
				return
				// fmt.Errorf("ошибка при выполнении запроса: %v", err)
			}
		}
	}

	// Фиксация транзакции
	if err := tx.Commit(); err != nil {
		panic(fmt.Errorf("ошибка при фиксации транзакции: %v", err))
		// return fmt.Errorf("ошибка при фиксации транзакции: %v", err)
	}

}

// TODO:Add support error handling
func (p *PostgresStorage) GetMetrics() map[storage.MetricType]map[string]interface{} {
	metrics := make(map[storage.MetricType]map[string]interface{})
	query := fmt.Sprintf("SELECT name, type, counter, gauge FROM %s", p.getTableName())

	// Выполнение запроса SQL
	rows, err := p.db.Query(query)
	if err != nil {
		return nil
		//, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()
	// Проверка наличия ошибок при чтении строк результата
	if rows.Err() != nil {
		return nil
		//, fmt.Errorf("ошибка при чтении строк результата: %v", rows.Err())
	}

	for rows.Next() {
		var name string
		var mType string
		var counter sql.NullInt64
		var gauge sql.NullFloat64

		// Сканирование значений строки результата
		if err := rows.Scan(&name, &mType, &counter, &gauge); err != nil {
			return nil
			//, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}

		// Преобразование строкового представления типа метрики в тип MetricType
		metricType, err := storage.GetMetricTypeByCode(mType)
		if err != nil {
			return nil
			//, fmt.Errorf("ошибка при преобразовании типа метрики: %v", err)
		}

		// Создание вложенной карты метрик для данного типа метрики, если она еще не существует
		if _, ok := metrics[metricType]; !ok {
			metrics[metricType] = make(map[string]interface{})
		}

		// Запись метрики в соответствующую карту
		if counter.Valid {
			metrics[metricType][name] = counter.Int64
		} else if gauge.Valid {
			metrics[metricType][name] = gauge.Float64
		}
	}

	return metrics
	//, nil
}
