# Tests

To run the smoke tests, run `make smoke-tests` from this directory.

## Adding Tests

Navigate to the `Makefile` and ensure your changes are either under the `smoke-tests` or `template-tests` suite, or add your own as needed.

Each case (image) has its own corresponding Makefile that must exist for the test suite to properly pick up the folder, build an image, and test off of it. Please ensure this file exists and is properly named and edited.

## Adding Template Tests

By default, the Makefile navigates all directories in the `/cases` subdirectory and builds images off of them. The templates tests exist in `/cases/template-ce` and `/cases/template-http`.

When you add a template in the root `templates/`, also symlink that template into the appropriate test directory. For example, a template at `templates/java/cloudevents-gradle` should be symlinked to `tests/cases/template-ce/java-cloudevents-gradle` 

Symlink example: 
```bash
cd tests/cases/template-ce
ln -s ../../../templates/python/cloudevents/ python-cloudevents
```

Use this command from the root of the project to view the current symlinks: `ls -lR tests/cases/ | grep ^l`
