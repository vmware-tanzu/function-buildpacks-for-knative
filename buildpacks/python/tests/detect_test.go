package tests

import (
	"bytes"
	"kn-fn/python-function-buildpack/python"
	"reflect"
	"testing"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
)

type result struct {
	plan libcnb.DetectResult
	err  error
}

func TestDetectNoEnvironmentWithValidFile(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withFunction("handler", "handler"))
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"handler": map[string]string{
				"module":   "handler",
				"function": "handler",
				"raw":      "handler.handler",
			},
		},
	}
	expectations.Check(t)
}

func TestDetectNoEnvironmentWithMissingFile(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withFunction("missing", "handler"))
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
	))

	expectations := DetectExpectations{
		Result:      result{plan, err},
		ExpectError: false,
		Pass:        false,
	}
	expectations.Check(t)
}

func TestDetectEnvironmentWithValidFile(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withFunction("other", "handler2"))
	defer app.Finish()

	plan, err := d.Detect(getContext(
		withApplicationPath(app.ApplicationPath),
		withEnvironmentVariable("PYTHON_HANDLER", "other.handler2"),
	))

	expectations := DetectExpectations{
		Result: result{plan, err},
		Pass:   true,
		Metadata: map[string]interface{}{
			"handler": map[string]string{
				"module":   "other",
				"function": "handler2",
				"raw":      "other.handler2",
			},
		},
	}
	expectations.Check(t)
}

func TestDetectEnironmentVariable(t *testing.T) {
	d := getDetector()
	app := createTestApplication(withFunction("handler", "handler"))
	defer app.Finish()

	cases := []struct {
		name    string
		context libcnb.DetectContext
	}{
		{
			name: "non matching module from environment",
			context: getContext(
				withApplicationPath(app.ApplicationPath),
				withEnvironmentVariable("PYTHON_HANDLER", "other.handler"),
			),
		},
		{
			name: "invalid environment variable (missing module)",
			context: getContext(
				withApplicationPath(app.ApplicationPath),
				withEnvironmentVariable("PYTHON_HANDLER", ".handler"),
			),
		},
		{
			name: "invalid environment variable (missing function)",
			context: getContext(
				withApplicationPath(app.ApplicationPath),
				withEnvironmentVariable("PYTHON_HANDLER", "handler."),
			),
		},
		{
			name: "too many items in environment variable",
			context: getContext(
				withApplicationPath(app.ApplicationPath),
				withEnvironmentVariable("PYTHON_HANDLER", "handler.handler.something"),
			),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			plan, err := d.Detect(c.context)

			expectations := DetectExpectations{
				Result:      result{plan, err},
				ExpectError: true,
			}
			expectations.Check(t)
		})
	}
}

func getDetector() python.Detect {
	buf := bytes.NewBuffer(nil)
	logger := bard.NewLogger(buf)
	return python.Detect{
		Logger: logger,
	}
}

func getContext(opts ...func(*libcnb.DetectContext)) libcnb.DetectContext {
	ctx := libcnb.DetectContext{
		Application: libcnb.Application{},
		Buildpack:   libcnb.Buildpack{},
		Platform: libcnb.Platform{
			Environment: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(&ctx)
	}

	return ctx
}

func withApplicationPath(path string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Application.Path = path
	}
}

func withEnvironmentVariable(key string, value string) func(*libcnb.DetectContext) {
	return func(dc *libcnb.DetectContext) {
		dc.Platform.Environment[key] = value
	}
}

type DetectExpectations struct {
	Result result

	ExpectError  bool
	Err          error
	Pass         bool
	Metadata     map[string]interface{}
	SkipProvides bool
	SkipRequires bool
}

func (e DetectExpectations) Check(t *testing.T) {
	if !e.ExpectError && e.Result.err != nil {
		if e.Err != nil && e.Result.err != e.Err {
			t.Errorf("expected error %v, but received error %v", e.Err, e.Result.err)
		} else {
			t.Errorf("unexpected error received: %v", e.Result.err)
		}
		return
	}

	if e.Pass != e.Result.plan.Pass {
		t.Errorf("expected detection to pass but it failed")
		return
	}

	// Find the invoker requires
	for _, plan := range e.Result.plan.Plans {
		for _, require := range plan.Requires {
			if require.Name == "kn-fn-python-invoker" {
				for k, v := range e.Metadata {
					val, ok := require.Metadata[k]
					if !ok {
						t.Errorf("plan requirement 'kn-fn-python-invoker' did not find metadata with key %s", k)
						return
					}

					if !reflect.DeepEqual(v, val) {
						t.Errorf("unexpected value in 'kn-fn-python-invoker' requires metadata for key %s", k)
						return
					}
				}
			}
		}
	}
}
