### 这个程序做什么
批量处理[Bilibili PC客户端](https://apps.microsoft.com/store/detail/%E5%93%94%E5%93%A9%E5%93%94%E5%93%A9/XPDDVC6XTQQKMM)缓存的视频，
还原b站对m4s文件修改的文件头信息，然后调用ffmpeg合并视频与音频，并将合成的mp4文件名指定为实际视频名，统一放置在执行根路径。

例：缓存番剧1-13集后，直接在缓存目录执行，即可获得命名好的的13个.mp4视频文件

#### 注意
1. 缓存视频的LOGO是覆盖的，如果需要无logo版请直接找字幕组的源

### 使用
1. 首先需要安装ffmpeg，见[ffmpeg的安装](#ffmpeg的安装)
2. 在视频缓存目录（可在客户端设置页面查看，默认为`%userprofile%\Videos\bilibili`）执行`.\BilibiliVideoCacheProcess.exe`

`BilibiliVideoCacheProcess.exe --help`显示可用参数


### ffmpeg的安装
可直接使用winget安装`winget install ffmpeg`<sup>[1]</sup>

注意在FFmpeg [Gyan.FFmpeg] 6.0版本，使用这种方式安装后，我测试的两台电脑（win10、win11）均遇到了新开控制台仍然无法找到ffmpeg可执行文件的问题，
解决方法是打开系统属性环境变量窗口后点击确定，然后重新打开控制台即可，猜测是有缓存问题


### 参考
1. https://www.gyan.dev/ffmpeg/builds/
