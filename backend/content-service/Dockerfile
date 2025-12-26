FROM python:3.12-slim

WORKDIR /app

COPY ./ ./

RUN pip install --no-cache-dir uv

EXPOSE 8081

CMD ["uv", "run", "server.py"]