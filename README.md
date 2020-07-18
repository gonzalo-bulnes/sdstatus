## sdstatus

This is a PoC tool to check the status of [SecureDrop](https://securedrop.org) instances.

All the addresses are mentioned in the `sdonion.txt` file, one address each line.
I have a copy of the addresses from the SecureDrop [directory](https://securedrop.org/directory)
to test.


## System requirements

The tool assumes that Tor is running with a SOCKS proxy at port 9050 in the system.

### Development

The dependencies are managed using [Go Modules][gomod] (your Go version must support them).

  [gomod]: https://blog.golang.org/using-go-modules

### How to build?

Clone the repo and run `make build`.

## Output format

```
[
	{
		"Info": {
			"sd_version": "0.6",
			"gpg_fpr": "3392A1CE68FE779A95FCAF04EDA0FB6F53FA9093"
		},
		"Url": "m4hynbhhctdk27jr.onion",
		"Available": true
	},
	{
		"Info": {
			"sd_version": "0.6",
			"gpg_fpr": "7C24A77EED0D50838E3315BD7A38590B2996F0C2"
		},
		"Url": "ftugftwajmgsmoau.onion",
		"Available": true
	}
]
```

Currently by default the tool prints JSON output on the STDOUT. It is a list of dictionaries telling if the site is available, and
the SecureDrop version and also the journalist GPG key fingerprint.

If you pass `--csv` flag to the tool, then it will print output in CSV format. First the Onion address, then SecureDrop version,
and then the journalist GPG key fingerprint. If only the Onion address is printed along with `,,`, it means those
instances could not be reached (maybe they down).

Remember that the CSV formatted output will be printed on the STDOUT as the network calls return the results, in an asynchronous manner.
For the JSON output format, the tool waits for all of the network calls to return the results, and then prints them on the STDOUT at
the end.

License: GPLv3+
