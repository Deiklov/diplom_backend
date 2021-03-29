import datetime
import logging
from concurrent.futures import ThreadPoolExecutor

import grpc
import numpy as np
# from analyze import StocksAnalyzer
import prediction_pb2 as pb2
from prediction_pb2_grpc import PredictAPIServicer, add_PredictAPIServicer_to_server


def find_outliers(data: np.ndarray):
    """Return indices where values more than 2 standard deviations from mean"""
    out = np.where(np.abs(data - data.mean()) > 2 * data.std())
    # np.where returns a tuple for each dimension, we want the 1st element
    return out[0]


class PredictServer(PredictAPIServicer):
    def __init__(self):
        # self.stockAnalyzer = StocksAnalyzer()
        pass

    def Predict(self, request, context):
        logging.info('detect request size: %s', request.stocks_name)
        logging.info('detect timeseries: %s', request.ended_time)
        logging.info('detect step: %s', request.step)
        ended_time: datetime.datetime = request.ended_time.ToDatetime()
        stocks: str = request.stocks_name
        step: int = request.step

        result = {'time': request.ended_time, 'value': 23.11}
        resp = pb2.Metric(**result)
        result2 = [resp, resp]
        return pb2.PredictionResp(time_series=result2)


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', )
    server = grpc.server(ThreadPoolExecutor())
    add_PredictAPIServicer_to_server(PredictServer(), server)
    port = 9999
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logging.info('server ready on port %r', port)
    server.wait_for_termination()
