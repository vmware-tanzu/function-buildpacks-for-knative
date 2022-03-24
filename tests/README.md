# Tests

To run the smoke tests, run `make smoke-tests` from this directory.

## Adding Tests

Navigate to the `Makefile` and ensure your changes are either under the `smoke-tests` or `template-tests` suite, or add your own as needed.

Each case (image) has its own corresponding Makefile that must exist for the test suite to properly pick up the folder, build an image, and test off of it. Please ensure this file exists and is properly named and edited.

## Adding Template Tests

By default, the Makefile navigates all directories in the `/cases` subdirectory and builds images off of them. The templates tests exist in `/cases/template-ce` and `/cases/template-http`, where the templates are copies of the root directory `templates/` folders' contents. 

Please keep in mind that if you add or edit a template in the root `templates/`, you will also need to copy the changes over to this directory. (In the future, if we do Symlinking, the same is true for renaming any folder.)
