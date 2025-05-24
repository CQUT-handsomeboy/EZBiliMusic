![](./.assets/github-header-image.png)

# 💡简介

本项目基于下载器项目[lux](https://github.com/iawia002/lux)，从B站上的MV得到了灵感，可通过保存B站MV的音频部分m4a的方式下载音乐。

通过Dockerfile部署，将output文件夹映射到宿主机与Navidrome Container的音乐目录保持一致。

添加一个对外的接口可以从云端控制下载器，间接添加Navidrome乐库。