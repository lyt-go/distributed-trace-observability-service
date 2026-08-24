package tracecollector

import (
	"context"

	"tracing/internal/chunkstream"
)

func Collect(ctx context.Context, chunks [][]string, failAt int) ([]string, error) {
	stream := chunkstream.Start(chunks, failAt)
	items := make([]string, 0)
	for {
		select {
		case item, ok := <-stream.Items:
			if !ok {
				return items, nil
			}
			items = append(items, item)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
