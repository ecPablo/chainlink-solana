package monitoring

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type MultiFeedMonitor interface {
	Start(ctx context.Context, wg *sync.WaitGroup, feeds []FeedConfig)
}

func NewMultiFeedMonitor(
	solanaConfig SolanaConfig,

	log logger.Logger,
	sourceFactory SourceFactory,
	producer Producer,
	metrics Metrics,

	configSetTopic string,
	configSetSimplifiedTopic string,
	transmissionTopic string,

	configSetSchema Schema,
	configSetSimplifiedSchema Schema,
	transmissionSchema Schema,
) MultiFeedMonitor {
	return &multiFeedMonitor{
		solanaConfig,

		log,
		sourceFactory,
		producer,
		metrics,

		configSetTopic,
		configSetSimplifiedTopic,
		transmissionTopic,

		configSetSchema,
		configSetSimplifiedSchema,
		transmissionSchema,
	}
}

type multiFeedMonitor struct {
	solanaConfig SolanaConfig

	log           logger.Logger
	sourceFactory SourceFactory
	producer      Producer
	metrics       Metrics

	configSetTopic           string
	configSetSimplifiedTopic string
	transmissionTopic        string

	configSetSchema           Schema
	configSetSimplifiedSchema Schema
	transmissionSchema        Schema
}

const bufferCapacity = 100

// Start should be executed as a goroutine.
func (m *multiFeedMonitor) Start(ctx context.Context, wg *sync.WaitGroup, feeds []FeedConfig) {
	wg.Add(len(feeds))
	for _, feedConfig := range feeds {
		go func(feedConfig FeedConfig) {
			defer wg.Done()

			feedLogger := m.log.With(
				"feed", feedConfig.GetName(),
				"network", m.solanaConfig.NetworkName,
			)

			sources, err := m.sourceFactory.NewSources(m.solanaConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create new sources", "error", err)
				return
			}

			transmissionPoller := NewSourcePoller(
				sources.NewTransmissionsSource(),
				feedLogger.With("component", "transmissions-poller"),
				m.solanaConfig.PollInterval,
				m.solanaConfig.ReadTimeout,
				bufferCapacity,
			)
			statePoller := NewSourcePoller(
				sources.NewConfigSource(),
				feedLogger.With("component", "config-poller"),
				m.solanaConfig.PollInterval,
				m.solanaConfig.ReadTimeout,
				bufferCapacity,
			)

			wg.Add(2)
			go func() {
				defer wg.Done()
				transmissionPoller.Start(ctx)
			}()
			go func() {
				defer wg.Done()
				statePoller.Start(ctx)
			}()

			exporters := []Exporter{
				NewPrometheusExporter(
					m.solanaConfig,
					feedConfig,
					feedLogger.With("component", "prometheus-exporter"),
					m.metrics,
				),
				NewKafkaExporter(
					m.solanaConfig,
					feedConfig,
					feedLogger.With("component", "kafka-exporter"),
					m.producer,

					m.configSetSchema,
					m.configSetSimplifiedSchema,
					m.transmissionSchema,

					m.configSetTopic,
					m.configSetSimplifiedTopic,
					m.transmissionTopic,
				),
			}

			feedMonitor := NewFeedMonitor(
				feedLogger.With("component", "feed-monitor"),
				transmissionPoller, statePoller,
				exporters,
			)
			feedMonitor.Start(ctx, wg)
		}(feedConfig)
	}
}
