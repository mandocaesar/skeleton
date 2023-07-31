package domain

import (
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	samplereserve "github.com/machtwatch/catalyst-go-skeleton/domain/sample_reserve"
)

type RepositoryCollection struct {
	SampleRepo        sample.SampleRepo
	SampleReserveRepo samplereserve.SampleReserveRepo
}
