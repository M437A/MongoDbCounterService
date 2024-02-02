package connections

import (
	"context"
)

func Init(ctx context.Context) {
	initEnvFile()
	connectToMongoDB(ctx)
}
