# ServiceName Feature - Implementation Summary

## Overview
Added support for an optional `ServiceName` configuration option that allows services to identify themselves in the User-Agent header of HTTP requests to ConfigCat.

## Changes Made

### 1. Config Struct (`configcat_client.go`)
Added new optional field:
```go
// ServiceName is an optional identifier for the service using the SDK.
// When set, it will be included in the User-Agent header of HTTP requests.
ServiceName string
```

### 2. ConfigFetcher Struct (`config_fetcher.go`)
- Added `serviceName` field to store the service name
- Updated `newConfigFetcher` to pass `cfg.ServiceName` to the fetcher
- Modified User-Agent header construction in `fetchHTTPWithoutRedirect`

### 3. User-Agent Format
**Without ServiceName:**
```
ConfigCat-Go/anativ-{pollingMode}-{version}
Example: ConfigCat-Go/anativ-a-9.0.7
```

**With ServiceName:**
```
ConfigCat-Go/{serviceName}-{pollingMode}-{version}
Example: ConfigCat-Go/my-service-a-9.0.7
```

Where:
- `a` = AutoPoll
- `l` = Lazy
- `m` = Manual

### 4. Tests (`config_fetcher_test.go`)
Created comprehensive test suite with 5 test cases covering:
- AutoPoll, Manual, and Lazy polling modes
- With and without ServiceName
- All tests passing ✅

## Usage Example

```go
client := configcat.NewCustomClient(configcat.Config{
    SDKKey:      "your-sdk-key",
    ServiceName: "my-backend-service",
})
```

This will send HTTP requests with:
```
X-ConfigCat-UserAgent: ConfigCat-Go/my-backend-service-a-9.0.7
```

## Benefits
- **Service Identification**: Easily identify which service is making requests
- **Debugging**: Trace requests back to specific services in distributed systems
- **Analytics**: Understand service usage patterns
- **Backward Compatible**: ServiceName is optional, existing code works unchanged

