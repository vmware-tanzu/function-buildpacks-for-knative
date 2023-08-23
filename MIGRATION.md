# Migrating Functions

The buildpacks we offer in this repository are going out of support and being archived.

This might put you in a situation where you need to migrate to a supported version. We are going to make our best effort to help you make this decision but this decisions are not going to be supported by the team; just our opinions.

## Function to Function Buildpacks

This repository is just one of the list of trusted builders by the Knatinve Function CLI. You can find more function buildpacks on https://github.com/boson-project/buildpacks/tree/main/buildpacks.

## Function to Languague Buildpack

If you want to delete your dependency on the function invoker you might need to do a few things depending on the language:

### Python

The way our invoker is set up wraps your function in a [Flask](https://flask.palletsprojects.com/) [handler](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/invokers/python/pyfunc/invoke.py#L94) with a default set of configurations.

You can take the [PyFunc](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/invokers/python/pyfunc/) (our invoker) as a starting point of your project and add your business logic in the [handler function](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/invokers/python/pyfunc/invoke.py#L33). 

### Java

For a Function in our Java/Springboot offering we recommend adding the [Spring Cloud Function](https://spring.io/projects/spring-cloud-function/) into your favorite gradle or maven project and follow the Spring documentation to add your business logic.

### Buildpacks

Now with any of these changes you can use a more generic [Python](https://github.com/paketo-buildpacks/python) or [Java/Springboot](https://github.com/paketo-buildpacks/spring-boot) buildpack from the Paketo buildpacks list or any of your like.


I hope this guide had helped you understand the different paths for migration your function project outside the set of deprecated functions this repository offers. 
