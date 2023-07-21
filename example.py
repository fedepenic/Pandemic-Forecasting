import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from statsmodels.tsa.arima.model import ARIMA

def read_csv(file_path):
    date_parser = lambda x: pd.to_datetime(x, format='%d %B %Y')
    return pd.read_csv(file_path, parse_dates=['Date'], date_parser=date_parser, index_col='Date')

def fit_arima(data, order):
    model = ARIMA(data, order=order)
    return model.fit()

def forecast_arima(model, steps):
    return model.forecast(steps)

def main():
    # Replace 'data.csv' with the path to your CSV file
    file_path = 'data.csv'
    df = read_csv(file_path)

    # Assuming the CSV file has columns named 'Date' and 'Value'
    value_col = 'Value'

    # Create a time series from the 'Value' column
    time_series = pd.Series(df[value_col].values, index=df.index)

    # ARIMA model order (p, d, q)
    order = (1, 1, 1)

    # Fit the ARIMA model
    model = fit_arima(time_series, order)

    # Number of steps ahead to forecast
    steps = 7

    # Forecast the next 'steps' values
    forecasted_values = forecast_arima(model, steps)

    # Generate dates for the forecasted values
    last_date = df.index[-1]
    forecasted_dates = pd.date_range(start=last_date + pd.DateOffset(days=1), periods=steps, freq='D')

    # Combine the forecasted dates and values into a dataframe
    forecast_df = pd.DataFrame({value_col: forecasted_values}, index=forecasted_dates)

    print("Forecasted values:")
    print(forecast_df)

    # Optionally, you can plot the forecasted values
    plt.figure(figsize=(10, 5))
    plt.plot(df, label='Original Data')
    plt.plot(forecast_df, label='Forecasted Data', linestyle='dashed', marker='o')
    plt.xlabel('Date')
    plt.ylabel('Value')
    plt.legend()
    plt.show()

if __name__ == "__main__":
    main()
