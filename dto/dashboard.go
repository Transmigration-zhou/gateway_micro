package dto

type PanelGroupDataOutput struct {
	ServiceNum      int64 `json:"service_num"`
	TenantNum       int64 `json:"tenant_num"`
	CurrentQPS      int64 `json:"current_qps"`
	TodayRequestNum int64 `json:"today_request_num"`
}

type DashboardStatisticsOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" validate:"required"`         //今日流量
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" validate:"required"` //昨日流量
}

type DashServiceStatisticsItemOutput struct {
	Name     string `json:"name"`
	LoadType int    `json:"load_type"`
	Value    int64  `json:"value"`
}

type DashServiceStatisticsOutput struct {
	Legend []string                          `json:"legend"`
	Series []DashServiceStatisticsItemOutput `json:"series"`
}
