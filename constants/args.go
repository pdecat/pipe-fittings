package constants

// Argument name constants
const (
	ArgVersion             = "version"
	ArgAll                 = "all"
	ArgBrowser             = "browser"
	ArgClear               = "clear"
	ArgCompact             = "compact"
	ArgServicePassword     = "database-password"
	ArgServiceShowPassword = "show-password"
	ArgSkipConfig          = "skip-config"
	ArgForeground          = "foreground"
	ArgInvoker             = "invoker"
	ArgSchemaComments      = "schema-comments"
	ArgCloudHost           = "cloud-host"
	ArgCloudToken          = "cloud-token"
	//nolint:gosec // This is not a hardcoded credential
	ArgDatabaseSSLPassword     = "database-ssl-password"
	ArgArg                     = "arg"
	ArgAutoComplete            = "auto-complete"
	ArgBaseUrl                 = "base-url"
	ArgBenchmarkTimeout        = "benchmark-timeout"
	ArgCacheMaxTtl             = "cache-max-ttl"
	ArgCacheTtl                = "cache-ttl"
	ArgClientCacheEnabled      = "client-cache-enabled"
	ArgConfigPath              = "config-path"
	ArgConnectionString        = "connection-string"
	ArgDashboardStartTimeout   = "dashboard-start-timeout"
	ArgDetectionTimeout        = "detection-timeout"
	ArgDashboardTimeout        = "dashboard-timeout"
	ArgDatabase                = "database"
	ArgDatabaseListenAddresses = "database-listen"
	ArgDatabasePort            = "database-port"
	ArgDatabaseQueryTimeout    = "query-timeout"
	ArgDatabaseStartTimeout    = "database-start-timeout"
	ArgDataDir                 = "data-dir"
	ArgDetach                  = "detach"
	ArgDisplayWidth            = "display-width"
	ArgDryRun                  = "dry-run"
	ArgEnvironment             = "environment"
	ArgExecutionId             = "execution-id"
	ArgExport                  = "export"
	ArgForce                   = "force"
	ArgFrom                    = "from"
	ArgTo                      = "to"
	ArgIndex                   = "index"
	ArgPartition               = "partition"
	ArgHeader                  = "header"
	ArgHelp                    = "help"
	ArgHost                    = "host"
	ArgInput                   = "input"
	ArgInsecure                = "insecure"
	ArgInstallDir              = "install-dir"
	ArgIntrospection           = "introspection"
	ArgLocal                   = "local"
	ArgListen                  = "listen"
	ArgLogLevel                = "log-level"
	ArgMaxCacheSizeMb          = "max-cache-size-mb"
	ArgMaxParallel             = "max-parallel"
	ArgMemoryMaxMb             = "memory-max-mb"
	ArgMemoryMaxMbPlugin       = "memory-max-mb-plugin"
	ArgModInstall              = "mod-install"
	ArgModLocation             = "mod-location"
	ArgMultiLine               = "multi-line"
	ArgOff                     = "off"
	ArgOn                      = "on"
	ArgOutput                  = "output"
	ArgPipesHost               = "pipes-host"
	ArgPipesInstallDir         = "pipes-install-dir"
	ArgPipesToken              = "pipes-token"
	ArgPluginStartTimeout      = "plugin-start-timeout"
	ArgPort                    = "port"
	ArgProgress                = "progress"
	ArgPrune                   = "prune"
	ArgPull                    = "pull"
	ArgRemote                  = "remote"
	ArgRemoteConnection        = "remote-connection"
	ArgSearchPath              = "search-path"
	ArgSearchPathPrefix        = "search-path-prefix"
	ArgSeparator               = "separator"
	ArgServiceCacheEnabled     = "service-cache-enabled"
	ArgShare                   = "share"
	ArgSnapshot                = "snapshot"
	ArgSnapshotLocation        = "snapshot-location"
	ArgSnapshotTag             = "snapshot-tag"
	ArgSnapshotTitle           = "snapshot-title"
	ArgTag                     = "tag"
	ArgTelemetry               = "telemetry"
	ArgTheme                   = "theme"
	ArgTiming                  = "timing"
	ArgUpdateCheck             = "update-check"
	ArgVarFile                 = "var-file"
	ArgVariable                = "var"
	ArgVerbose                 = "verbose"
	ArgWatch                   = "watch"
	ArgWhere                   = "where"
	ArgWorkspaceProfile        = "workspace"
	ArgWorkspaceDatabase       = "workspace-database"
	ArgResume                  = "resume"
	ArgResumeInput             = "resume-input"
	// Flowpipe concurrency
	ArgMaxConcurrencyHttp      = "max-concurrency-http"
	ArgMaxConcurrencyQuery     = "max-concurrency-query"
	ArgMaxConcurrencyContainer = "max-concurrency-container"
	ArgMaxConcurrencyFunction  = "max-concurrency-function"
	ArgProcessRetention        = "process-retention"
)

// BoolToOnOff converts a boolean value onto the string "on" or "off"
func BoolToOnOff(val bool) string {
	if val {
		return ArgOn
	}
	return ArgOff
}

// BoolToEnableDisable converts a boolean value onto the string "enable" or "disable"
func BoolToEnableDisable(val bool) string {
	if val {
		return "enable"
	}
	return "disable"

}
