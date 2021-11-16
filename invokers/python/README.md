# PyFunc

This package provides a mechanism to write Function-as-a-Service style code in
Python for handling HTTP events, including CloudEvents delivered via HTTP.

This framework is primarily intended to work with
[Knative](https://knative.dev/), but also works to provide a generic server for
handling CloudEvents over HTTP (e.g. from Kubernetes or on a local machine).

The framework uses reflection to find a suitable function to wrap; it should not
be necessary to import any of the following modules in your own code unless you
want (e.g. for type definitions):

- `pyfunc` (this module; on PyPi as `pyfunc-invoker`)
- `flask`
- `cloudevents`

Instead, simply ensure that you have a single non-`_` prefixed function which
uses some combination of the following:

- HTTP request arguments (named `req`, `request`, `body`, `headers` or of the
  `flask.Request` type)
- CloudEvent arguments (named `event`, `payload`, `data`, `attributes` or of the
  `cloudevents.sdk.event.v1.Event` type)

Usage:

```python
import logging
from typing import Any

counter = 0

def handler(data: Any, attributes: dict, req: Any):
    global counter
    counter = counter + 1

    logging.info(f"Got data: {data}")
    logging.info(f"From {req.origin}, my {counter}th request!")

    attributes["type"] = "com.example.reply"
    attributes["datacontenttype"] = "text/plain"

    return attributes, "It's a demo"

```

### Usage

To check the current working directory for a module called `handler` and function `handler`:
```
python -m pyfunc
```

To check the current working directory for a module called `myhandler` and function `func`:
```
PYTHON_HANDLER=myhandler.func python -m pyfunc
```

To check a specific directory `./path/to/function/dir` for a module called `handler` and function `handler`:
```
python -m pyfunc ./path/to/function/dir
```

To check a specific directory `./path/to/function/dir` for a module called `mymodule` and function `fname`:
```
PYTHON_HANDLER=mymodule.fname python -m pyfunc ./path/to/function/dir
```
