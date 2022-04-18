# Templates

Navigate to the subdirectory that contains the language, then event type, of your choice.

## Usage

Find the template of your need, modify the function(s) as needed, and follow the steps in `README.md` and `DEPLOYING.md` to experience your function.
## Testing

When you add a template in this directory, also add a symlink to the `tests/cases/template-<eventtype>` directory. The [README for tests](../tests/README.md) has details. 

To run the tests, you can navigate to `tests/` and run `make template-tests`, or push any changing commit and check GitHub Actions output if authenticated.
