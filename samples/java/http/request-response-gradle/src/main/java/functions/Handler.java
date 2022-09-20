/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package functions;

import java.math.BigDecimal;
import java.util.function.Function;

import functions.models.CelsiusToFahrenheitConverter;

public class Handler implements Function<CelsiusToFahrenheitConverter, BigDecimal> {
    @Override
    public BigDecimal apply(CelsiusToFahrenheitConverter celsiusToFahrenheitConverter) {
        return celsiusToFahrenheitConverter.getFahrenheit();
    }
}
