#!/bin/bash
export CUDA_VISIBLE_DEVICES=3
gunicorn -w 8 -b 0.0.0.0:5000 main:app