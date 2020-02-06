"""Aggregations and statistics of processed invoices"""

import logging
import decimal

import pandas as pd


class InvoiceStats:
    """Calculates mean and median value of processed invoices"""

    def __init__(self, maximum_invoices=20_000_000):
        decimal.getcontext().rounding = decimal.ROUND_HALF_DOWN
        logging.info("Creating new invoice table")
        self.invoices = pd.DataFrame(columns=["value"])
        self.maximum_invoices = maximum_invoices
        self.invoice_count = 0

    def add_invoice(self, value):
        """add single invoice to collection
        
        Parameters
        ----------
        value : float
            Two decimal point precision representing dollars and cents value of 
            invoice.
        """
        logging.info(f"Adding single invoice")
        if self.invoice_count + 1 > self.maximum_invoices:
            raise RuntimeError("Adding another invoice would exceed maximum")
        self.invoices = self.invoices.append(
            {"value": self._validate(value)}, ignore_index=True
        )
        self.invoice_count += 1

    def add_invoices(self, invoices):
        """add multiple invoices to collection
        
        Parameters
        ----------
        invoices : list
            List of two decimal point precision floats representing dollars and cents 
            value of invoice

        """
        num_invoices = len(invoices)
        logging.info(f"Adding {num_invoices} invoices")
        if self.invoice_count + num_invoices > self.maximum_invoices:
            raise RuntimeError(
                "Adding another {num_invoices} invoices would exceed maximum"
            )
        self.invoices = self.invoices.append(
            [{"value": self._validate(value)} for value in invoices], ignore_index=True,
        )
        self.invoice_count += num_invoices

    def clear(self):
        """Resets state of class and removes all invoices from collection"""
        logging.info("Removing all invoice data")
        self.__init__()

    def get_median(self):
        """calculate median of all values of invoices
        
        Returns
        -------
        float
            Two decimal point precision representing dollars and cents value of median
            of invoice values. Half a cent is rounded down post calculation.
        """
        logging.info("Calculating median")
        median = self.invoices["value"].median()

        return self._round_half_down(median)

    def get_mean(self):
        """calculate mean of all values of invoices
        
        Returns
        -------
        float
            Two decimal point precision representing dollars and cents value of mean
            of invoice values. Half a cent is rounded down post calculation.
        """
        logging.info("Calculating mean")
        mean = self.invoices["value"].mean()

        return self._round_half_down(mean)

    @staticmethod
    def _round_half_down(value):
        """rounds float to two decimal places and rounds half cent down.
        
        To avoid float errors when rounding half a cent down, the dollar amount is first
        rounded to six decimal places. i.e. One millionth of a dollar is deemed 
        insignificant.

        Parameters
        ----------
        value : float
            value to round
        
        Returns
        -------
        float
            rounded value
        """
        decimal_value = decimal.Decimal(round(value, 6))
        rounded_value = float(decimal_value.quantize(decimal.Decimal("1.00")))
        return rounded_value

    @staticmethod
    def _validate(value):
        """Raise a value error for invalid invoice values"""
        if value <= 0 or value >= 200_000_000:
            raise ValueError(f"Invoice: {value} not between 0 and 200,000,000")
        if round(value, 2) != value:
            raise ValueError(f"Invoice: {value} is not valid (dollars and cents) value")
        return value
