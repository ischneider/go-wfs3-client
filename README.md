go-wfs-client
=============

This package provides a POC client to the WFS-3 spec: https://rawgit.com/opengeospatial/WFS_FES/master/docs/17-069.html

Dependencies
------------

Listed in `vendor/vendor.json`. Use `govendor sync` or your favorite other
means of getting the dependencies.

NOTE: This effort was done with go-1.10 - it's not clear what the minimum
requirement is at this time.

Approach
--------

The client provides a way to connect to a WFS-3 service defined by a single
root URL. When connecting, a request to the /api path is made to obtain the
OpenAPI specification. The specification paths and operation IDs are used to
guide discovery. Operation parameters are adhered to as well but at the moment
there is only support for 'path' parameters.

NOTE: There is support in the client for the current path guidance
(as of 03-08-2018) and the newer proposal arrived at during the recent
hackathon. The CLI currently does not support switching between the two and
uses the older layout (where collectionInfo is at the root URL).
See: `wfs/paths.go`.

NOTE: If the specification does not parse correctly the first time, the
core WFS-3 'Components' section of an OpenAPI spec is 'patched' in to the
original definition and a 2nd attempt at parsing is made.

API/Spec Verification
---------------------

The intention of this effort was to include a means of validation and
verification. This currently only validates the OpenAPI specification using
a 3rd party library and does not contain WFS-3 specific checks.

Driving
-------

A simple CLI is provided as a driver for the client.

NOTE: The CLI uses a HTTP disk cache that, by default, will be created in
$TMP/wfs-http-cache

Examples:

Print out service summary (text):

    go run cmd/cli/main.go info <URL>

Print out collections metadata (raw response):

    go run cmd/cli/main.go coll <URL>

Invoke a simple operation (prints raw response):

    go run cmd/cli/main.go op <URL> <OPERATION>

Invoke an operation with parameter (prints raw response):

    go run cmd/cli/main.go op <URL> <OPERATION> <PARAM>=<VALUE>
