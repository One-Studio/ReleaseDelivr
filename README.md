# ReleaseDelivr
> 基于 GoFrame + GitHub Actions + JsDelivr 的**免费**程序发布版本（Release）自动搬运与整理、更新API和加速下载的工具。如果工具对您有所帮助，请考虑捐赠 JsDelivr 以回报它所提供的服务。

尽管这是一个网络发达的时代，但在如中国大陆和火星等地访问、下载一些好用的工具却非常的扑街，幸苦的搬运工们总是做着重复劳动。不如把这种重复的工作交给机器去做，全自动搬运整理，用户可以极速下载程序，服务器和工具也可很方便的获取最新版本和下载链接等信息。

重点是，这个过程是**免费的、持续的、自动的**，无需手动配置服务器而是依靠 GitHub Actions 定时更新，又得益于JsDelivr的服务，几乎任何地方都应可高速访问这些资源（火星idk）。

# 功能
定时查询或者用dispatch触发，获取并更新release文件到最新版本，修改对应的API文件，利用GitHub Actions的功能Push和发布Release

# 使用方法
1. 前往release页面下载linux64版本，打开 `./.github/workflow/CI.yml` ，修改邮箱和GitHub用户名，最后一行的更新信息可以魔改一下
2. 如果可以联系目标仓库添加GitHub Actions，可以按照 `./.github/workflow/publish.yml` 的说明进行配置，yml文件放在目标仓库，之后就不用没多少分钟查看一次了
3. 配置 `config.json` ：
    - 目标为GitHub仓库则填写仓库的 `owner` 和 `repo`
    - 目标非GitHub仓库则修改 `TargetGH` 为 `false`，`TargetAPI` 是能获取版本号的接口链接，`TargetDLink` 是release文件的下载链接
    - 同理给归档/搬运的仓库仓库填写对应参数，`ArchiverAPI` 是获取搬运仓库的版本号API
    - `Format` = 1 代表.7z | 2 代表.zip
    - `CompRatio` = 压缩率 1 快速 ｜ 2 标准 ｜ 3 极限
    - `DistPath` = 搬运得到的文件存放的位置，默认 `dist`，不用修改
    - `Filter` = release附件的文件过滤器，包含它们的文件才会被下载
    - `DeleteFilter` = 对20MB以上文件生效，名字包含 `Index` 的文件才会触发，先解压到临时位置，删除 `List` 中指定的所有文件/文件夹，然后再次打包，保持文件名不变
4. 上传到仓库，**务必发布一个release更新，版本号越低越好，比如v0.0.1**