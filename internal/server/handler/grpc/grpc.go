package grpc

import (
	"context"
	"fmt"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/grpclog"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	pb "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/protometrics"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/validate"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServiceServer
	s   storage.MetricsStorage
	cfg *flags.InitedFlags
}

// NewHandler создает экземпляр Handler
func NewServer(mStorage storage.MetricsStorage, cnf *flags.InitedFlags) *MetricsServer {
	return &MetricsServer{
		s:   mStorage,
		cfg: cnf,
	}
}

func (s *MetricsServer) UpdateJSONMetrics(ctx context.Context, req *pb.UpdateJSONMetricsRequest) (*pb.UpdateJSONMetricsResponse, error) {
	md := metadata.Pairs(
		"content-type", "application/json",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Конвертируем gRPC MetricsList в формат, используемый в вашем хранилище
	var metrics []models.Metrics
	for _, m := range req.GetMetrics() {
		delta := m.GetDelta()
		value := m.GetValue()
		metric := models.Metrics{
			ID:    m.GetId(),
			MType: m.GetType(),
			Delta: &delta,
			Value: &value,
		}
		metrics = append(metrics, metric)
	}

	// Обновляем метрики
	err := s.s.UpdateJSONMetrics(ctx, &metrics)
	if err != nil {
		return nil, fmt.Errorf("error updating metrics: %v", err)
	}

	// Формируем ответ в формате MetricsList
	updatedMetrics := &pb.UpdateJSONMetricsResponse{}
	for _, m := range metrics {
		updatedMetric := &pb.Metric{
			Id:    m.ID,
			Type:  m.MType,
			Delta: *m.Delta,
			Value: *m.Value,
		}
		updatedMetrics.Metrics = append(updatedMetrics.Metrics, updatedMetric)
	}

	return updatedMetrics, nil
}

func (s *MetricsServer) UpdateJSONMetric(ctx context.Context, req *pb.Metric) (*pb.Metric, error) {
	// Преобразование *pb.Metric в *models.Metrics
	metric := models.Metrics{
		ID:    req.Id,
		MType: req.Type,
		Delta: &req.Delta,
		Value: &req.Value,
	}

	// Обновляем метрику
	err := s.s.UpdateJSONMetric(ctx, &metric)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update metric: %v", err)
	}

	// Получаем метрику после обновления
	err = s.s.GetJSONMetric(ctx, &metric)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get updated metric: %v", err)
	}

	// Преобразование *models.Metrics в *pb.Metric
	resp := &pb.Metric{
		Id:    metric.ID,
		Type:  metric.MType,
		Delta: *metric.Delta,
		Value: *metric.Value,
	}

	// Возвращаем обновленную метрику в качестве ответа
	return resp, nil
}

func (s *MetricsServer) GetJSONMetrics(ctx context.Context, req *pb.Metric) (*pb.GetJSONMetricsResponse, error) {
	// metrics := s.s.GetMetrics(ctx)
	var metricsData []*pb.Metric
	// for name, value := range metricsMap {
	//     // Преобразуем каждый элемент в pb.Metric и добавляем в срез
	//     metricsData = append(metricsData, &pb.Metric{
	//         Name:  name,
	//         Value: value,
	//     })
	// }

	return &pb.GetJSONMetricsResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) GetMetric(ctx context.Context, req *pb.Metric) (*pb.GetMetricResponse, error) {
	// Валидация входящих данных
	if err := validate.CheckMetricTypeAndName(req.GetType(), req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid metric type or name: %v", err)
	}

	// Получение данных для вывода
	metricTypeID, err := utils.GetMetricTypeByCode(req.GetType())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get metric type by code: %v", err)
	}

	currentMetrics := s.s.GetMetric(ctx, metricTypeID)

	value, ok := currentMetrics[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Metric not found")
	}

	return &pb.GetMetricResponse{Value: fmt.Sprintf("%v", value)}, nil
}

func (s *MetricsServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	// Здесь мы получаем все метрики
	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.GetMetricsResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) GetMetricsRoot(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	// Здесь мы получаем все метрики
	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.GetMetricsResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) UpdateMetric(ctx context.Context, req *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	// Проверка типа и имени метрики
	if err := validate.CheckMetricTypeAndName(req.GetType(), req.GetName()); err != nil {
		return nil, fmt.Errorf("invalid metric type or name: %w", err)
	}

	err := s.s.UpdateMetric(ctx, req.GetType(), req.GetName(), req.GetValue())
	if err != nil {
		return nil, fmt.Errorf("failed to update metric: %w", err)
	}

	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.UpdateMetricResponse{Metrics: metricsData}, nil
}

func convertMetrics(metrics map[storage.MetricType]map[string]interface{}) []*pb.Metric {
	var metricsData []*pb.Metric
	for mtype, data := range metrics {
		// Преобразуем каждый элемент в pb.Metric и добавляем в срез
		for id, value := range data {
			metric := pb.Metric{
				Type:  string(mtype),
				Id:    id,
				Value: 0,
				Delta: 0,
			}
			// Пробуем преобразовать значение метрики
			switch v := value.(type) {
			case int64:
				// Если значение int64, присваиваем его в Delta
				metric.Delta = v
			case float64:
				// Если значение float64, присваиваем его в Value
				metric.Value = v
			}

			metricsData = append(metricsData, &metric)
		}
	}
	return metricsData
}

func (s *MetricsServer) PingDB(ctx context.Context, req *pb.PingDBRequest) (*pb.PingDBResponse, error) {
	err := s.s.Ping()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	return &pb.PingDBResponse{Message: "Database connection is successful!"}, nil
}
