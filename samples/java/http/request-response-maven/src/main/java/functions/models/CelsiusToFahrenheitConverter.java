/*
 * Copyright 2021-2022 VMware, Inc.
 * SPDX-License-Identifier: BSD-2-Clause
 */

package functions.models;

import java.math.BigDecimal;

public class CelsiusToFahrenheitConverter {

    public BigDecimal celsius;

    public BigDecimal getFahrenheit() {
        BigDecimal nine = new BigDecimal("9");
        BigDecimal five = new BigDecimal("5");
        BigDecimal nineOverFive = nine.divide(five);
        BigDecimal thirtyTwo = new BigDecimal("32");

        BigDecimal result = celsius.multiply(nineOverFive);
        return result.add(thirtyTwo);
    }

    public BigDecimal getCelsius() {
        return celsius;
    }

    public void setCelsius(BigDecimal celsius) {
        this.celsius = celsius;
    }

}
