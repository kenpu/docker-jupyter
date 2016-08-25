FROM continuumio/anaconda
COPY jupyter_notebook_config.py /
RUN mkdir /notebooks
WORKDIR /notebooks
ENTRYPOINT jupyter notebook --config /jupyter_notebook_config.py --NotebookApp.notebook_dir=/notebooks
EXPOSE 8888

