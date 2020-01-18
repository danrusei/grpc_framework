package grpcklog

import (
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc/codes"
)

//KlogLevel is the Klog Level type
type KlogLevel int

const (
	//InfoLog is info level
	InfoLog KlogLevel = iota
	//WarningLog is warn level
	WarningLog
	//ErrorLog is error level
	ErrorLog
	//FatalLog is fatal level
	FatalLog
)

var klogLevelName = []string{
	InfoLog:    "INFO",
	WarningLog: "WARNING",
	ErrorLog:   "ERROR",
	FatalLog:   "FATAL",
}

func (level KlogLevel) String() string {
	return klogLevelName[level]
}

var (
	defaultOptions = &options{
		levelFunc:    nil,
		shouldLog:    grpc_logging.DefaultDeciderMethod,
		codeFunc:     grpc_logging.DefaultErrorToCode,
		durationFunc: DefaultDurationToField,
	}
)

type options struct {
	levelFunc    CodeToLevel
	shouldLog    grpc_logging.Decider
	codeFunc     grpc_logging.ErrorToCode
	durationFunc DurationToField
}

func evaluateServerOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	optCopy.levelFunc = DefaultCodeToLevel
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

func evaluateClientOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	optCopy.levelFunc = DefaultClientCodeToLevel
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option for options
type Option func(*options)

// CodeToLevel function defines the mapping between gRPC return codes and interceptor log level.
type CodeToLevel func(code codes.Code) KlogLevel

// DurationToField function defines how to produce duration fields for logging
type DurationToField func(duration time.Duration) (key string, value interface{})

// WithDecider customizes the function for deciding if the gRPC interceptor logs should log.
func WithDecider(f grpc_logging.Decider) Option {
	return func(o *options) {
		o.shouldLog = f
	}
}

// WithLevels customizes the function for mapping gRPC return codes and interceptor log level statements.
func WithLevels(f CodeToLevel) Option {
	return func(o *options) {
		o.levelFunc = f
	}
}

// WithCodes customizes the function for mapping errors to error codes.
func WithCodes(f grpc_logging.ErrorToCode) Option {
	return func(o *options) {
		o.codeFunc = f
	}
}

// WithDurationField customizes the function for mapping request durations to Zap fields.
func WithDurationField(f DurationToField) Option {
	return func(o *options) {
		o.durationFunc = f
	}
}

// DefaultCodeToLevel is the default implementation of gRPC return codes to log levels for server side.
func DefaultCodeToLevel(code codes.Code) KlogLevel {
	switch code {
	case codes.OK:
		return InfoLog
	case codes.Canceled:
		return InfoLog
	case codes.Unknown:
		return ErrorLog
	case codes.InvalidArgument:
		return WarningLog
	case codes.DeadlineExceeded:
		return WarningLog
	case codes.NotFound:
		return InfoLog
	case codes.AlreadyExists:
		return InfoLog
	case codes.PermissionDenied:
		return WarningLog
	case codes.Unauthenticated:
		return InfoLog
	case codes.ResourceExhausted:
		return WarningLog
	case codes.FailedPrecondition:
		return WarningLog
	case codes.Aborted:
		return WarningLog
	case codes.OutOfRange:
		return WarningLog
	case codes.Unimplemented:
		return ErrorLog
	case codes.Internal:
		return ErrorLog
	case codes.Unavailable:
		return WarningLog
	case codes.DataLoss:
		return ErrorLog
	default:
		return ErrorLog
	}
}

// DefaultClientCodeToLevel is the default implementation of gRPC return codes to log levels for client side.
func DefaultClientCodeToLevel(code codes.Code) KlogLevel {
	switch code {
	case codes.OK:
		return InfoLog
	case codes.Canceled:
		return ErrorLog
	case codes.Unknown:
		return InfoLog
	case codes.InvalidArgument:
		return WarningLog
	case codes.DeadlineExceeded:
		return WarningLog
	case codes.NotFound:
		return ErrorLog
	case codes.AlreadyExists:
		return ErrorLog
	case codes.PermissionDenied:
		return WarningLog
	case codes.Unauthenticated:
		return InfoLog
	case codes.ResourceExhausted:
		return ErrorLog
	case codes.FailedPrecondition:
		return WarningLog
	case codes.Aborted:
		return WarningLog
	case codes.OutOfRange:
		return ErrorLog
	case codes.Unimplemented:
		return WarningLog
	case codes.Internal:
		return WarningLog
	case codes.Unavailable:
		return ErrorLog
	case codes.DataLoss:
		return WarningLog
	default:
		return InfoLog
	}
}

// DefaultDurationToField is the default implementation of converting request duration to a log field (key and value).
var DefaultDurationToField = DurationToTimeMillisField

// DurationToTimeMillisField converts the duration to milliseconds and uses the key `grpc.time_ms`.
func DurationToTimeMillisField(duration time.Duration) (key string, value interface{}) {
	return "grpc.time_ms", durationToMilliseconds(duration)
}

// DurationToDurationField uses the duration value to log the request duration.
func DurationToDurationField(duration time.Duration) (key string, value interface{}) {
	return "grpc.duration", duration
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}
