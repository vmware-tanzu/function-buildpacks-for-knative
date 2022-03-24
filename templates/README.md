# Templates

Navigate to the subdirectory that contains the language, then event type, of your choice.

## Usage

Find the template of your need, modify the function(s) as needed, and follow the steps in `README.md` and `DEPLOYING.md` to experience your function.
## Testing

If you add, edit, or rename a template in this directory, you will also need to manually copy the changes over to `tests/cases/template-<eventtype>` directory. In the future, we may symlink the folders.

To run the tests, you can navigate to `tests/` and run `make template-tests`, or push any changing commit and check GitHub Actions output if authenticated.
