package app

import "go.uber.org/zap"

func setLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)

	}
	defer logger.Sync()

	Sugar = *logger.Sugar()
}
