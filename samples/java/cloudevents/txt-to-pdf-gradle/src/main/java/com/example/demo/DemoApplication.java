/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package com.example.demo;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.function.Function;

import lombok.Data;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.function.cloudevent.CloudEventMessageBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.messaging.Message;

import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.auth.credentials.EnvironmentVariableCredentialsProvider;
import software.amazon.awssdk.core.ResponseBytes;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.model.*;
import software.amazon.awssdk.services.s3.S3Client;

@SpringBootApplication
public class DemoApplication {

	public static void main(String[] args) {
		SpringApplication.run(DemoApplication.class, args);
	}

	@Bean
	public Function<Message<S3Event>, Message<S3Event>> func() {
		return in -> {
			// Print payload
			System.out.println("Message:" + in);
			String cloudEventType = (String) in.getHeaders().get("ce-type");
			System.out.println("Type: " + cloudEventType);
			System.out.println("Payload:" + in.getPayload());

			// Configure S3 client
			Region region = Region.US_WEST_2;
			EnvironmentVariableCredentialsProvider credentialsProvider = EnvironmentVariableCredentialsProvider.create();
			S3Client s3Client = S3Client.builder()
				.region(region)
				.credentialsProvider(credentialsProvider)
				.build();

			// Parse CE to obtain S3 information
			S3Event s3Event = in.getPayload();
			String bucket = s3Event.s3.bucket.name;
			System.out.println("Bucket: " + bucket);
        	String key = s3Event.s3.object.key;
			System.out.println("Key: " + key);

			if (!key.endsWith(".txt")) {
				System.out.println("Key does not end in .txt, skipping");
				return CloudEventMessageBuilder.withData(in.getPayload()).setId("TEST")
					.setSource(URI.create("https://test.cnr")).setType("file-processed").build();
			}

			// Get an object and print its contents.
            System.out.println("Downloading an object");
			GetObjectRequest getObjectRequest = GetObjectRequest.builder()
				.bucket(bucket)
				.key(key)
				.build();
			ResponseBytes<GetObjectResponse> objectBytes = s3Client.getObjectAsBytes(getObjectRequest);
            byte[] data = objectBytes.asByteArray();

			// Convert byte array into string array
			String s = new String(data, StandardCharsets.UTF_8);
			String lines[] = s.split("\\r?\\n");
			System.out.println("Lines: " + lines);

			// Convert string array to formatted PDF
			String pdf = "";
			try {
				pdf = PDFBuilder.convert(lines);
			} catch(Exception e) {
				System.out.println("Something went wrong.");
			}

			// Upload PDF to S3
    		System.out.println("Uploading object...");
			String convertedKey = key.substring(0, key.length() - 3) + "pdf";
			PutObjectRequest putObjectRequest = PutObjectRequest.builder()
				.bucket(bucket)
				.key(convertedKey)
				.contentType("application/pdf")
				.build();
			
			s3Client.putObject(putObjectRequest, RequestBody.fromString(pdf));

			// Exit out
			System.out.println("Upload complete");
			System.out.printf("%n");
			System.out.println(" in bucket " + bucket + " converted to PDF " + key);

			return CloudEventMessageBuilder.withData(in.getPayload()).setId("TEST")
					.setSource(URI.create("https://test.cnr")).setType("file-processed").build();
		};
	}

	@Data
	static class S3Event {
		private String awsRegion;
		private String eventName;
		private String eventSource;
		private S3 s3;
	}

	@Data
	static class S3 {
		private Bucket bucket;
		private Object object;
	}

	@Data
	static class Bucket {
		private String name;
		private String arn;
	}

	@Data
	static class Object {
		private String key;
		private String eTag;
		private int size;
	}
}
