/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package com.vmware.functions;

import java.util.function.Function;

public class Handler implements Function<Integer, Double> {

	@Override
	public Double apply(Integer input) {
		return input / 2.0;
	}
}
