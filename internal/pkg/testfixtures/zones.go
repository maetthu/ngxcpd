package testfixtures

import "github.com/maetthu/ngxcpd/pkg/proxycache"

// TestdataCacheFiles is a collection of test fixtures for different cache zones
var TestdataCacheFiles = map[string][]*proxycache.Entry{
	"zone1": TestdataCacheFilesZone1,
	"zone2": TestdataCacheFilesZone2,
}
