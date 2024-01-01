package internal

import (
	"context"
	"strings"
)

func (s *SpotPriceExporter) RenderPrometheus(ctx context.Context) string {
	builder := strings.Builder{}

	for _, price := range s.GetPriceHistory(ctx) {
		builder.WriteString(s.config.MetricName)
		builder.WriteString("{")
		builder.WriteString("region=\"")
		builder.WriteString(price.Region)
		builder.WriteString("\",")
		builder.WriteString("instance_type=\"")
		builder.WriteString(price.InstanceType)
		builder.WriteString("\",")
		builder.WriteString("product_description=\"")
		builder.WriteString(price.ProductDescription)
		builder.WriteString("\",")
		builder.WriteString("availability_zone=\"")
		builder.WriteString(price.AvailabilityZone)
		builder.WriteString("\"")

		if s.config.MapZoneId {
			builder.WriteString(",zone_id=\"")
			builder.WriteString(s.zoneIdMap[price.Region][price.AvailabilityZone])
			builder.WriteString("\"")
		}

		builder.WriteString("} ")
		builder.WriteString(price.SpotPrice)
		builder.WriteString("\n")
	}

	return builder.String()
}
