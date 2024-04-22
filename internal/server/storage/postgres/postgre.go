package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
	"github.com/Arcadian-Sky/musthave-metrics/migrations"
)

// PostgresStorage представляет хранилище метрик в PostgreSQL
type PostgresStorage struct {
	db                *sql.DB
	retriableErrorMap map[string]bool
	maxRetries        int
	initialDelay      time.Duration
}

// NewPostgresStorage создает новый экземпляр PostgresStorage
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	p := &PostgresStorage{
		db: db,
		retriableErrorMap: map[string]bool{
			pgerrcode.UniqueViolation: true,
		},
		maxRetries:   3,
		initialDelay: time.Second,
	}
	err := p.migrateDB()
	if err != nil {
		log.Fatal(err)
	}
	return p
}

// // executeWithRetry - функция для выполнения операции с повторными попытками
// func (p *PostgresStorage) executeWithRetry(ctx context.Context, operation func() error) error {
// 	var err error
// 	delay := p.initialDelay
// 	for attempt := 0; attempt <= p.maxRetries; attempt++ {
// 		if err = operation(); err == nil || !p.isRetriableError(err) {
// 			return err
// 		}

// 		// Если ошибка retriable, ждем перед следующей попыткой
// 		fmt.Printf("Retriable error occurred: %v. Retrying after %s...\n", err, delay)
// 		select {
// 		case <-time.After(delay):
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		}
// 		delay *= 2 // Увеличиваем задержку перед следующей попыткой
// 	}
// 	return err
// }

func (p *PostgresStorage) executeWithRetry(_ context.Context, operation func() error) error {
	// Создаем экземпляр стратегии повторных попыток
	backoffStrategy := backoff.NewExponentialBackOff()
	backoffStrategy.MaxElapsedTime = time.Duration(p.maxRetries) * p.initialDelay

	// Обертываем операцию в функцию, которую backoff будет повторять
	retryOperation := func() error {
		err := operation()
		if err != nil && p.isRetriableError(err) {
			fmt.Printf("Retriable error occurred: %v\n", err)
			return err
		}
		return nil
	}

	// Выполняем операцию с использованием backoff
	return backoff.Retry(retryOperation, backoffStrategy)
}

// isRetryableError - функция для проверки, является ли ошибка retriable
func (p *PostgresStorage) isRetriableError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return p.retriableErrorMap[pgErr.Code]
	}
	return false
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

func (p *PostgresStorage) migrateDB() error {
	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := goose.RunContext(ctx, "up", p.db, ".")
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetMetric(ctx context.Context, mtype storage.MetricType) map[string]interface{} {
	var query string
	var metrics map[string]interface{}
	fmt.Printf("mtype: %v\n", mtype)
	query = fmt.Sprintf("SELECT name, gauge FROM %s WHERE type = '$1'", p.getTableName())
	var queryType string
	switch mtype {
	case storage.Gauge:
		queryType = "gauge"
		// query = fmt.Sprintf("SELECT name, gauge FROM %s WHERE type = 'gauge'", p.getTableName())
	case storage.Counter:
		queryType = "counter"
		// query = fmt.Sprintf("SELECT name, counter FROM %s WHERE type = 'counter'", p.getTableName())
	default:
		return nil //, fmt.Errorf("неподдерживаемый тип метрики: %s", metricType)
	}

	rows, err := p.db.Query(query, queryType)
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

func (p *PostgresStorage) GetJSONMetric(ctx context.Context, metric *models.Metrics) error {
	query := fmt.Sprintf("SELECT name, type, counter, gauge FROM %s WHERE name = $1 AND type = $2", p.getTableName())
	row := p.db.QueryRowContext(ctx, query, metric.ID, metric.MType)

	var counter sql.NullInt64
	var gauge sql.NullFloat64
	err := row.Scan(&metric.ID, &metric.MType, &counter, &gauge)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// return nil
		} else {
			// Возникла другая ошибка
			return err
		}
	}

	metricType, err := utils.GetMetricTypeByCode(metric.MType)
	if err != nil {
		return err
	}
	switch metricType {
	case storage.Gauge:
		if gauge.Valid {
			metric.Value = &gauge.Float64
		} else {
			metric.Value = new(float64)
		}
	case storage.Counter:
		if counter.Valid {
			metric.Delta = &counter.Int64
		} else {
			metric.Delta = new(int64)
		}
	}

	return nil
}

