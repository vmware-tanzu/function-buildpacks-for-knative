package com.vmware.functions;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.UUID;
import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;

@SpringBootApplication
public class Handler {
	public static void main(String[] args) {
		SpringApplication.run(Handler.class, args);
	}
    
    @Bean
    public Function<CloudEvent, CloudEvent> hello() {
    return event -> CloudEventBuilder.from(event)
        .withId(UUID.randomUUID().toString())
        .withSource(URI.create("/hello-world"))
        .withType("function-reply")
        .withData("{ \"msg\" : \"Hello World!\" }".getBytes(StandardCharsets.UTF_8))
        .build();
    }

}
