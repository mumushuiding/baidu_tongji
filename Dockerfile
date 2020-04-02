FROM scratch
ADD /baidu_tongji //
ADD /config.json //
EXPOSE 8080
ENTRYPOINT [ "/baidu_tongji" ]