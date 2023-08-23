# Transitioning from Deprecated Function Buildpacks

The Function Buildpacks currently hosted in this repository are no longer supported and are being archived. To ensure the continued functionality of your applications, follow the migration strategies outlined below.

## Migration from Function Buildpacks

This repository is just one of the list of trusted builders by the Knatinve Function CLI. You can find more function buildpacks on https://github.com/boson-project/buildpacks/tree/main/buildpacks.

## Function to Languague Buildpack

If you are aiming to remove your reliance on the function invoker and migrate towards a more sustainable solution, the migration process will differ based on the programming language you are using.

### Python
For Python-based applications, transitioning can be achieved by adopting Flask, a popular web framework. The following steps are recommended:

1. **Develop Web App:** Create your web application using the Flask framework.
2. **Choose Production-ready Server:** Opt for a production-ready web server such as Waitress, Hypercorn, or Uvicorn to serve your Flask application efficiently.
3. **Cloudevents Integration:** Refer to the official [Cloudevents](https://cloudevents.io/) documentation for guidance on integrating Cloudevents into your Python application.

### Java
For Java-based applications, particularly those utilizing Spring Boot, the migration process involves incorporating [Spring Cloud Function](https://spring.io/projects/spring-cloud-function/) into your project. Follow these steps:

1. **Integrate Spring Cloud Function:** Add the Spring Cloud Function library to your project's Gradle or Maven configuration.
2. **Implement Business Logic:** Utilize Spring documentation to seamlessly integrate your existing business logic with Spring Cloud Function.

After transitioning your codebase to align with the recommended language-based solutions, you can benefit from more versatile buildpacks:

- **Python Buildpack:** If you migrated to Python with Flask, consider using the [Python Buildpack](https://github.com/paketo-buildpacks/python) from Paketo.
- **Java/Spring Boot Buildpack:** If you migrated to Java with Spring Boot, leverage the [Java/Spring Boot Buildpack](https://github.com/paketo-buildpacks/spring-boot) from Paketo.

I hope this guide had helped you understand the different paths for migration your function project outside the set of deprecated functions this repository offers. 
