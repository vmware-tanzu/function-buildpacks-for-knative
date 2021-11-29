// package com.vmware.functions;

// import java.util.function.Function;

// import io.cloudevents.CloudEvent;

// public class Echo implements Function<CloudEvent, String> {
// 	@Override
// 	public String apply(CloudEvent event) {
// 		String payload = event.getData().toString();
// 		String header = (String) event.getAttribute("my-header");
// 		if ("test" != header) {
// 			return "Incorrect 'my-header' value in request";
// 		}
// 		return payload;
// 	}
// }
