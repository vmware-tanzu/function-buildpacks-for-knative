package org.example;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.function.Function;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.v1.CloudEventBuilder;

public class Greeter implements Function<CloudEvent, CloudEvent> {

   @Override
   public CloudEvent apply(CloudEvent cloudEvent) {
       System.out.println("Received CE:" + cloudEvent);

       CloudEvent ceToSend = new CloudEventBuilder()
               .withId("my-id")
               .withSource(URI.create("/my-test"))
               .withType("function-reply")
               .withDataContentType("application/json")
               .withData("{ \"msg\" : \"hello\" }".getBytes(StandardCharsets.UTF_8))
               .build();
       return ceToSend;
   }

}
