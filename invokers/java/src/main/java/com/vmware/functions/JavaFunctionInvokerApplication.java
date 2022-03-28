/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package com.vmware.functions;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.function.Function;

@SpringBootApplication
public class JavaFunctionInvokerApplication {
	public static void main(String[] args) {
		SpringApplication.run(JavaFunctionInvokerApplication.class, args);
	}
}
