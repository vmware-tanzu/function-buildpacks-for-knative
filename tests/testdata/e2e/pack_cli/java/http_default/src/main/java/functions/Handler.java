/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package functions;

import java.util.function.Function;

public class Handler implements Function<String,String>{
	@Override
	public String apply(String in) {
		return "Hello World!";
	}
}
