package com.vmware.functions;

import java.util.Map;
import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.messaging.Message;
import org.springframework.cloud.function.cloudevent.CloudEventMessageUtils;

@SpringBootApplication
public class Func {
	public static void main(String[] args) {
		SpringApplication.run(Func.class, args);
	}

	@Bean
	public Function<Message<String>, String> echo() {
		return event -> {
			if (!CloudEventMessageUtils.isCloudEvent(event)) {
				return "Did not receive cloudevent";
			}

			String payload = CloudEventMessageUtils.getData(event);
			// Map<String, Object> attrs = CloudEventMessageUtils.getAttributes(event);
			// Object attr = attrs.get("ce-my-attr");
			//
			// if (null != attr && "test" != attr) {
			// 	return "Invalid attribute value for 'my-attr'";
			// }
			return payload;
		};
	}
}
