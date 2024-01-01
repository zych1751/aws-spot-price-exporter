package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/exp/maps"
	"log"
	"sync"
	"time"
)

type SpotPrice struct {
	Region             string
	AvailabilityZone   string
	InstanceType       string
	ProductDescription string
	SpotPrice          string
}

func (s *SpotPriceExporter) GetPriceHistory(ctx context.Context) []SpotPrice {
	result := make([]SpotPrice, 0)

	waitGroup := &sync.WaitGroup{}
	for region := range s.ec2Clients {
		r := region
		waitGroup.Add(1)
		go func() {
			prices, err := s.getPriceHistory(ctx, r)
			if err != nil {
				log.Printf("unable to get spot price history, %v", err)
			}

			result = append(result, prices...)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()

	return result
}

func (s *SpotPriceExporter) getPriceHistory(ctx context.Context, region string) ([]SpotPrice, error) {
	ec2Client := s.ec2Clients[region]

	result := make([]SpotPrice, 0)

	var currentTime = time.Now()
	var nextToken *string = nil

	for {
		response, err := ec2Client.DescribeSpotPriceHistory(ctx, &ec2.DescribeSpotPriceHistoryInput{
			StartTime: &currentTime,
			NextToken: nextToken,
			Filters: []types.Filter{
				{
					Name:   &[]string{"availability-zone"}[0],
					Values: maps.Keys[map[string]string](s.zoneIdMap[region]),
				},
			},
			MaxResults: &[]int32{100000}[0],
		})

		if err != nil {
			return nil, err
		}

		for _, price := range response.SpotPriceHistory {
			result = append(result, SpotPrice{
				Region:             region,
				AvailabilityZone:   *price.AvailabilityZone,
				InstanceType:       string(price.InstanceType),
				ProductDescription: string(price.ProductDescription),
				SpotPrice:          *price.SpotPrice,
			})
		}

		if response.NextToken == nil || len(*response.NextToken) == 0 {
			break
		}

		nextToken = response.NextToken
	}

	return result, nil
}
