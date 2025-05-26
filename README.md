![](./.assets/github-header-image.png)

# 💡简介

本项目基于下载器项目[lux](https://github.com/iawia002/lux)，从B站上的MV得到了灵感，可通过保存B站MV的音频部分m4a的方式下载音乐。

通过Dockerfile部署，将output文件夹映射到宿主机与Navidrome Container的音乐目录保持一致。

添加一个对外的接口可以从云端控制下载器，间接添加Navidrome乐库。

# 🚀开始使用

## POST请求格式

1.  下载(`/download`)

```jsonc
{
    "aid":113893791760835,
    "bvid":"BV1YmFPe4EnY",
    "cid":28086437741,
    "title":"孤独患者",
    "artist":"陈奕迅"
}
```

2.  请求元信息(`/metadata`)

```jsonc
{
    "url":"https://www.bilibili.com/video/BV1YmFPe4EnY"
}
```

3.  启动

```shell
docker build -t ezbili-music .
docker run -d -p 8080:8080 -v /path/to/host/music:/usr/src/app/output ezbili-music
```