// TODO:Add support error handling
func (p *PostgresStorage) GetMetrics(ctx context.Context) map[storage.MetricType]map[string]interface{} {
	metrics := make(map[storage.MetricType]map[string]interface{})
	query := fmt.Sprintf("SELECT name, type, counter, gauge FROM %s", p.getTableName())

	// Выполнение запроса SQL
	rows, err := p.db.Query(query)
	if err != nil {
		fmt.Printf("1err: %v\n", err)
		return nil
		//, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()
	// Проверка наличия ошибок при чтении строк результата
	if rows.Err() != nil {
		fmt.Printf("2err: %v\n", err)
		return nil
		//, fmt.Errorf("ошибка при чтении строк результата: %v", rows.Err())
	}

	for rows.Next() {
		var row struct {
			name    string `field:"name"`
			mType   string `field:"type"`
			counter sql.NullInt64
			gauge   sql.NullFloat64
		}

		// Сканирование значений строки результата
		if err := rows.Scan(&row.name, &row.mType, &row.counter, &row.gauge); err != nil {
			fmt.Printf("3err: %v\n", err)
			return nil
			//, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}
		// Преобразование строкового представления типа метрики в тип MetricType
		metricType, err := utils.GetMetricTypeByCode(row.mType)
		if err != nil {
			fmt.Printf("3err: %v\n", err)
			return nil
			//, fmt.Errorf("ошибка при преобразовании типа метрики: %v", err)
		}

		// Создание вложенной карты метрик для данного типа метрики, если она еще не существует
		if _, ok := metrics[metricType]; !ok {
			metrics[metricType] = make(map[string]interface{})
		}

		// Запись метрики в соответствующую карту
		if row.counter.Valid {
			metrics[metricType][row.name] = row.counter.Int64
		} else if row.gauge.Valid {
			metrics[metricType][row.name] = row.gauge.Float64
		}
	}

	return metrics
}

func (p *PostgresStorage) GetJSONMetrics(ctx context.Context, metric *[]models.Metrics) error {
	return nil
}

func (p *PostgresStorage) UpdateMetric(ctx context.Context, mtype string, name string, value string) error {
	// Получаем тип метрики
	metricType, err := utils.GetMetricTypeByCode(mtype)
	if err != nil {
		return err
	}

	// Определяем запрос SQL в зависимости от типа метрики
	var query string
	var reValue interface{}
	switch metricType {
	case storage.Gauge:
		query = "INSERT INTO " + p.getTableName() + " (name, type, gauge)" +
			"VALUES ($1, 'gauge', $2)" +
			"ON CONFLICT (name, type) DO UPDATE" +
			"SET gauge = EXCLUDED.gauge"
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			reValue = floatValue
		}
	case storage.Counter:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			var currentCounter sql.NullInt64
			err := p.executeWithRetry(context.Background(), func() error {
				return p.db.QueryRowContext(ctx, "SELECT counter FROM "+p.getTableName()+" WHERE name = $1 AND type = 'counter'", name).Scan(&currentCounter)
			})
			if err != nil {
				return err
			}
			reValue = currentCounter.Int64 + intValue
			query = "INSERT INTO " + p.getTableName() + " (name, type, counter)" +
				" VALUES ($1, 'counter', $2)" +
				" ON CONFLICT (name, type) DO UPDATE" +
				" SET counter = EXCLUDED.counter"
		} else {
			return fmt.Errorf("invalid metric value: %v", err)
		}

	default:
		return fmt.Errorf("неподдерживаемый тип метрики: %s", metricType)
	}

	// Выполняем SQL запрос
	err = p.executeWithRetry(context.Background(), func() error {
		_, err = p.db.Exec(query, name, reValue)
		return err
	})

	if err != nil {
		return fmt.Errorf("ошибка при обновлении метрики в базе данных: %v", err)
	}
	return nil
}

