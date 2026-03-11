package main

import (
	"crypto/tls"
	"sort"
)

func mapTLSVersions(tlsVersions []string) []int {
	var versions []int
	for _, tlsVersion := range tlsVersions {
		switch tlsVersion {
		case "TLSv1.3":
			versions = append(versions, tls.VersionTLS13)
		case "TLSv1.2":
			versions = append(versions, tls.VersionTLS12)
		case "TLSv1.1":
			versions = append(versions, tls.VersionTLS11)
		case "TLSv1":
			versions = append(versions, tls.VersionTLS10)
		default:
			continue
		}
	}
	sort.Ints(versions)
	return versions
}
