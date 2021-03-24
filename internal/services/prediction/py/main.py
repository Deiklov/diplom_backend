import logging
from concurrent.futures import ThreadPoolExecutor

import grpc
import numpy as np
from prediction_pb2 import PredictionResp
from prediction_pb2_grpc import PredictAPIServicer, add_PredictAPIServicer_to_server


def find_outliers(data: np.ndarray):
    """Return indices where values more than 2 standard deviations from mean"""
    out = np.where(np.abs(data - data.mean()) > 2 * data.std())
    # np.where returns a tuple for each dimension, we want the 1st element
    return out[0]


class PredictServer(PredictAPIServicer):
    def Predict(self, request, context):
        logging.info('detect request size: %s', request.stocks_name)
        logging.info('detect timeseries: %s', request.ended_time)
        logging.info('detect step: %s', request.step)
        # Convert metrics to numpy array of values only

        data = np.fromiter((m.value for m in request.metrics), dtype='float64')
        indices = find_outliers(data)
        logging.info('found %d outliers', len(indices))
        resp = PredictionResp(indices=indices)
        return resp


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', )
    server = grpc.server(ThreadPoolExecutor())
    add_PredictAPIServicer_to_server(PredictServer(), server)
    port = 9999
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logging.info('server ready on port %r', port)
    server.wait_for_termination()
