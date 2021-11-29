package com.vmware.functions;

import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import io.cloudevents.CloudEvent;
import org.springframework.context.annotation.Bean;
// import org.springframework.messaging.Message;

@SpringBootApplication
public class Func {
	public static void main(String[] args) {
		SpringApplication.run(Func.class, args);
	}

	@Bean
	public Function<CloudEvent, String> echo() {
		return event -> {
			String payload = event.getData().toString();
			String header = (String) event.getAttribute("my-header");
			if ("test" != header) {
				return "Incorrect 'my-header' value in request";
			}
			return payload;
		};
	}

	// public String apply(Message<String> message) {
	// 	String payload = message.getPayload();
	// 	MessageHeaders headers = message.getHeaders();
	// 	String header = headers.get("my-header", String.class);
	// 	if ("test" != header) {
	// 		return "Incorrect 'my-header' value in request";
	// 	}
	// 	return payload;
	// }
}
