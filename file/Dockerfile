FROM jrottenberg/ffmpeg:4.1-alpine

# ENTRYPOINT ["/bin/bash"]


RUN mkdir /app
RUN mkdir /app/storage
RUN mkdir /app/storage/temp
RUN mkdir /app/storage/photos
RUN mkdir /app/storage/videos
RUN mkdir /app/storage/musics

COPY bin/fileService /app


WORKDIR /app
RUN chmod +x /app/fileService
# CMD [ "./fileService" ]
ENTRYPOINT [ "./fileService" ]