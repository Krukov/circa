FROM python:3.9-slim

ENV APP_HOME="/app"
WORKDIR ${APP_HOME}

COPY requirements.txt requirements.txt

RUN pip install -U wheel pip && \
    pip install -r requirements.txt

ENV PYTHONUNBUFFERED=1 PYTHONOPTIMIZE=2 PYTHONFAULTHANDLER=1 PYTHONDONTWRITEBYTECODE=1

COPY . ${APP_HOME}

CMD ["python", "app.py"]
