package com.vmware.functions;

import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

@SpringBootApplication
public class Handler {

	public static void main(String[] args) {
		SpringApplication.run(Handler.class, args);
	}

	@Bean
	public Function<String, String> hello() {
		return in -> "Hello World!";
	}
}
