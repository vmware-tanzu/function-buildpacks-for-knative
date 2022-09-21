/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package com.vmware.functions;

import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import java.net.URI;
import java.net.http.HttpResponse;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;

@SpringBootApplication
public class Handler {

	public static void main(String[] args) {
		SpringApplication.run(Handler.class, args);
	}

	@Bean
	public Function<String, String> hello() {
		return in -> {
			HttpClient client = HttpClient.newHttpClient();
			// Replace with your request
			HttpRequest request = HttpRequest.newBuilder()
					.uri(URI.create("http://example.com"))
					.build();
	
			try {
				return client.send(request, HttpResponse.BodyHandlers.ofString()).body();
			} catch (Exception e) {
				throw new RuntimeException("Failed to fetch website");
			}
		};
	}

	@Bean
	public Function<String, String> bye() {
		return in -> {
			return "bye";
		};
	}
}
