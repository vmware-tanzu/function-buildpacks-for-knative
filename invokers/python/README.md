# PyFunc

This package provides a mechanism to write Function-as-a-Service style code in
Python for handling HTTP events, including CloudEvents delivered via HTTP.

This framework is primarily intended to work with
[Knative](https://knative.dev/), but also works to provide a generic server for
handling CloudEvents over HTTP (e.g. from Kubernetes or on a local machine).

The framework uses reflection to find a suitable function to wrap; it should not
be necessary to import any of the following modules in your own code unless you
want (e.g. for type definitions):

- `pyfunc` (this module)
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

def main(data: Any, attributes: dict, req: Any):
    global counter
    counter = counter + 1

    logging.info(f"Got data: {data}")
    logging.info(f"From {req.origin}, my {counter}th request!")

    attributes["type"] = "com.example.reply"
    attributes["datacontenttype"] = "text/plain"

    return attributes, "It's a demo"

```

### Usage

To validate whether the function can be loaded see the following examples:

To check the current working directory for a module with default values (module=`func`, function=`main`):
```
python -m pyfunc check
```

To check for a different module, define an environment variable with the name `MODULE_NAME`.
To check for a different function, define an environment variable with the name `FUNCTION_NAME`.
For example, to check the current working directory with (module=`my_handler`, function=`my_func`):
```
MODULE_NAME=my_handler FUNCTION_NAME=my_func python -m pyfunc check
```

To check a different directory, the `-s <path_to_search>` flag can be used to specify the directory to search for the function:
```
python -m pyfunc check -s ./path/to/function/dir
```

To run the function, instead of using `check` as above, we will be using `start`:
```
python -m pyfunc start
```