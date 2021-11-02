package java

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

type Detect struct{}

func (Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	result := libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "kn-fn-java"},
					{Name: "jvm-application"},
				},
				Requires: []libcnb.BuildPlanRequire{
					// {Name: "jdk" },
					{Name: "jre", Metadata: map[string]interface{}{"launch": true}},
					{Name: "jvm-application"}, //jvm-application-package
				},
			},
		},
	}

	fmt.Println("SWAP: in detect, context ", context)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	//fmt.Println("SWAP: ConfigResolver", cr)
	cr.Resolve("KN_FN")

	// if ok, err := libfnbuildpack.IsRiff(context.Application.Path, cr); err != nil {
	// 	return libcnb.DetectResult{}, fmt.Errorf("unable to determine if application is riff\n%w", err)
	// } else if !ok {
	// 	return result, nil
	// }

	// metadata, err := libfnbuildpack.Metadata(context.Application.Path, cr)
	// if err != nil {
	// 	return libcnb.DetectResult{}, fmt.Errorf("unable to read riff metadata\n%w", err)
	// }

	result.Plans[0].Requires = append(result.Plans[0].Requires, libcnb.BuildPlanRequire{
		Name: "kn-fn-java",
		// Metadata: metadata,
	})

	fmt.Println("SWAP: result: ", result)
	return result, nil
}
