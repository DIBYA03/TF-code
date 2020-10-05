package shared

type StreamProviderType string

func (s StreamProviderType) String() string {
	return string(s)
}

const (
	StreamProviderTypeKinesis = StreamProviderType("kinesis")
	StreamProviderTypeKafka   = StreamProviderType("kafka")
)

type StreamProviderRegion string

func (s StreamProviderRegion) String() string {
	return string(s)
}

type StreamProvider interface {
	Type() StreamProviderType
	Region() StreamProviderRegion
	StreamName() string
}

func NewKinesisStreamProvider(region StreamProviderRegion, streamName string) StreamProvider {
	return &kinesisStreamProvider{region, streamName}
}

type kinesisStreamProvider struct {
	region     StreamProviderRegion
	streamName string
}

func (p *kinesisStreamProvider) Type() StreamProviderType {
	return StreamProviderTypeKinesis
}

func (p *kinesisStreamProvider) Region() StreamProviderRegion {
	return p.region
}

func (p *kinesisStreamProvider) StreamName() string {
	return p.streamName
}
