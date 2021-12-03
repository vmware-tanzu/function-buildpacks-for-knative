package com.vmware.functions;

import java.util.Map;
import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.messaging.Message;
import org.springframework.cloud.function.cloudevent.CloudEventMessageUtils;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.function.Function;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.v1.CloudEventBuilder;

@SpringBootApplication
public class Handler {
	public static void main(String[] args) {
		SpringApplication.run(Handler.class, args);
	}

	@Bean
	public Function<CloudEvent, CloudEvent> hello() {
        CloudEvent outgoingCloudEvent = new CloudEventBuilder()
               .withId("my-id")
               .withSource(URI.create("/my-test"))
               .withType("function-reply")
               .withDataContentType("application/json")
               .withData("{ \"msg\" : \"hello\" }".getBytes(StandardCharsets.UTF_8))
               .build();
        return outgoingCloudEvent;
    }
}
