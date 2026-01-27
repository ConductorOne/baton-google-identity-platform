package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	sdkConfig "github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/ConductorOne/baton-google-identity-platform/pkg/config"
	"github.com/ConductorOne/baton-google-identity-platform/pkg/connector"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := sdkConfig.DefineConfiguration(
		ctx,
		"baton-google-identity-platform",
		getConnector,
		config.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.GoogleIdentityPlatform{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cfg *config.GoogleIdentityPlatform) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	cb, err := connector.New(ctx, cfg.CredentialsJSONFilePath)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
