package grpc

import (
	"context"
	"fmt"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/grpclog"

	"github.com/Arcadian-Sky/musthave-metrics/internal/server/flags"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage/utils"

	// pb "github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/protometrics"
	pb "github.com/Arcadian-Sky/musthave-metrics/gen/proto/api/metrics/v1"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/handler/validate"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/models"
	"github.com/Arcadian-Sky/musthave-metrics/internal/server/storage"
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
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}
	md := metadata.Pairs(
		"content-type", "application/json",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Конвертируем gRPC MetricsList в формат, используемый в вашем хранилище
	var metrics []models.Metrics
	for _, m := range req.GetMetrics() {
		delta := m.GetDelta()
		value := m.GetValue()
		fmt.Printf("m.GetType(): %v\n", m.GetType())
		fmt.Printf("m: %v\n", m)
		metric := models.Metrics{
			ID:    m.GetId(),
			MType: convertMetricTypeToString(m.GetType()),
			Delta: &delta,
			Value: &value,
		}
		metrics = append(metrics, metric)
	}

	fmt.Printf("metrics: %v\n", metrics)
	// Обновляем метрики
	err := s.s.UpdateJSONMetrics(ctx, &metrics)
	if err != nil {
		return nil, fmt.Errorf("error updating metrics: %v", err)
	}

	// Формируем ответ в формате MetricsList
	updatedMetrics := &pb.UpdateJSONMetricsResponse{}
	for _, m := range metrics {

		metricType := pb.Type_TYPE_UNSPECIFIED

		updatedMetric := pb.Metric{
			Id:   m.ID,
			Type: metricType,
		}

		switch m.MType {
		case "counter":
			updatedMetric.Type = pb.Type_TYPE_COUNTER
			updatedMetric.Mvalue = &pb.Metric_Delta{
				Delta: *m.Delta,
			}
		case "gauge":
			updatedMetric.Type = pb.Type_TYPE_GAUGE
			updatedMetric.Mvalue = &pb.Metric_Value{
				Value: *m.Value,
			}
		default:
			fmt.Println("Unknown MetricType")
		}

		updatedMetrics.Metrics = append(updatedMetrics.Metrics, &updatedMetric)
	}

	return updatedMetrics, nil
}

func (s *MetricsServer) UpdateJSONMetric(ctx context.Context, req *pb.UpdateJSONMetricRequest) (*pb.UpdateJSONMetricResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	if req.GetMetric() == nil {
		return nil, status.Error(codes.InvalidArgument, "Metric field is required")
	}
	rmetric := req.GetMetric()
	delta := rmetric.GetDelta()
	value := rmetric.GetValue()
	// Преобразование *pb.Metric в *models.Metrics
	metric := models.Metrics{
		ID:    rmetric.Id,
		MType: convertMetricTypeToString(rmetric.GetType()),
		Delta: &delta,
		Value: &value,
	}

	fmt.Printf("metric: %v\n", metric)
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
	resp := convertModelMetricsToRPCMetrics(&metric)

	// Возвращаем обновленную метрику в качестве ответа
	return &pb.UpdateJSONMetricResponse{Metric: resp}, nil
}

