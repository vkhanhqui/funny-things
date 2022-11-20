# Copy Command Dockerfile


### COPY has two forms:
```bash
COPY [--chown=<user>:<group>] <src>... <dest>
COPY [--chown=<user>:<group>] ["<src>",... "<dest>"]
```


### The problem
Seems we can just use one COPY command in each Dockerfile. It's quite inconvenient, so I thought a trick (Please correct me if I'm wrong by creating a PR)
```bash
COPY ["src_1", "src_2", "root"]
RUN mv root/src_1 /dest_src1/ && \
    mv root/src_2 /dest_src2/
```


### FYI, RUN has 2 forms:
```bash
RUN <command> (shell form, the command is run in a shell, which by default is /bin/sh -c on Linux or cmd /S /C on Windows)
RUN ["executable", "param1", "param2"] (exec form)
```