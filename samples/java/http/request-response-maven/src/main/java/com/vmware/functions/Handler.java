/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package com.vmware.functions;

import java.util.function.Function;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import java.math.BigDecimal;

import com.vmware.functions.models.CelsiusToFahrenheitConverter;

@SpringBootApplication
public class Handler {

    public static void main(String[] args) {
        SpringApplication.run(Handler.class, args);
    }

    @Bean
    public Function<CelsiusToFahrenheitConverter, BigDecimal> convert() {
        return CelsiusToFahrenheitConverter::getFahrenheit;
    }

    @Bean
    public Function<String, String> uppercase() {
        return String::toUpperCase;
    }
}
