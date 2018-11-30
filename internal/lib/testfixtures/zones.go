// +build fixtures

package testfixtures

import "github.com/maetthu/ngxcpd/internal/lib/proxycache"

var TestdataCacheFiles = map[string][]*proxycache.Entry{
	"zone1": TestdataCacheFilesZone1,
	"zone2": TestdataCacheFilesZone2,
}
