package functions;

import java.net.http.HttpResponse;
import java.util.function.Function;

public class HTTPFunc implements Function<String, String> {
    @Override
    public String apply(String s) {
        HttpClient client = HttpClient.newHttpClient();
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://example.com"))
                .build();

        HttpResponse<String> response = client.send(request,
                HttpResponse.BodyHandlers.ofString());

        System.out.println("HTTP GET request sent");
        return response.body();
    }
}
