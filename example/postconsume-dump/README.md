# Post-consumption hook logging debug information

This example is a [post-consumption hook][paperless-hooks] writing a bit debug
information to its output.

Paperless always appends a few arguments while all relevant information comes
from environment variables. A wrapper script is the easiest way to ignore the
additional arguments them:

```shell
#!/bin/bash

set -e -u -o pipefail

exec /usr/local/bin/postconsume-dump
```

Then configure Paperless to use the post-consumption hook:

```shell
PAPERLESS_POST_CONSUME_SCRIPT=/usr/local/hooks/postconsume
```

[paperless-hooks]: https://docs.paperless-ngx.com/advanced_usage/#consume-hooks

<!-- vim: set sw=2 sts=2 et : -->
