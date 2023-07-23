import pandas as pd
import matplotlib.pyplot as plt
from statsmodels.tsa.arima.model import ARIMA

def read_csv(file_path):
    date_parser = lambda x: pd.to_datetime(x, format='%b %d, %Y')  
    return pd.read_csv(file_path, parse_dates=['Date'], date_parser=date_parser, index_col='Date')

def fit_arima(data, order):
    model = ARIMA(data, order=order)
    return model.fit()

def forecast_arima(model, steps):
    return model.forecast(steps)

def main():
    file_path = 'data.csv'
    df = read_csv(file_path)

    value_col = 'Value'

    time_series = pd.Series(df[value_col].values, index=df.index)

    order = (1, 1, 1)

    model = fit_arima(time_series, order)

    steps = 7

    forecasted_values = forecast_arima(model, steps)

    last_date = df.index[-1]
    forecasted_dates = pd.date_range(start=last_date + pd.DateOffset(days=1), periods=steps, freq='D')

    forecast_df = pd.DataFrame({value_col: forecasted_values}, index=forecasted_dates)

    print("Forecasted values:")
    print(forecast_df)

    output_csv_path = 'forecasted_values.csv'
    forecast_df.to_csv(output_csv_path, index_label='date', header=['new_cases'])

    plt.figure(figsize=(10, 5))
    plt.plot(df, label='Original Data')
    plt.plot(forecast_df, label='Forecasted Data', linestyle='dashed', marker='o')
    plt.xlabel('Date')
    plt.ylabel('Value')
    plt.legend()

    output_plot_path = 'forecast_plot.png'
    plt.savefig(output_plot_path)

    plt.show()

if __name__ == "__main__":
    main()
