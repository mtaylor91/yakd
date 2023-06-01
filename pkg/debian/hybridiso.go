package debian

import (
	"context"

	"github.com/mtaylor91/yakd/pkg/util/executor"
)

type HybridISOSourceBuilder struct {
	FSDir  string
	ISODir string
}

func (b *HybridISOSourceBuilder) BuildISOFS(
	ctx context.Context, chroot executor.Executor,
) error {
	panic("not implemented")
}

func (b *HybridISOSourceBuilder) BuildISOSources(ctx context.Context) error {
	panic("not implemented")
}

type HybridISOBuilder struct {
	ISODir string
	Target string
}

func (b *HybridISOBuilder) BuildISO(ctx context.Context) error {
	panic("not implemented")
}
