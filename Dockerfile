## 打包前需要 set GOOS=linux,然后go build
FROM scratch
ADD /baidu_tongji //
ADD /config.json //
EXPOSE 8080
ENTRYPOINT [ "/baidu_tongji" ]