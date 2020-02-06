"""unit tests for InvoiceStats"""

import pandas as pd
import pytest

from stats import InvoiceStats


def test_InvoiceStats_instance_sets_default_values():
    instance = InvoiceStats()

    assert isinstance(instance.invoices, pd.DataFrame)
    assert instance.invoices.empty
    assert instance.maximum_invoices == 20_000_000


def test_add_invoice_adds_single_row_to_invoice_table():
    instance = InvoiceStats()
    length_before = len(instance.invoices)
    instance.add_invoice(1.56)
    length_after = len(instance.invoices)

    assert length_after - length_before == 1


def test_add_invoices_adds_multiple_rows_to_invoice_table():
    invoices = [1.21, 1.99, 202.98, 69909.01]
    instance = InvoiceStats()
    instance.add_invoices(invoices)

    assert len(instance.invoices) == len(invoices)


def test_clear_removes_all_stored_invoice_data():
    instance = InvoiceStats()
    instance.add_invoices([765.19, 11.01])
    length_before = len(instance.invoices)
    instance.clear()
    length_after = len(instance.invoices)

    assert length_before != 0
    assert length_after == 0


@pytest.mark.parametrize(
    "invoices,expected",
    [
        ([11.05, 11.06, 11.07, 11.08], 11.06),
        ([191.00, 150.00, 101.99], 150.00),
        ([1.00, 110.00, 110.02, 200_000.00], 110.01),
    ],
)
def test_get_median_calculates_median_and_rounds_down_half_cents(invoices, expected):
    instance = InvoiceStats()
    instance.add_invoices(invoices)

    median = instance.get_median()

    assert median == expected


@pytest.mark.parametrize(
    "invoices,expected",
    [([11.05, 11.06, 11.07, 11.08], 11.06), ([100.00, 200.00], 150.00),],
)
def test_get_mean_calculates_mean_and_rounds_down_half_cents(invoices, expected):
    instance = InvoiceStats()
    instance.add_invoices(invoices)

    mean = instance.get_mean()

    assert mean == expected


@pytest.mark.parametrize("invoice", [-101.00, 200_000_000.00, 0.0, 2_918_918_999.00])
def test_add_invoice_raises_value_error_for_out_of_range_value(invoice):
    instance = InvoiceStats()
    with pytest.raises(ValueError):
        instance.add_invoice(invoice)


@pytest.mark.parametrize("invoices", [[1.00, -1.00], [6069.00, 200_000_001]])
def test_add_invoices_raises_value_error_for_any_out_of_range_values(invoices):
    instance = InvoiceStats()
    with pytest.raises(ValueError):
        instance.add_invoices(invoices)


@pytest.mark.parametrize("invoice", [101.001, 200.000005, 0.00981, 2_918_918_999.882])
def test_add_invoice_raises_value_error_for_invalid_value(invoice):
    instance = InvoiceStats()
    with pytest.raises(ValueError):
        instance.add_invoice(invoice)


@pytest.mark.parametrize("invoices", [[1.00, 1.001], [6069.00, 200.999]])
def test_add_invoices_raises_value_error_for_any_invalid_value(invoices):
    instance = InvoiceStats()
    with pytest.raises(ValueError):
        instance.add_invoices(invoices)


def test_add_invoices_raises_runtime_error_if_maximum_exceeded():
    maximum = 200
    instance = InvoiceStats(maximum_invoices=maximum)
    with pytest.raises(RuntimeError):
        instance.add_invoices([1.00 for _ in range(maximum + 1)])


def test_add_invoice_raises_runtime_error_if_maximum_exceeded():
    maximum = 200
    instance = InvoiceStats(maximum_invoices=maximum)
    with pytest.raises(RuntimeError):
        for _ in range(maximum + 1):
            instance.add_invoice(1.00)
