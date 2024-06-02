
VENV = st_env

setup:
	python -m venv $(VENV)
	. ./$(VENV)/bin/activate && pip install streamlit

start:
	. ./$(VENV)/bin/activate && streamlit run $(shell pwd)/app.py