func (s *MetricsServer) GetJSONMetrics(ctx context.Context, req *pb.GetJSONMetricsRequest) (*pb.GetJSONMetricsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}
	metricsMap := s.s.GetMetrics(ctx)
	metricsData := convertMetrics(metricsMap)
	return &pb.GetJSONMetricsResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) GetMetric(ctx context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Валидация входящих данных
	if err := validate.CheckMetricTypeAndName(convertMetricTypeToString(req.GetType()), req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid metric type or name: %v", err)
	}

	// Получение данных для вывода
	metricTypeID, err := utils.GetMetricTypeByCode(convertMetricTypeToString(req.GetType()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get metric type by code: %v", err)
	}

	currentMetrics := s.s.GetMetric(ctx, metricTypeID)

	value, ok := currentMetrics[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Metric not found")
	}
	resp := pb.GetMetricResponse{}
	switch string(metricTypeID) {
	case "counter":
		resp.Mvalue = &pb.GetMetricResponse_Delta{
			Delta: value.(int64),
		}
	case "gauge":
		resp.Mvalue = &pb.GetMetricResponse_Value{
			Value: value.(float64),
		}
	default:
		fmt.Println("Unknown MetricType")
	}

	fmt.Printf("value: %v\n", value)
	return &resp, nil
}

func (s *MetricsServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Здесь мы получаем все метрики
	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.GetMetricsResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) GetMetricsRoot(ctx context.Context, req *pb.GetMetricsRootRequest) (*pb.GetMetricsRootResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Здесь мы получаем все метрики
	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.GetMetricsRootResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) UpdateMetric(ctx context.Context, req *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Проверка типа и имени метрики
	if err := validate.CheckMetricTypeAndName(convertMetricTypeToString(req.GetType()), req.GetId()); err != nil {
		return nil, fmt.Errorf("invalid metric type or name: %w", err)
	}

	err := s.s.UpdateMetric(ctx, convertMetricTypeToString(req.GetType()), req.GetId(), req.GetValue())
	if err != nil {
		return nil, fmt.Errorf("failed to update metric: %w", err)
	}

	metrics := s.s.GetMetrics(ctx)
	var metricsData = convertMetrics(metrics)

	return &pb.UpdateMetricResponse{Metrics: metricsData}, nil
}

func (s *MetricsServer) PingDB(ctx context.Context, req *pb.PingDBRequest) (*pb.PingDBResponse, error) {
	err := s.s.Ping()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	return &pb.PingDBResponse{Message: "Database connection is successful!"}, nil
}

func convertMetrics(metrics map[storage.MetricType]map[string]interface{}) []*pb.Metric {
	var metricsData []*pb.Metric
	for mtype, data := range metrics {
		// Преобразуем каждый элемент в pb.Metric и добавляем в срез
		for id, value := range data {

			metric := pb.Metric{
				Id:   id,
				Type: pb.Type_TYPE_UNSPECIFIED,
			}
			switch string(mtype) {
			case "counter":
				metric.Type = pb.Type_TYPE_COUNTER
				metric.Mvalue = &pb.Metric_Delta{
					Delta: value.(int64),
				}
			case "gauge":
				metric.Type = pb.Type_TYPE_GAUGE
				metric.Mvalue = &pb.Metric_Value{
					Value: value.(float64),
				}
			default:
				fmt.Println("Unknown MetricType")
			}

			metricsData = append(metricsData, &metric)
		}
	}
	return metricsData
}

func convertModelMetricsToRPCMetrics(metric *models.Metrics) *pb.Metric {

	resp := pb.Metric{
		Id:   metric.ID,
		Type: convertStringToMetricType(metric.MType),
	}

	switch string(metric.MType) {
	case "counter":
		m := pb.Metric_Delta{}
		if metric.Delta != nil {
			m.Delta = *metric.Delta
		}
		resp.Mvalue = &m
	case "gauge":
		m := pb.Metric_Value{}
		if metric.Value != nil {
			m.Value = *metric.Value
		}
		resp.Mvalue = &m
		// default:
		// fmt.Println("Unknown MetricType")
	}

	return &resp
}

func convertStringToMetricType(s string) pb.Type {
	switch s {
	case string(storage.Gauge):
		return pb.Type_TYPE_GAUGE
	case string(storage.Counter):
		return pb.Type_TYPE_COUNTER
	default:
		return pb.Type_TYPE_UNSPECIFIED
	}
}

func convertMetricTypeToString(s pb.Type) string {
	switch s.Number() {
	case pb.Type_TYPE_GAUGE.Number():
		return string(storage.Gauge)
	case pb.Type_TYPE_COUNTER.Number():
		return string(storage.Counter)
	default:
		return ""
	}
}
