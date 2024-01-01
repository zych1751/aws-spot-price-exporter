package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
	"sync"
)

type SpotPriceExporter struct {
	ec2Clients map[string]*ec2.Client
	zoneIdMap  map[string]map[string]string
	config     *Config
}

func NewSpotPriceExporter(cfg *Config) *SpotPriceExporter {
	defaultRegion := "us-east-1"

	if len(cfg.Regions) > 0 {
		defaultRegion = cfg.Regions[0]
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(defaultRegion))

	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := ec2.NewFromConfig(awsConfig)

	ec2Clients := make(map[string]*ec2.Client)
	targetRegions := loadRegions(client, cfg)

	log.Printf("target regions: %v", targetRegions)

	for _, region := range targetRegions {
		ec2Clients[region] = ec2.NewFromConfig(awsConfig, func(o *ec2.Options) {
			o.Region = region
		})
	}

	zoneIdMap := make(map[string]map[string]string)
	if cfg.MapZoneId {
		waitGroup := &sync.WaitGroup{}
		for _, region := range targetRegions {
			r := region
			waitGroup.Add(1)
			go func() {
				zoneIdMap[r] = loadZoneIdMap(ec2Clients[r], cfg, r)
				log.Printf("Load zone id map for region %s", r)
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()
	}

	return &SpotPriceExporter{
		ec2Clients: ec2Clients,
		zoneIdMap:  zoneIdMap,
		config:     cfg,
	}
}

func loadRegions(ec2Client *ec2.Client, config *Config) []string {
	if len(config.Regions) > 0 {
		return config.Regions
	}

	response, err := ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{
		AllRegions: &config.IncludeAllRegions,
	})

	if err != nil {
		log.Fatalf("unable to load region info, %v", err)
	}

	result := make([]string, 0)

	for _, region := range response.Regions {
		result = append(result, *region.RegionName)
	}

	return result
}

func loadZoneIdMap(ec2Client *ec2.Client, config *Config, region string) map[string]string {
	localZones, err := ec2Client.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: &config.IncludeAllZones,
	})

	if err != nil {
		log.Fatalf("unable to load zone info for region %s, %v", region, err)
	}

	result := make(map[string]string)

	for _, zone := range localZones.AvailabilityZones {
		result[*zone.ZoneName] = *zone.ZoneId
	}

	return result
}
