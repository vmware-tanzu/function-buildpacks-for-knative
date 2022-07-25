# Building Function images for Spring Boot Applications

This guide will explain how to build functions using the Spring Boot Functions.

The Spring Boot ecosystem already have [Spring Cloud Functions](https://spring.io/projects/spring-cloud-function) to manage Functions. But, the function's experience we offer takes away the management of the Spring Boot server and allow you to focus on the business model of your applications.

## Implementing your function

You need to build an implementation of the `Function` interface. The Function Interface is a part of the `java.util.function` package which has been introduced since Java 8, to implement functional programming in Java.

The Function interface has a abstract method `apply()`. Applies this function to the given argument.

```java
R apply(T t)
```

It only takes one parameter which is the input of your function and the return type R of your function result.

```Java
public class MyFunction implements Function<Integer, Double> {

    @Override
    public Double apply(Integer n) {
        return n / 2.0;
    }
}
```

To understand how to configure and extend your application please visit the Spring Boot buildpack [documentation](../buildpacks/java/README.md).
