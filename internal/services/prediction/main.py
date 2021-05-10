from flask import Flask
from flask import request
import numpy as np
import tinvest as ti
from datetime import datetime, timedelta
import math
from flask import jsonify
import pandas as pd
from sklearn.preprocessing import MinMaxScaler
import os
import tensorflow as tf
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dropout, Dense
from tensorflow.python.client import device_lib

app = Flask(__name__)

SANDBOX_TOKEN = "t.UE-TeGMgnVeOVaoBYl7uk33-QtM9k2KwZTc7VyI1ubJErMxsVQvmYb92eRa157bm6XPjx74NGDIYfxSecNrdEQ"


@app.route('/predict/<figi>', methods=['GET'])
def predict(figi: str):
    interval = request.args.get('interval')
    end_date = request.args.get('to')
    client = ti.SyncClient(SANDBOX_TOKEN, use_sandbox=True)
    print(datetime.now())
    # get candles for last year
    response = client.get_market_candles(figi=figi,
                                         from_=datetime.now() - timedelta(days=365),
                                         to=datetime.now(),
                                         interval=ti.CandleResolution.day)
    candles_data = response.payload.candles
    close_full_y = [float(x.c.real) for x in candles_data]
    time_full_x = [x.time for x in candles_data]
    num_shape = math.ceil(0.9 * len(close_full_y))

    df = pd.DataFrame(list(zip(close_full_y, time_full_x)), columns=['Close', 'Date'])

    train = df.iloc[:num_shape, 0:1].values
    test = df.iloc[num_shape:, 0:1].values

    sc = MinMaxScaler(feature_range=(0, 1))
    train_scaled = sc.fit_transform(train)
    X_train = []
    # Price on next day
    y_train = []
    window = 30
    for i in range(window, num_shape):
        X_train_ = np.reshape(train_scaled[i - window:i, 0], (window, 1))
        X_train.append(X_train_)
        y_train.append(train_scaled[i, 0])
    X_train = np.stack(X_train)
    y_train = np.stack(y_train)
    model = Sequential()

    model.add(LSTM(units=50, return_sequences=True, input_shape=(X_train.shape[1], 1)))
    model.add(Dropout(0.2))

    model.add(LSTM(units=50, return_sequences=True))
    model.add(Dropout(0.2))

    model.add(LSTM(units=50, return_sequences=True))
    model.add(Dropout(0.2))

    model.add(LSTM(units=50))
    model.add(Dropout(0.2))

    model.add(Dense(units=1))

    model.compile(optimizer='adam', loss='mean_squared_error')

    model.fit(X_train, y_train, epochs=30, batch_size=24, verbose=0)
    df_volume = np.vstack((train, test))
    num_2 = df_volume.shape[0] - num_shape + window

    pred_ = df['Close'].iloc[-1].copy()
    prediction_full = []
    df_copy = df.iloc[:, 0:1][1:].values

    for j in range(20):
        df_ = np.vstack((df_copy, pred_))
        train_ = df_[:num_shape]
        test_ = df_[num_shape:]

        df_volume_ = np.vstack((train_, test_))

        inputs_ = df_volume_[df_volume_.shape[0] - test_.shape[0] - window:]
        inputs_ = inputs_.reshape(-1, 1)
        inputs_ = sc.transform(inputs_)

        X_test_2 = []

        for k in range(window, num_2):
            X_test_3 = np.reshape(inputs_[k - window:k, 0], (window, 1))
            X_test_2.append(X_test_3)

        X_test_ = np.stack(X_test_2)
        predict_ = model.predict(X_test_)
        pred_ = sc.inverse_transform(predict_)
        prediction_full.append(pred_[-1][0])
        df_copy = df_[j:]
    df_date = df[['Date']]

    for h in range(20):
        df_date_add = pd.to_datetime(df_date['Date'].iloc[-1]) + pd.DateOffset(days=1)
        df_date_add = pd.DataFrame([df_date_add.strftime("%Y-%m-%d")], columns=['Date'])
        df_date = df_date.append(df_date_add)
    df_date = df_date.reset_index(drop=True)

    result = []
    for ind, val in zip(df_date['Date'][len(close_full_y):], prediction_full):
        mda = dict({"pred_close": float(val), "time": ind})
        result.append(mda)
    return jsonify(result)
    # print(json.dumps(result))


if __name__ == '__main__':
    os.environ["CUDA_VISIBLE_DEVICES"] = "3"  # Or 2, 3, etc. other than 0

    # On CPU/GPU placement
    config = tf.compat.v1.ConfigProto(allow_soft_placement=True, log_device_placement=True)
    config.gpu_options.allow_growth = True
    tf.compat.v1.Session(config=config)

    print(device_lib.list_local_devices())

    app.run(host='0.0.0.0')
