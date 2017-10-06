package har

import "time"

type Har struct {
	log struct {
		creator struct {
			name    string
			version string
			comment string
		}
		pages   []*Page
		entries []*Entry
	}
}

type Value struct {
	name  string
	value string
}

type CacheItem struct {
	expires    *time.Time
	lastAccess *time.Time
	eTag       string
	hitCount   int64
	comment    string
}

type Cookie struct {
	name     string
	value    string
	path     string
	domain   string
	expires  *time.Time
	httpOnly bool
	secure   bool
	comment  string
}

type Entry struct {
	pageref         string
	startedDateTime *time.Time
	time            int64
	request         struct {
		method      string
		url         string
		httpVersion string
		cookies     []*Cookie
		headers     []*Value
		queryString []*Value
		postData    *struct {
			mimeType string
			params   []struct {
				name        string
				value       string
				fileName    string
				contentType string
				comment     string
			}
			text    string
			comment string
		}
		headersSize int64
		bodySize    int64
		comment     string
	}
	response struct {
		status      int64
		statusText  string
		httpVersion string
		cookies     []*Cookie
		headers     []*Value
		content     *struct {
			size        int64
			compression int64
			mimeType    string
			text        string
			comment     string
		}
		redirectURL string
		headersSize int64
		bodySize    int64
		comment     string
	}
	cache struct {
		beforeRequest CacheItem
		afterRequest  CacheItem
	}
	timings struct {
		blocked int64
		dns     int64
		connect int64
		send    int64
		wait    int64
		receive int64
		ssl     int64
		comment string
	}
	serverIPAddress string
	connection      string
	comment         string
}

type Page struct {
	startedDateTime *time.Time
	id              string
	title           string
	pageTimings     struct {
		onContentLoad int64
		onLoad        int64
		comment       string
	}
	comment string
}
