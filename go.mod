module github.com/hlandau/acmetool

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000
	github.com/fatih/color v1.9.0 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/hlandau/buildinfo v0.0.0-20161112115716-337a29b54997 // indirect
	github.com/hlandau/dexlogconfig v0.0.0-20161112114350-244f29bd2608
	github.com/hlandau/goutils v0.0.0-20160722130800-0cdb66aea5b8
	github.com/hlandau/xlog v1.0.0
	github.com/jmhodges/clock v0.0.0-20160418191101-880ee4c33548
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/ogier/pflag v0.0.2-0.20160129220114-45c278ab3607 // indirect
	github.com/peterhellberg/link v1.1.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/hlandau/acmeapi.v2 v2.0.1
	gopkg.in/hlandau/configurable.v1 v1.0.1 // indirect
	gopkg.in/hlandau/easyconfig.v1 v1.0.17
	gopkg.in/hlandau/service.v2 v2.0.16
	gopkg.in/hlandau/svcutils.v1 v1.0.10
	gopkg.in/square/go-jose.v1 v1.1.2
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/tylerb/graceful.v1 v1.2.15
	gopkg.in/yaml.v2 v2.3.0
)

exclude (
	github.com/btcsuite/winsvc v1.0.0
	github.com/erikdubbelboer/gspt v0.0.0-20190125194910-e68493906b83
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
