package handler

import (
	"fliqt/internal/util"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer(util.GetFileNameFromCaller())
