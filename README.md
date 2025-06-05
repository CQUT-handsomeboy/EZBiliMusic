![](./.assets/github-header-image.png)

![Go](https://img.shields.io/badge/Language-Go-blue?logo=go) ![GitHub stars](https://img.shields.io/github/stars/cqut-handsomeboy/EZBiliMusic?style=flat)


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

3.  删除音乐(`/delete`)'

```jsonc
{
    "title":"孤独患者",
    "artist":"陈奕迅"
}
```

## 启动


1. 构建镜像

```shell
docker build -t ezbili-music .
```

2. 运行容器

```shell
docker run -d -p 8080:8080 -v /path/to/host/music:/usr/src/app/output ezbili-music
```

3.  运行音乐服务器

推荐[Navidrome](https://www.navidrome.org/docs/installation/docker/)

```shell
$ docker run -d \
   --name navidrome \
   --restart=unless-stopped \
   --user $(id -u):$(id -g) \
   -v /path/to/host/music:/music \
   -v /path/to/host/data:/data \
   -p 4533:4533 \
   -e ND_LOGLEVEL=info \
   deluan/navidrome:latest
```

# 🙁美中不足

1. m4a添加元信息能力较弱，尤其是配合Taglib，专辑封面和歌词无法添加

2. 即使可以添加专辑封面，如何巧妙实现歌词的刮削问题？非官方音源能被刮削软件识别吗？（考虑直接将Bilibili的字幕转为LRC）

> [!NOTE]  
> 如果您有很好的建议，欢迎为我提供PR和Issue。