func (p *PostgresStorage) UpdateJSONMetric(ctx context.Context, metric *models.Metrics) error {
	mType, err := utils.GetMetricTypeByCode(metric.MType)
	if err != nil {
		return err
	}
	var query string
	var value any
	switch mType {
	case storage.Gauge:
		if metric.Value == nil {
			return nil
		}
		query = fmt.Sprintf(`
				INSERT INTO %s (name, type, gauge)
				VALUES ($1, 'gauge', $2)
				ON CONFLICT (name, type) DO UPDATE
				SET gauge = EXCLUDED.gauge
			`, p.getTableName())
		value = metric.Value
	case storage.Counter:
		if metric.Delta == nil {
			return nil
		}
		var currentCounter sql.NullInt64
		_ = p.db.QueryRowContext(ctx, "SELECT counter FROM "+p.getTableName()+" WHERE name = $1 AND type = 'counter'", metric.ID).Scan(&currentCounter)

		*metric.Delta += currentCounter.Int64
		fmt.Printf("metric.Delta: %v\n", metric.Delta)
		query = "INSERT INTO " + p.getTableName() + " (name, type, counter)" +
			" VALUES ($1, 'counter', $2)" +
			" ON CONFLICT (name, type) DO UPDATE" +
			" SET counter = EXCLUDED.counter"
		value = metric.Delta
	default:
		return fmt.Errorf("неподдерживаемый тип метрики: %s", mType)
	}

	_, err = p.db.Exec(query, metric.ID, value)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении метрики в базе данных: %v", err)
	}

	return nil
}

// TODO:Add support error handling
func (p *PostgresStorage) SetMetrics(ctx context.Context, metrics map[storage.MetricType]map[string]interface{}) {
	// Начало транзакции
	tx, err := p.db.Begin()
	if err != nil {
		panic(fmt.Errorf("ошибка при начале транзакции: %v", err))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Ошибка отката транзакции:", rollbackErr)
			}
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

func (p *PostgresStorage) UpdateJSONMetrics(ctx context.Context, metrics *[]models.Metrics) error {
	// ctxB := context.Background()

	// Начинаем транзакцию
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Ошибка отката транзакции:", rollbackErr)
			}
		}
	}()

	// Создаем карты для хранения сумм дельт метрик типа "counter" и значений метрик типа "gauge"
	counterDeltas := make(map[string]int64)
	gaugeValues := make(map[string]float64)
	// Обновляем значения в карты в соответствии с типом метрики
	for _, metric := range *metrics {
		switch metric.MType {
		case "counter":
			counterDeltas[metric.ID] += *metric.Delta
		case "gauge":
			gaugeValues[metric.ID] = *metric.Value
		}
	}

	var queryString = "INSERT INTO %s (name, type, %[2]s)" +
		" VALUES ($1, '%[2]s', $2)" +
		" ON CONFLICT (name, type) DO UPDATE" +
		" SET %[2]s = EXCLUDED.%[2]s;"

	// Вставляем значения метрик типа "counter" из карты в базу данных
	for id, delta := range counterDeltas {
		fmt.Printf("id, value: %v, %v\n", id, delta)
		query := fmt.Sprintf(queryString, p.getTableName(), "counter")
		_, err := tx.ExecContext(ctx, query, id, delta)
		if err != nil {
			fmt.Printf("2 err.Error(): %v\n", err.Error())
			return err
		}
	}

	// Вставляем значения метрик типа "gauge" из карты в базу данных
	for id, value := range gaugeValues {
		fmt.Printf("id, value: %v, %v\n", id, value)
		query := fmt.Sprintf(queryString, p.getTableName(), "gauge")
		_, err := tx.ExecContext(ctx, query, id, value)
		if err != nil {
			fmt.Printf("1 err.Error(): %v\n", err.Error())
			return err
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		fmt.Printf("3 err.Error(): %v\n", err.Error())
		return err
	}

	return nil
}
